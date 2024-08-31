package main

import (
	"fmt"
	"testing"
)

func TestBlockChain(t *testing.T) {
	difficulty := 5

	myChain := NewBlockchain(difficulty)

	// 生成两个交易者身份的密钥对，也就是对应了钱包地址
	senderPrivateKey, senderPublicKey := generateKeyPair()
	_, receiverPublicKey := generateKeyPair()

	//公钥作为钱包的地址，标记转账时哪个钱包地址->另外一个钱包地址

	t1 := Transaction{to: senderPublicKey, from: receiverPublicKey, amount: 100}
	t2 := Transaction{to: senderPublicKey, from: receiverPublicKey, amount: 99}

	//使用发送者的密钥对里(其实只用到了私钥)来进行签名
	t1.Sign(senderPrivateKey)
	t2.Sign(senderPrivateKey)

	//尝试添加交易记录到chain的交易池子transactionPool里，等待"挖出来"的block来保存这些交易记录
	err := myChain.addTransction2Pool(t1)
	if err != nil {
		t.Errorf("Failed to add transaction to pool: %v", err)
	}

	err = myChain.addTransction2Pool(t2)
	if err != nil {
		t.Errorf("Failed to add transaction to pool: %v", err)
	}

	//准备矿工的身份
	_, minerPublicKey := generateKeyPair()

	//挖矿
	fmt.Println("正在挖矿...")
	myChain.mineTransctionFromPool(minerPublicKey)
	fmt.Println("挖完矿了")

	fmt.Println("hh")
}
