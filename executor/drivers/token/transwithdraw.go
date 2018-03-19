package token

import (
	"code.aliyun.com/chain33/chain33/account"
	dbm "code.aliyun.com/chain33/chain33/common/db"
	"code.aliyun.com/chain33/chain33/executor/drivers"
	"code.aliyun.com/chain33/chain33/types"
	"fmt"
)

func (t *Token) ExecTransWithdraw(accountDB *account.AccountDB, tx *types.Transaction, action *types.TokenAction) (*types.Receipt, error) {
	if (action.Ty == types.ActionTransfer) && action.GetTransfer() != nil {
		transfer := action.GetTransfer()
		from := account.From(tx).String()
		//to 是 execs 合约地址
		if drivers.IsDriverAddress(tx.To) {
			return accountDB.TransferToExec(from, tx.To, transfer.Amount)
		}
		return accountDB.Transfer(from, tx.To, transfer.Amount)
	} else if (action.Ty == types.ActionWithdraw) && action.GetWithdraw() != nil {
		withdraw := action.GetWithdraw()
		from := account.PubKeyToAddress(tx.Signature.Pubkey).String()
		//to 是 execs 合约地址
		if drivers.IsDriverAddress(tx.To) {
			return accountDB.TransferWithdraw(from, tx.To, withdraw.Amount)
		}
		return nil, types.ErrActionNotSupport
	} else if (action.Ty == types.ActionGenesis) && action.GetGenesis() != nil {
		genesis := action.GetGenesis()
		if t.GetHeight() == 0 {
			if drivers.IsDriverAddress(tx.To) {
				return accountDB.GenesisInitExec(genesis.ReturnAddress, genesis.Amount, tx.To)
			}
			return accountDB.GenesisInit(tx.To, genesis.Amount)
		} else {
			return nil, types.ErrReRunGenesis
		}
	} else {
		return nil, types.ErrActionNotSupport
	}
}

//0: all tx
//1: from tx
//2: to tx

func (t *Token) ExecLocalTransWithdraw(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecLocal(tx, receipt, index)
	if err != nil {
		return nil, err
	}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//执行成功
	var action types.CoinsAction
	err = types.Decode(tx.GetPayload(), &action)
	if err != nil {
		panic(err)
	}
	var kv *types.KeyValue
	if action.Ty == types.ActionTransfer && action.GetTransfer() != nil {
		transfer := action.GetTransfer()
		kv, err = updateAddrReciver(t.GetLocalDB(), transfer.Cointoken, tx.To, transfer.Amount, true)
	} else if action.Ty == types.ActionWithdraw && action.GetWithdraw() != nil {
		withdraw := action.GetWithdraw()
		from := account.PubKeyToAddress(tx.Signature.Pubkey).String()
		kv, err = updateAddrReciver(t.GetLocalDB(), withdraw.Cointoken, from, withdraw.Amount, true)
	} else if action.Ty == types.ActionGenesis && action.GetGenesis() != nil {
		gen := action.GetGenesis()
		kv, err = updateAddrReciver(t.GetLocalDB(), "token", tx.To, gen.Amount, true)
	}
	if err != nil {
		return set, nil
	}
	if kv != nil {
		set.KV = append(set.KV, kv)
	}
	return set, nil
}

func (t *Token) ExecDelLocalLocalTransWithdraw(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecDelLocal(tx, receipt, index)
	if err != nil {
		return nil, err
	}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//执行成功
	var action types.CoinsAction
	err = types.Decode(tx.GetPayload(), &action)
	if err != nil {
		panic(err)
	}
	var kv *types.KeyValue
	if action.Ty == types.ActionTransfer && action.GetTransfer() != nil {
		transfer := action.GetTransfer()
		kv, err = updateAddrReciver(t.GetLocalDB(), transfer.Cointoken, tx.To, transfer.Amount, false)
	} else if action.Ty == types.ActionWithdraw && action.GetWithdraw() != nil {
		withdraw := action.GetWithdraw()
		from := account.PubKeyToAddress(tx.Signature.Pubkey).String()
		kv, err = updateAddrReciver(t.GetLocalDB(), withdraw.Cointoken, from, withdraw.Amount, false)
	}
	if err != nil {
		return set, nil
	}
	if kv != nil {
		set.KV = append(set.KV, kv)
	}
	return set, nil
}

//存储地址上收币的信息
func CalcAddrKey(token string, addr string) []byte {
	return []byte(fmt.Sprintf("Token:%s-Addr:%s", token, addr))
}

type AddrRecv struct {
	addr   string
	amount int64
}

func geAddrReciverKV(token string, addr string, reciverAmount int64) *types.KeyValue {
	reciver := &types.Int64{reciverAmount}
	amountbytes := types.Encode(reciver)
	kv := &types.KeyValue{CalcAddrKey(token, addr), amountbytes}
	return kv
}

func getAddrReciver(db dbm.KVDB, token string, addr string) (int64, error) {
	reciver := types.Int64{}
	addrReciver, err := db.Get(CalcAddrKey(token, addr))
	if err != nil && err != types.ErrNotFound {
		return 0, err
	}
	if len(addrReciver) == 0 {
		return 0, nil
	}
	err = types.Decode(addrReciver, &reciver)
	if err != nil {
		return 0, err
	}
	return reciver.Data, nil
}

func setAddrReciver(db dbm.KVDB, token string, addr string, reciverAmount int64) error {
	kv := geAddrReciverKV(addr, token, reciverAmount)
	return db.Set(kv.Key, kv.Value)
}

func updateAddrReciver(cachedb dbm.KVDB, token string, addr string, amount int64, isadd bool) (*types.KeyValue, error) {
	recv, err := getAddrReciver(cachedb, token, addr)
	if err != nil && err != types.ErrNotFound {
		return nil, err
	}
	if isadd {
		recv += amount
	} else {
		recv -= amount
	}
	setAddrReciver(cachedb, token, addr, recv)
	//keyvalue
	return geAddrReciverKV(token, addr, recv), nil
}
