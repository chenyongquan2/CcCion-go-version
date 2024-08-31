package test

import (
	"CcCoin-go-version/internal/blockchain"
	"CcCoin-go-version/internal/encryption"
	"fmt"
	"testing"
)

func TestBlockChain(t *testing.T) {
	difficulty := 3

	myChain := blockchain.NewBlockchain(difficulty)

	// 生成两个交易者身份的密钥对，也就是对应了钱包地址
	senderPrivateKey, senderPublicKey := encryption.GenerateKeyPair()
	_, receiverPublicKey := encryption.GenerateKeyPair()

	//公钥作为钱包的地址，标记转账时哪个钱包地址->另外一个钱包地址

	t1, err := blockchain.NewTransaction(senderPublicKey, senderPrivateKey, receiverPublicKey, 100)
	if err != nil {
		t.Errorf("NewTransaction failed err: %v", err)
	}

	t2, err := blockchain.NewTransaction(senderPublicKey, senderPrivateKey, receiverPublicKey, 99)
	if err != nil {
		t.Errorf("NewTransaction failed err: %v", err)
	}

	//尝试添加交易记录到chain的交易池子transactionPool里，等待"挖出来"的block来保存这些交易记录
	err = myChain.AddTransction2Pool(t1)
	if err != nil {
		t.Errorf("Failed to add transaction to pool: %v", err)
	}

	err = myChain.AddTransction2Pool(t2)
	if err != nil {
		t.Errorf("Failed to add transaction to pool: %v", err)
	}

	//准备矿工的身份
	_, minerPublicKey := encryption.GenerateKeyPair()

	//挖矿
	fmt.Println("正在挖矿...")
	myChain.MineTransctionFromPool(minerPublicKey)
	fmt.Println("挖完矿了")
}
