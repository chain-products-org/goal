package solana

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gophero/goal/httpx"

	"github.com/gagliardetto/solana-go/rpc"
)

func TestCreateWallet(t *testing.T) {
	wallet := CreateWallet()
	fmt.Println(wallet)
}

func TestImportFromPrivateKey(t *testing.T) {
	wallet := ImportFromPrivateKey("/Home/ubuntu/.config/solana/id.json")
	fmt.Println(wallet)
}

func TestSOLWallet_GetAirdrop(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: "DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", PrivateKey: "xxxx"}

	// sWallet.GetAirdrop(rpc.DevNet_RPC)
	sWallet.GetAirdrop(rpc.DevNet_RPC)
}

func TestSOLWallet_TransferToWaitConfirm(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: "CVVQeRGEi1bR4a7CVqNxBPCxScVMYPR86SiNPkeFfLFs", PrivateKey: "xxxx"}

	sWallet.TransferToWaitConfirm("DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", uint64(32231112), rpc.DevNet_RPC, rpc.DevNet_WS)
}

func TestSOLWallet_TransferToWithoutConfirm(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: "CVVQeRGEi1bR4a7CVqNxBPCxScVMYPR86SiNPkeFfLFs", PrivateKey: "xxxxx"}
	hex, err := sWallet.TransferToWithoutConfirm("DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", uint64(3324147), rpc.MainNetBeta_RPC)
	if err != nil {
		panic(err)
	}
	fmt.Println(hex)
}

func TestGetBalance(t *testing.T) {
	bal := GetBalance("DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", rpc.MainNetBeta_RPC)
	// bal := GetBalance("DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", rpc.DevNet_RPC)
	fmt.Println("◎", bal.Text('f', 10))
	//	0.0051575720
	//  0.00515757
}

func TestGetSolPrice(t *testing.T) {
	price := GetSolPriceMobula("your api key")
	fmt.Println(fmt.Sprintf("%f", price))
}

func TestGetTransaction(t *testing.T) {
	// prd confirm
	tx := "2uMqKF3MjiQkYy8rnvTyS3aGw3VHefJZJpnUzqWGZchENaoi31dV9uyihWuqmW7p9s8YgacCMjiiSQF2uhwDtu2C"
	state, from, to, amount, err := GetTransactionInfo(tx, "https://mainnet.helius-rpc.com/?api-key=621a800a-5716-4e3e-96ea-c9f954c99ba3")
	fmt.Println(state, from, to, amount, err)
}

func TestTransferSPLToken(t *testing.T) {
	sWallet := &SOLWallet{PublicKey: "DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", PrivateKey: "xxxxx"}
	to := "CVVQeRGEi1bR4a7CVqNxBPCxScVMYPR86SiNPkeFfLFs"
	sWallet.TransferSPLToken("CVVQeRGEi1bR4a7CVqNxBPCxScVMYPR86SiNPkeFfLFs", to, 10, rpc.DevNet_RPC)
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

func TestCallNode(t *testing.T) {
	params := fmt.Sprintf(`{"txId":%d, "collection":"%s", "tokenId":"%s", "type":%d, "rpcUrl":"%s", "to":"%s" }`, 12, "CVVQeRGEi1bR4a7CVqNxBPCxScVMYPR86SiNPkeFfLFs", "0002", 2, rpc.MainNetBeta_RPC, "DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP")
	// TODO 通知node mint nft
	mintHashByte := httpx.PostJson("http://10.10.10.106:3000"+"/inner/mint_to", strings.NewReader(params), func(err error) {
		// err = srv.MintReport(nft, "", model.NFTMintFail)
		fmt.Println("mint nft failed")
	}, nil)

	mp := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(mintHashByte), &mp)
	if err1 != nil {
		fmt.Println("node mint respose err", mintHashByte)
	}
	mintHash := mp["hash"].(string)

	fmt.Println("mintHash", mintHash)
}

func TestGetTransactionList(t *testing.T) {
	GetTransactionList("DqWiEQscZgoxDdJjAoVkZLRKL8EZSL9f2gk1cHwXcPBP", rpc.MainNetBeta_RPC)
}

func TestURLDecode(t *testing.T) {
	str := "https%3A%2F%2Fblink-flip.onrender.com%2Fapi%2F0.1%2F8MaDk3Nou9jRVturbfnt3egf1aP9p1AjL8wiJ98kH1F%2Fheads"
	decode, _ := url.QueryUnescape(str)
	fmt.Println(decode)
}
