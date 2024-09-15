package mempool_test

import (
	"github.com/gophero/goal/btcx"
	"github.com/gophero/goal/btcx/mempool"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func TestGetRawTransaction(t *testing.T) {
	client := mempool.NewClient(btcx.TestNet)
	txId, _ := chainhash.NewHashFromStr("xxx")
	transaction, err := client.GetRawTransaction(txId)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(transaction.TxHash().String())
	}
}
