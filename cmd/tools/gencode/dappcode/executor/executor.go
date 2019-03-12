// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/cmd/tools/gencode/base"
	"github.com/33cn/chain33/cmd/tools/types"
)

func init() {

	base.RegisterCodeFile(executorCodeFile{})
}


type executorCodeFile struct {
	base.DappCodeFile
}


func (c executorCodeFile)GetDirName() string {

	return "executor"
}


func (c executorCodeFile)GetFiles() map[string]string {

	return map[string]string{
		executorName: executorContent,
	}
}


func (c executorCodeFile) GetReplaceTags() []string {

	return []string{types.TagExecName, types.TagClassName}
}



var (

	executorName = "${EXECNAME}.go"
	executorContent =
`package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	ptypes "github.com/33cn/plugin/plugin/dapp/${EXECNAME}/types"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
)

var (
	ptylog = log.New("module", "execs.${EXECNAME}")
)

var driverName = ptypes.${CLASSNAME}X

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&${EXECNAME}{}))
}

func Init(name string, sub []byte) {
	drivers.Register(GetName(), new${CLASSNAME}, types.GetDappFork(driverName, "Enable"))
}

type ${EXECNAME} struct {
	drivers.DriverBase
}

func new${CLASSNAME}() drivers.Driver {
	t := &${EXECNAME}{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

func GetName() string {
	return new${CLASSNAME}().GetName()
}

func (*${EXECNAME}) GetDriverName() string {
	return driverName
}

`
)