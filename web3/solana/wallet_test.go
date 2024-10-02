package solana

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gophero/goal/testx"
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
	tl.NoErr(err, "import should no error")
	tl.Require(wallet != nil, "wallet should not be nil")
	tl.Require(wallet.PublicKey != "", "wallet public key should not be empty")
	tl.Require(wallet.PrivateKey != "", "wallet private key should not be empty")
}

func TestSOLWallet_GetAirdrop(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("GetAirdrop")
	sWallet := &SOLWallet{PublicKey: "9Fjk6CFddDfuc1hVeWUc19LnnrbgX3yVF9wC5fkrJpB9", PrivateKey: ""}
	hash, err := sWallet.GetAirdrop(2, rpc.DevNet_RPC)
	tl.Log(hash)
	tl.NoErr(err, "should no error")
	tl.Require(hash != "", "should return hash")
}

func TestSOLWallet_TransferToWaitConfirm(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: puk, PrivateKey: prk}
	balance, acc := GetBalance(sWallet.PublicKey, rpc.DevNet_RPC).Float64()
	fmt.Println(acc)
	fmt.Println(balance)
	ret, err := sWallet.TransferToWaitConfirm(to, uint64(balance*float64(solana.LAMPORTS_PER_SOL))-5000, rpc.DevNet_RPC, rpc.DevNet_WS)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(ret)
}

func TestSOLWallet_TransferToWithoutConfirm(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: puk, PrivateKey: prk}
	hex, err := sWallet.TransferToWithoutConfirm(to, uint64(3324147), rpc.DevNet_RPC)
	if err != nil {
		t.Errorf("failed: %v", err)
		t.Fail()
	}
	fmt.Println(hex)
}

func TestGetBalance(t *testing.T) {
	bal := GetBalance(puk, rpc.DevNet_RPC)
	fmt.Println("◎", bal.Text('f', 10))
}

func TestGetSolPrice(t *testing.T) {
	price := GetSolPriceMobula(mobula_apikey)
	fmt.Println(fmt.Sprintf("%f", price))
}

func TestGetTransaction(t *testing.T) {
	tx := "4YBb8sCBNZXPsfabt1dR7e78NFSyE7zZanfPU144Zx173cT6r2WhjohAsxNtwWt8Fd5U79LfGWj36cQDDsuY7yhr"
	state, from, to, amount, err := GetTransactionInfo(tx, "https://mainnet.helius-rpc.com/?api-key="+helius_apiKey)
	fmt.Println(state, from, to, amount, err)
}

func TestTransferSPLToken(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: puk, PrivateKey: prk}
	sWallet.TransferSPLToken(puk, to, 10, rpc.DevNet_RPC)
}

func TestTemp(t *testing.T) {
	// data := []byte{0xcc, 0xdd}
	fmt.Println(Hex2Dec("dd"))
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

func TestGetTransactionList(t *testing.T) {
	GetTransactionList(puk, rpc.DevNet_RPC)
}

func TestURLDecode(t *testing.T) {
	str := "https%3A%2F%2Fblink-flip.onrender.com%2Fapi%2F0.1%2F8MaDk3Nou9jRVturbfnt3egf1aP9p1AjL8wiJ98kH1F%2Fheads"
	decode, _ := url.QueryUnescape(str)
	fmt.Println(decode)
}
