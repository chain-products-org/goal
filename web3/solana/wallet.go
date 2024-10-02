package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/httpx"
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
		panic(err)
	}
	rpcClient := rpc.New(rpcUrl)
	tx := signTransaction(to, amount, rpcClient, accountFrom)
	// Pretty print the transaction:
	wsClient, err := ws.Connect(context.Background(), wsUrl)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	spew.Dump(sig)
	return sig.String(), nil
}

func (w *SOLWallet) TransferToWithoutConfirm(to string, amount uint64, rpcUrl string) (string, error) {
	accountFrom, err := solana.PrivateKeyFromBase58(w.PrivateKey)
	if err != nil {
		panic(err)
	}
	rpcClient := rpc.New(rpcUrl)
	tx := signTransaction(to, amount, rpcClient, accountFrom)
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

func signTransaction(to string, amount uint64, rpcClient *rpc.Client, accountFrom solana.PrivateKey) *solana.Transaction {
	// Get the recent blockhash:
	recent, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	return tx
}

func (w *SOLWallet) TransferSPLToken(tokenSource string, to string, amount uint64, rpcUrl string) {
	rpcClient := rpc.New(rpcUrl)
	tx := w.createSPLToken(tokenSource, to, amount, rpcClient)
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
		panic(err)
	}
	spew.Dump(sig)
}

func (w *SOLWallet) createSPLToken(tokenSource string, to string, amount uint64, rpcClient *rpc.Client) *solana.Transaction {
	accountFrom, err := solana.PrivateKeyFromBase58(w.PrivateKey)
	if err != nil {
		panic(err)
	}
	// Get the recent blockhash:
	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	params := new(token.Transfer)
	source := &solana.AccountMeta{PublicKey: solana.MustPublicKeyFromBase58(tokenSource), IsSigner: false, IsWritable: true}
	toAct := &solana.AccountMeta{PublicKey: solana.MustPublicKeyFromBase58(to), IsSigner: false, IsWritable: true}
	signer := &solana.AccountMeta{PublicKey: accountFrom.PublicKey(), IsSigner: true, IsWritable: false}
	//
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
		panic(err)
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	return tx
}

func GetBalance(addr string, rpcUrl string) *big.Float {
	client := rpc.New(rpcUrl)
	pubKey := solana.MustPublicKeyFromBase58(addr)
	out, err := client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}
	// Convert lamports to sol:
	spew.Dump(out.Value) // total lamports on the account; 1 sol = 1000000000 lamports
	lamportsOnAccount := new(big.Float).SetUint64(out.Value)
	solBalance := new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
	// WARNING: this is not a precise conversion.
	// fmt.Println("◎", solBalance.Text('f', 10))
	return solBalance
}

func FromLamports(value uint64) *big.Float {
	lamportsOnAccount := new(big.Float).SetUint64(value)
	return new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
}

func ToLamports(value *big.Float) uint64 {
	result, _ := new(big.Float).Mul(value, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL)).Uint64()
	return result
}

func GetSolPriceMobula(key string) float64 {
	const CoinPriceUrl = "https://api.mobula.io/api/1/market/data?asset=%s"
	const SOL = "Solana"
	var rerr error
	sol := make(map[string]interface{})
	_, err := httpx.GetJson(fmt.Sprintf(CoinPriceUrl, SOL), &sol, map[string]string{"Content-Type": "application/json", "Authorization": key})
	if err != nil {
		log.Printf("request failed: %v", rerr)
		return 0
	}
	sol_price := (sol["data"]).(map[string]interface{})["price"]
	if sol_price.(float64) > 0 {
		return sol_price.(float64)
	}
	return 0
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

func GetTransactionList(addr string, rpcUrl string) {
	client := rpc.New(rpcUrl)
	pubKey := solana.MustPublicKeyFromBase58(addr)
	lim := uint64(10)
	out, err := client.GetConfirmedSignaturesForAddress2(
		context.TODO(),
		pubKey,
		&rpc.GetConfirmedSignaturesForAddress2Opts{
			Limit: &lim,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

type MoralisPriceResponse struct {
	UsdPrice float64 `json:"usdPrice"`
}

func GetSolPriceFromMoralis(key string) float64 {
	const CoinPriceUrl = "https://deep-index.moralis.io/api/v2/erc20/%s/price?chain=eth"
	const WETH = "0xD31a59c85aE9D8edEFeC411D448f90841571b89c"
	ret := &MoralisPriceResponse{}
	var rerr error
	_, err := httpx.GetJson(fmt.Sprintf(CoinPriceUrl, WETH), ret, map[string]string{"accept": "application/json", "X-API-Key": key})
	if err != nil {
		log.Printf("request failed: %v", rerr)
		return 0
	}
	if ret.UsdPrice > 0 {
		return ret.UsdPrice
	}
	return 0
}

type TransactionModel struct {
	Signature string  `json:"signature"`
	Message   Message `json:"message"`
}

type Message struct {
	AccountKeys []string `json:"accountKeys"`
}
