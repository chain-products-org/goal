package util

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/gophero/goal/httpx"
)

func TestImportWallet(t *testing.T) {
	legacy := ImportPrivateKeyToLegacy("your private key", &chaincfg.TestNet3Params)
	p2tr := ImportPrivateKeyToP2TR("your private key", &chaincfg.TestNet3Params)
	fmt.Printf("legacy addr: %s \n", legacy)
	fmt.Printf("P2TR addr: %s \n", p2tr)
}

func TestCreateWalletP2TR(t *testing.T) {
	legacyWallet := CreateWalletP2TR(&chaincfg.TestNet3Params)
	fmt.Println(legacyWallet)
	p2tr := ImportPrivateKeyToP2TR(legacyWallet.privateKey, &chaincfg.TestNet3Params)
	fmt.Println(p2tr)
}

func TestGetAddressTransactions(t *testing.T) {
	httpx.Get("https://mempool.space/testnet/api/address/xxx/txs", func(resp *http.Response) {
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bs))
	}, func(err error) {
		if err != nil {
			panic(err)
		}
	})
}
