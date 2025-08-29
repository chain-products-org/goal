package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
	"github.com/chain-products-org/goal/errorx"
	"github.com/chain-products-org/goal/httpx"
)

type SOLWallet struct {
	PublicKey  string
	PrivateKey string
}

func CreateWallet() *SOLWallet {
	account := solana.NewWallet()
	return &SOLWallet{PrivateKey: account.PrivateKey.String(), PublicKey: account.PublicKey().String()}
}

func ImportFromPrivateKey(filePath string) (*SOLWallet, error) {
	privKey, err := solana.PrivateKeyFromSolanaKeygenFile(filePath)
	if err != nil {
		return nil, err
	}
	account, err := solana.WalletFromPrivateKeyBase58(privKey.String())
	if err != nil {
		return nil, err
	}
	return &SOLWallet{PrivateKey: account.PrivateKey.String(), PublicKey: account.PublicKey().String()}, nil
}

func (w *SOLWallet) GetAirdrop(sol float64, rpcUrl string) (string, error) {
	client := rpc.New(rpcUrl)
	// Airdrop 1 SOL to the new account:
	out, err := client.RequestAirdrop(
		context.TODO(),
		solana.MustPublicKeyFromBase58(w.PublicKey),
		uint64(float64(solana.LAMPORTS_PER_SOL)*sol),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// TransferToWaitConfirm transfer SOL to another account and wait for confirmation
// amount: in lamports (1 SOL = 1000000000 lamports)
func (w *SOLWallet) TransferToWaitConfirm(to string, amount uint64, rpcUrl, wsUrl string) (string, error) {
	accountFrom, err := solana.PrivateKeyFromBase58(w.PrivateKey)
	if err != nil {
		return "", err
	}
	rpcClient := rpc.New(rpcUrl)
	tx, err := signTransaction(to, amount, rpcClient, accountFrom)
	if err != nil {
		return "", nil
	}
	// Pretty print the transaction:
	wsClient, err := ws.Connect(context.Background(), wsUrl)
	if err != nil {
		return "", nil
	}
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))
	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		wsClient,
		tx,
	)
	if err != nil {
		return "", nil
	}
	spew.Dump(sig)
	return sig.String(), nil
}

func (w *SOLWallet) TransferToWithoutConfirm(to string, amount uint64, rpcUrl string) (string, error) {
	accountFrom, err := solana.PrivateKeyFromBase58(w.PrivateKey)
	if err != nil {
		return "", nil
	}
	rpcClient := rpc.New(rpcUrl)
	tx, err := signTransaction(to, amount, rpcClient, accountFrom)
	if err != nil {
		return "", nil
	}
	// Or just send the transaction WITHOUT waiting for confirmation:
	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentConfirmed,
	}
	sig, err := rpcClient.SendTransactionWithOpts(
		context.TODO(),
		tx,
		opts,
	)
	if err != nil {
		return "", err
	}
	spew.Dump(sig)
	return sig.String(), nil
}

func signTransaction(to string, amount uint64, rpcClient *rpc.Client, accountFrom solana.PrivateKey) (*solana.Transaction, error) {
	// Get the recent blockhash:
	recent, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	// Create the transaction:
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				accountFrom.PublicKey(),
				solana.MustPublicKeyFromBase58(to),
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		return nil, err
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	return tx, nil
}

func (w *SOLWallet) TransferSPLToken(tokenSource string, to string, amount uint64, rpcUrl string) (string, error) {
	rpcClient := rpc.New(rpcUrl)
	tx, err := w.createSPLToken(tokenSource, to, amount, rpcClient)
	if err != nil {
		return "", err
	}
	// Pretty print the transaction:
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SPL Token"))
	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		nil,
		tx,
	)
	if err != nil {
		return "", err
	}
	spew.Dump(sig)
	return sig.String(), nil
}

func (w *SOLWallet) createSPLToken(tokenSource string, to string, amount uint64, rpcClient *rpc.Client) (*solana.Transaction, error) {
	accountFrom, err := solana.PrivateKeyFromBase58(w.PrivateKey)
	if err != nil {
		return nil, err
	}
	// Get the recent blockhash:
	recent, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	params := new(token.Transfer)
	source := &solana.AccountMeta{PublicKey: solana.MustPublicKeyFromBase58(tokenSource), IsSigner: false, IsWritable: true}
	toAct := &solana.AccountMeta{PublicKey: solana.MustPublicKeyFromBase58(to), IsSigner: false, IsWritable: true}
	signer := &solana.AccountMeta{PublicKey: accountFrom.PublicKey(), IsSigner: true, IsWritable: false}
	params.SetAccounts([]*solana.AccountMeta{source, toAct, signer})
	params.SetAmount(amount)
	// param := params.SetSourceAccount(solana.MustPublicKeyFromBase58(tokenSource)).SetDestinationAccount(solana.MustPublicKeyFromBase58(to)).SetOwnerAccount(accountFrom.PublicKey()).SetAmount(amount)

	// Create the transaction:
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			params.Build(),
		},
		recent.Value.Blockhash,
		// solana.TransactionPayer(solana.MustPublicKeyFromBase58(tokenSource)),
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		return nil, err
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	return tx, nil
}

func GetBalance(addr string, rpcUrl string) (*big.Float, error) {
	client := rpc.New(rpcUrl)
	pubKey := solana.MustPublicKeyFromBase58(addr)
	out, err := client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}
	spew.Dump(out.Value)
	lamportsOnAccount := new(big.Float).SetUint64(out.Value)
	solBalance := new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
	return solBalance, nil
}

func FromLamports(v uint64) *big.Float {
	lamports := new(big.Float).SetUint64(v)
	return new(big.Float).Quo(lamports, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
}

func ToLamports(v *big.Float) uint64 {
	result, _ := new(big.Float).Mul(v, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL)).Uint64()
	return result
}

func GetSolPriceMobula(key string) (float64, error) {
	const CoinPriceUrl = "https://api.mobula.io/api/1/market/data?asset=%s"
	const SOL = "Solana"
	sol := make(map[string]interface{})
	_, err := httpx.GetJson(fmt.Sprintf(CoinPriceUrl, SOL), &sol, map[string]string{"Content-Type": "application/json", "Authorization": key})
	if err != nil {
		return 0, err
	}
	sol_price := (sol["data"]).(map[string]interface{})["price"]
	if sol_price.(float64) > 0 {
		return sol_price.(float64), nil
	}
	return 0, nil
}

const (
	Failed  = -1
	Pending = 1
	Success = 2
)

func GetTransactionInfo(signature string, rpcUrl string) (state int, fromVal string, toVal string, amountVal uint64, err error) {
	client := rpc.New(rpcUrl)
	sigure := solana.MustSignatureFromBase58(signature)
	sd := uint64(0)
	resp, err := client.GetTransaction(context.Background(), sigure, &rpc.GetTransactionOpts{MaxSupportedTransactionVersion: &sd})
	if err != nil {
		// TODO add error log
		return Pending, "", "", 0, err
	}
	chainErr := resp.Meta.Status["Err"]
	// chainErr := resp.Meta.Status["Ok"]

	// receiptAmount
	amount := resp.Meta.PostBalances[1] - resp.Meta.PreBalances[1]
	if chainErr != nil {
		return Failed, "", "", 0, nil
	}
	by, _ := resp.Transaction.MarshalJSON()
	mp := &TransactionModel{}
	err = json.Unmarshal(by, &mp)
	if err != nil {
		return Pending, "", "", 0, err
	}
	// receiptAddress is the to address
	to := mp.Message.AccountKeys[1]
	from := mp.Message.AccountKeys[0]
	return Success, from, to, amount, nil
}

func GetNFTAddrFromTransaction(signature, rpcUrl string) (int, string, error) {
	client := rpc.New(rpcUrl)
	sigure := solana.MustSignatureFromBase58(signature)
	sd := uint64(0)
	resp, err := client.GetTransaction(context.Background(), sigure, &rpc.GetTransactionOpts{MaxSupportedTransactionVersion: &sd})
	if err != nil {
		// TODO add error log
		return Pending, "", err
	}
	chainErr := resp.Meta.Status["Err"]
	// chainErr := resp.Meta.Status["Ok"]

	if chainErr != nil {
		return Failed, "", nil
	}
	by, _ := resp.Transaction.MarshalJSON()
	mp := &TransactionModel{}
	err = json.Unmarshal(by, &mp)
	if err != nil {
		return Pending, "", err
	}
	// receiptAddress is the to address
	tokenAddr := mp.Message.AccountKeys[1]
	return Success, tokenAddr, nil
}

func GetTransactionSate(sig string, rpcUrl string) (bool, error) {
	client := rpc.New(rpcUrl)
	sigure := solana.MustSignatureFromBase58(sig)
	out, err := client.GetSignatureStatuses(context.Background(), true, sigure)
	if err != nil {
		return false, err
	}
	if out.Value[0] == nil {
		return false, errorx.New("transaction not found")
	}
	if out.Value[0].ConfirmationStatus == "finalized" {
		return true, nil
	}
	return false, nil
}

func GetTransactionList(addr string, rpcUrl string) (rpc.GetConfirmedSignaturesForAddress2Result, error) {
	client := rpc.New(rpcUrl)
	pubKey := solana.MustPublicKeyFromBase58(addr)
	lim := 10
	out, err := client.GetSignaturesForAddressWithOpts(
		context.TODO(),
		pubKey,
		&rpc.GetSignaturesForAddressOpts{
			Commitment: rpc.CommitmentConfirmed, // query confirmed
			Limit:      &lim,
		},
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type MoralisPriceResponse struct {
	UsdPrice float64 `json:"usdPrice"`
}

func GetSolPriceFromMoralis(key string) (float64, error) {
	const CoinPriceUrl = "https://deep-index.moralis.io/api/v2/erc20/%s/price?chain=eth"
	const WETH = "0xD31a59c85aE9D8edEFeC411D448f90841571b89c"
	ret := &MoralisPriceResponse{}
	_, err := httpx.GetJson(fmt.Sprintf(CoinPriceUrl, WETH), ret, map[string]string{"accept": "application/json", "X-API-Key": key})
	if err != nil {
		return 0, err
	}
	if ret.UsdPrice > 0 {
		return ret.UsdPrice, nil
	}
	return 0, nil
}

type TransactionModel struct {
	Signature string  `json:"signature"`
	Message   Message `json:"message"`
}

type Message struct {
	AccountKeys []string `json:"accountKeys"`
}
