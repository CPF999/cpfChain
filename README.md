[![pipeline status](https://api.travis-ci.org/bityuan/bityuan.svg?branch=master)](https://travis-ci.org/bityuan/bityuan/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bityuan/bityuan)](https://goreportcard.com/report/github.com/bityuan/bityuan)

# 基于 chain33 区块链开发 框架 开发的 CPF平行链 系统

#### 编译

```
git clone https://github.com/CPF999/cpfChain $GOPATH/src/github.com/CPF999/cpfChain
cd $GOPATH/src/github.com/CPF999/cpfChain
go build -i -o cpf
go build -i -o cpf-cli github.com/CPF999/cpfChain/cli
```

#### 运行

拷贝编译好的cpf, cpf-cli, cpfChain.toml这三个文件置于同一个文件夹下，执行：

```
./cpf -f cpfChain.toml
```


