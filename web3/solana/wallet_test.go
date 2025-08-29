package solana

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/chain-products-org/goal/testx"
)

var (
	puk           = "GLfZ1AbfccfsKpLiwPxG83KKWfnPzDvpTkoZ3XZmwMBA"
	prk           = "jUG6je2Be7qGoQkecqYywQT5pjjhQ7M4FtXbLB5BsHkFyskEgM6Ryv4WAS7bxBv6Wjd4wyg5kZniDrtx9Eb6GEx"
	to            = "9Fjk6CFddDfuc1hVeWUc19LnnrbgX3yVF9wC5fkrJpB9"
	helius_apiKey = "21af6b88-9dae-45b5-88db-39dc20a0a6db"
	mobula_apikey = "8f42eea7-8167-4d12-9cf3-e4fb18de2d6e"
)

func TestCreateWallet(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("create wallet")
	wallet := CreateWallet()
	tl.Log(wallet)
	tl.Require(wallet != nil, "wallet should not be nil")
	tl.Require(wallet.PublicKey != "", "wallet public key should not be empty")
	tl.Require(wallet.PrivateKey != "", "wallet private key should not be empty")
}

func TestImportFromPrivateKey(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("import from private key")
	wallet, err := ImportFromPrivateKey("/Home/ubuntu/.config/solana/id.json")
	tl.Log(wallet)
	tl.NoErrf(err, "import should no error")
	tl.Require(wallet != nil, "wallet should not be nil")
	tl.Require(wallet.PublicKey != "", "wallet public key should not be empty")
	tl.Require(wallet.PrivateKey != "", "wallet private key should not be empty")
}

func TestSOLWallet_GetAirdrop(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("GetAirdrop")
	sWallet := &SOLWallet{PublicKey: to}
	hash, err := sWallet.GetAirdrop(2, rpc.DevNet_RPC)
	tl.Log(hash)
	tl.NoErrf(err, "should be no error")
	tl.Require(hash != "", "should return hash")
}

func TestSOLWallet_TransferToWaitConfirm(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("TransferToWaitConfirm")
	tl.Title("check balance is sufficent")
	sWallet := &SOLWallet{PublicKey: puk, PrivateKey: prk}
	bret, err := GetBalance(sWallet.PublicKey, rpc.DevNet_RPC)
	tl.NoErrf(err, "should be no error")
	balance, acc := bret.Float64()
	tl.Log("balance:", balance, "acc:", acc)
	if balance*float64(solana.LAMPORTS_PER_SOL) < 5000 {
		tl.Log("insufficient balance, try to get airdrop")
		var sol float64 = 1
		_, err := sWallet.GetAirdrop(sol, rpc.DevNet_RPC)
		tl.NoErrf(err, "should be no error when GetAirdrop")
	}
	tl.Title("do transfer")
	amt := 0.0000001 // transfer amt
	ret, err := sWallet.TransferToWaitConfirm(to, uint64(amt*float64(solana.LAMPORTS_PER_SOL)), rpc.DevNet_RPC, rpc.DevNet_WS)
	tl.NoErrf(err, "should be no error")
	tl.Require(ret != "", "should return hash")
}

func TestSOLWallet_TransferToWithoutConfirm(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("TransferToWithoutConfirm")
	tl.Title("check balance is sufficent")
	sWallet := &SOLWallet{PublicKey: puk, PrivateKey: prk}
	bret, err := GetBalance(sWallet.PublicKey, rpc.DevNet_RPC)
	tl.NoErrf(err, "should be no error")
	balance, acc := bret.Float64()
	tl.Log("balance:", balance, "acc:", acc)
	if balance*float64(solana.LAMPORTS_PER_SOL) < 5000 {
		tl.Log("insufficient balance, try to get airdrop")
		var sol float64 = 1
		_, err := sWallet.GetAirdrop(sol, rpc.DevNet_RPC)
		tl.NoErrf(err, "should be no error when GetAirdrop")
	}
	tl.Title("do transfer")
	amt := 0.0000001 // transfer amt
	ret, err := sWallet.TransferToWithoutConfirm(to, uint64(amt*float64(solana.LAMPORTS_PER_SOL)), rpc.DevNet_RPC)
	tl.NoErr(err)
	tl.Require(ret != "", "should return hash")
}

func TestGetBalance(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("GetBalance")
	bal, err := GetBalance(puk, rpc.DevNet_RPC)
	tl.NoErr(err)
	fmt.Println("◎", bal.Text('f', 10))
	tl.Log("balance:", bal.Text('f', 10))
	v, _ := bal.Float64()
	tl.Require(v >= 0, "balance should >= 0")
}

func TestGetSolPrice(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("GetSolPriceMobula")
	price, err := GetSolPriceMobula(mobula_apikey)
	tl.NoErr(err)
	tl.Log(price)
	tl.Require(price > 0, "price should be valid")
}

func TestGetTransaction(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("GetTransactionInfo")
	tx := "4YBb8sCBNZXPsfabt1dR7e78NFSyE7zZanfPU144Zx173cT6r2WhjohAsxNtwWt8Fd5U79LfGWj36cQDDsuY7yhr"
	state, from, to, amount, err := GetTransactionInfo(tx, "https://mainnet.helius-rpc.com/?api-key="+helius_apiKey)
	tl.NoErr(err)
	tl.Log(state, from, to, amount)
	tl.Require(state == Success, "state should be success")
}

func TestTransferSPLToken(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("TransferSPLToken")
	sWallet := &SOLWallet{PublicKey: puk, PrivateKey: prk}
	hash, err := sWallet.TransferSPLToken(puk, to, 10, rpc.DevNet_RPC)
	tl.NoErr(err)
	tl.Log(hash)
	tl.Require(hash != "", "should return tx hash")
}

func TestGetTransactionList(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("GetTransactionList")
	r, err := GetTransactionList(puk, rpc.DevNet_RPC)
	tl.NoErr(err)
	tl.Log(len(r))
	tl.Require(len(r) > 0, "tx list should have more elements")
}
