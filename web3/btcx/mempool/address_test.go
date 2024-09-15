package mempool_test

import (
	"bytes"
	"encoding/hex"
	"github.com/gophero/goal/web3/btcx"
	"github.com/gophero/goal/web3/btcx/mempool"
	"github.com/gophero/goal/web3/btcx/util"
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/stretchr/testify/assert"
)

var addr = "your addr"

func TestGetAddress(t *testing.T) {
	client := mempool.NewClient(btcx.TestNet)
	r, err := client.GetAddress(addr)
	if err != nil {
		t.Error(err)
	}
	t.Logf("result: %v", r)
	assert.True(t, r.Address != "")
}

func TestListTxs(t *testing.T) {
	client := mempool.NewClient(btcx.TestNet)
	addr := addr
	txs, err := client.ListTxs(addr, "")
	if err != nil {
		panic(err)
	}
	assert.True(t, len(txs) > 0)
}

func TestListTxsConfirmed(t *testing.T) {
	client := mempool.NewClient(btcx.TestNet)
	addr := addr
	txs, err := client.ListTxsConfirmed(addr, "")
	if err != nil {
		t.Errorf("test failed: %v", err)
	}
	assert.True(t, len(txs) == 25)
	txs, err = client.ListTxsConfirmed(addr, addr)
	if err != nil {
		t.Errorf("test failed: %v", err)
	}
	assert.True(t, len(txs) == 25)
	assert.True(t, txs[0].TxHash().String() == addr)
	var firstTx = txs[0]
	var w bytes.Buffer
	if err := firstTx.Serialize(&w); err != nil {
		t.Errorf("serialize error: %v", err)
	}
	assert.True(t, firstTx.TxIn[0].PreviousOutPoint.Hash.String() == "xxx")
	assert.True(t, firstTx.TxIn[0].PreviousOutPoint.Index == 1)
}

func TestListUnspent(t *testing.T) {
	// https://mempool.space/signet/api/address/tb1p8lh4np5824u48ppawq3numsm7rss0de4kkxry0z70dcfwwwn2fcspyyhc7/utxo
	client := mempool.NewClient(btcx.TestNet)
	addr := addr
	unspentList, err := client.ListUnspent(addr)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(len(unspentList))
		for _, output := range unspentList {
			t.Log(output.Outpoint.Hash.String(), "    ", output.Outpoint.Index)
			assert.True(t, output.Outpoint.Hash.String() != "")
		}
	}
}

func TestListUnconfirmed(t *testing.T) {
	client := mempool.NewClient(btcx.TestNet)
	addr := addr
	list, err := client.ListTxsUnconfirmed(addr)
	if err != nil {
		t.Error(err)
	} else {
		assert.True(t, len(list) == 0)
	}
}

func TestValidAddress(t *testing.T) {
	client := mempool.NewClient((btcx.TestNet))
	addr := addr
	va, err := client.ValidAddress(addr)
	if err != nil {
		t.Errorf("test failed: %v", err)
	} else {
		assert.True(t, va.Isvalid)
		assert.True(t, va.Iswitness)
		assert.True(t, va.Isscript)
		ad, _ := util.ParseAddress(addr, btcx.TestNet)
		pkScript, _ := txscript.PayToAddrScript(ad)
		pubscriptstr := hex.EncodeToString(pkScript)
		t.Log(pubscriptstr)
		assert.True(t, pubscriptstr == va.ScriptPubKey)
	}
}
