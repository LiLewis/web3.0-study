package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	// for demo

	"init_project/store"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// 合约地址
	contractAddr = "0x45678B5273144509a8fd003aaB75C575bf58c426"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/my_url_here")
	if err != nil {
		log.Fatal(err)
	}

	// 创建合约实例
	storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("my_private_key_here")
	if err != nil {
		log.Fatal(err)
	}

	/* -------------------------- 创建合约开始 -------------------------- */
	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// chainId, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(0)     // in wei
	// auth.GasLimit = uint64(300000) // in units
	// auth.GasPrice = gasPrice

	// input := "1.0"
	// address, tx, instance, err := store.DeployStore(auth, client, input)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(address.Hex())   //0x45678B5273144509a8fd003aaB75C575bf58c426
	// fmt.Println(tx.Hash().Hex()) //0x034ef2de6db7e0d59bbc57d61bb969c0785e0a1fd214d442cd3f0721a6c7ce74

	// _ = instance
	/* -------------------------- 创建合约结束 -------------------------- */

	/* -------------------------- 调用合约开始 -------------------------- */
	// 准备数据
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("lewis_save_key"))
	copy(value[:], []byte("lewis_save_value11111"))

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex()) //0x0629dee762d1adea60fab142c25ca50bfd7553afb0b234b9bd8036d55ac70cd9

	callOpt := &bind.CallOpts{Context: context.Background()}
	valueInContract, err := storeContract.Items(callOpt, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("is value saving in contract equals to origin value:", valueInContract == value) //true
	//sepolia链上已经成功执行:https://sepolia.etherscan.io/address/0x45678b5273144509a8fd003aab75c575bf58c426
	/* -------------------------- 调用合约结束 -------------------------- */
}
