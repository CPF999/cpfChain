// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package para

import (
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	typesmocks "github.com/33cn/chain33/types/mocks"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
	"github.com/stretchr/testify/mock"
)

var (
	Amount = int64(1 * types.Coin)
	Title  = string("user.p.para.")
	Title2 = string("user.p.test.")
)

func TestFilterTxsForPara(t *testing.T) {
	types.Init(Title, nil)

	//only main txs
	tx0, _ := createMainTx("ticket", "to")
	//only main txs group
	tx1, _ := createMainTx("ticket", "to")
	tx2, _ := createMainTx("token", "to")
	tx12 := []*types.Transaction{tx1, tx2}
	txGroup12, err := createTxsGroup(tx12)
	assert.Nil(t, err)

	//para cross tx group succ
	tx3, _ := createCrossMainTx("toA")
	tx4, err := createCrossParaTx("toB", 4)
	assert.Nil(t, err)
	tx34 := []*types.Transaction{tx3, tx4}
	txGroup34, err := createTxsGroup(tx34)
	assert.Nil(t, err)

	//all para tx group
	tx5, err := createCrossParaTx("toB", 5)
	assert.Nil(t, err)
	tx6, err := createCrossParaTx("toB", 6)
	assert.Nil(t, err)
	tx56 := []*types.Transaction{tx5, tx6}
	txGroup56, err := createTxsGroup(tx56)
	assert.Nil(t, err)

	//para cross tx group fail
	tx7, _ := createCrossMainTx("toA")
	tx8, err := createCrossParaTx("toB", 8)
	assert.Nil(t, err)
	tx78 := []*types.Transaction{tx7, tx8}
	txGroup78, err := createTxsGroup(tx78)
	assert.Nil(t, err)

	tx9, _ := createMainTx("relay", "to")
	//single para tx
	txA, err := createCrossParaTx("toB", 10)
	assert.Nil(t, err)

	//all para tx group
	txB, err := createCrossParaTx("toB", 11)
	assert.Nil(t, err)
	txC, err := createCrossParaTx("toB", 12)
	assert.Nil(t, err)
	txBC := []*types.Transaction{txB, txC}
	txGroupBC, err := createTxsGroup(txBC)
	assert.Nil(t, err)

	//single para tx
	txD, err := createCrossParaTempTx("toB", 10)
	assert.Nil(t, err)

	txs := []*types.Transaction{tx0}
	txs = append(txs, txGroup12...)
	txs = append(txs, txGroup34...)
	txs = append(txs, txGroup56...)
	txs = append(txs, txGroup78...)
	txs = append(txs, tx9, txA)
	txs = append(txs, txGroupBC...)
	txs = append(txs, txD)

	//for i, tx := range txs {
	//	t.Log("tx exec name", "i", i, "name", string(tx.Execer))
	//}

	recpt0 := &types.ReceiptData{Ty: types.ExecOk}
	recpt1 := &types.ReceiptData{Ty: types.ExecOk}
	recpt2 := &types.ReceiptData{Ty: types.ExecOk}
	recpt3 := &types.ReceiptData{Ty: types.ExecOk}
	recpt4 := &types.ReceiptData{Ty: types.ExecOk}
	recpt5 := &types.ReceiptData{Ty: types.ExecPack}
	recpt6 := &types.ReceiptData{Ty: types.ExecPack}

	log7 := &types.ReceiptLog{Ty: types.TyLogErr}
	logs := []*types.ReceiptLog{log7}
	recpt7 := &types.ReceiptData{Ty: types.ExecPack, Logs: logs}
	recpt8 := &types.ReceiptData{Ty: types.ExecPack}

	recpt9 := &types.ReceiptData{Ty: types.ExecOk}
	recptA := &types.ReceiptData{Ty: types.ExecPack}
	recptB := &types.ReceiptData{Ty: types.ExecPack}
	recptC := &types.ReceiptData{Ty: types.ExecPack}
	recptD := &types.ReceiptData{Ty: types.ExecPack}
	receipts := []*types.ReceiptData{recpt0, recpt1, recpt2, recpt3, recpt4, recpt5,
		recpt6, recpt7, recpt8, recpt9, recptA, recptB, recptC, recptD}

	block := &types.Block{Txs: txs}
	detail := &types.BlockDetail{
		Block:    block,
		Receipts: receipts,
	}

	para := &client{}
	rst := para.FilterTxsForPara(detail)
	filterTxs := []*types.Transaction{tx3, tx4, tx5, tx6, txA, txB, txC}
	assert.Equal(t, filterTxs, rst)

}

func createMainTx(exec string, to string) (*types.Transaction, error) {
	param := types.CreateTx{
		To:          to,
		Amount:      Amount,
		Fee:         0,
		Note:        []byte("test"),
		TokenSymbol: "",
		ExecName:    exec,
	}
	transfer := &pt.ParacrossAction{}
	v := &pt.ParacrossAction_AssetTransfer{AssetTransfer: &types.AssetsTransfer{
		Amount: param.Amount, Note: param.GetNote(), To: param.GetTo()}}
	transfer.Value = v
	transfer.Ty = pt.ParacrossActionAssetTransfer

	tx := &types.Transaction{
		Execer:  []byte(param.GetExecName()),
		Payload: types.Encode(transfer),
		To:      address.ExecAddress(param.GetExecName()),
		Fee:     param.Fee,
		Nonce:   rand.New(rand.NewSource(time.Now().UnixNano())).Int63(),
	}

	return tx, nil
}

func createCrossMainTx(to string) (*types.Transaction, error) {
	param := types.CreateTx{
		To:          string(to),
		Amount:      Amount,
		Fee:         0,
		Note:        []byte("test asset transfer"),
		IsWithdraw:  false,
		IsToken:     false,
		TokenSymbol: "",
		ExecName:    pt.ParaX,
	}
	transfer := &pt.ParacrossAction{}
	v := &pt.ParacrossAction_AssetTransfer{AssetTransfer: &types.AssetsTransfer{
		Amount: param.Amount, Note: param.GetNote(), To: param.GetTo()}}
	transfer.Value = v
	transfer.Ty = pt.ParacrossActionAssetTransfer

	tx := &types.Transaction{
		Execer:  []byte(param.GetExecName()),
		Payload: types.Encode(transfer),
		To:      address.ExecAddress(param.GetExecName()),
		Fee:     param.Fee,
		Nonce:   rand.New(rand.NewSource(time.Now().UnixNano())).Int63(),
	}

	return tx, nil
}

func createCrossParaTx(to string, amount int64) (*types.Transaction, error) {
	param := types.CreateTx{
		To:          string(to),
		Amount:      amount,
		Fee:         0,
		Note:        []byte("test asset transfer"),
		IsWithdraw:  false,
		IsToken:     false,
		TokenSymbol: "",
		ExecName:    types.ExecName(pt.ParaX),
	}
	tx, err := pt.CreateRawAssetTransferTx(&param)

	return tx, err
}

func createCrossParaTempTx(to string, amount int64) (*types.Transaction, error) {
	param := types.CreateTx{
		To:          string(to),
		Amount:      amount,
		Fee:         0,
		Note:        []byte("test asset transfer"),
		IsWithdraw:  false,
		IsToken:     false,
		TokenSymbol: "",
		ExecName:    Title2 + pt.ParaX,
	}
	tx, err := pt.CreateRawAssetTransferTx(&param)

	return tx, err
}

func createTxsGroup(txs []*types.Transaction) ([]*types.Transaction, error) {

	group, err := types.CreateTxGroup(txs)
	if err != nil {
		return nil, err
	}
	err = group.Check(0, types.GInt("MinFee"), types.GInt("MaxFee"))
	if err != nil {
		return nil, err
	}
	return group.Txs, nil
}

func TestGetBlockHashForkHeightOnMainChain(t *testing.T) {
	para := new(client)
	grpcClient := &typesmocks.Chain33Client{}
	grpcClient.On("GetFork", mock.Anything, &types.ReqKey{Key: []byte("ForkBlockHash")}).Return(&types.Int64{Data: 1}, errors.New("err")).Once()
	para.grpcClient = grpcClient
	_, err := para.GetForkHeightOnMainChain("ForkBlockHash")
	assert.NotNil(t, err)
	grpcClient.On("GetFork", mock.Anything, &types.ReqKey{Key: []byte("ForkBlockHash")}).Return(&types.Int64{Data: 1}, nil).Once()
	ret, err := para.GetForkHeightOnMainChain("ForkBlockHash")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), ret)

}
