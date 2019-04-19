package plugin

import (
	_ "github.com/CPF999/cpfChain/plugin/consensus/init" //consensus init
	_ "github.com/CPF999/cpfChain/plugin/crypto/init"    //crypto init
	_ "github.com/CPF999/cpfChain/plugin/dapp/init"      //dapp init
	_ "github.com/CPF999/cpfChain/plugin/mempool/init"   //mempool init
	_ "github.com/CPF999/cpfChain/plugin/store/init"     //store init
)
