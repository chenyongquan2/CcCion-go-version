package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func generateKeyPair() (privateKey, publicKey string) {
	// 使用椭圆曲线加密算法生成密钥对
	privateKeyECDSA, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// 将私钥转换为十六进制字符串
	privateKeyBytes := privateKeyECDSA.D.Bytes()
	privateKey = hex.EncodeToString(privateKeyBytes)

	// 将公钥转换为十六进制字符串
	publicKeyBytes := elliptic.Marshal(privateKeyECDSA.PublicKey.Curve, privateKeyECDSA.PublicKey.X, privateKeyECDSA.PublicKey.Y)
	publicKey = hex.EncodeToString(publicKeyBytes)

	return privateKey, publicKey
}
func main() {
	difficulty := 5

	myChain := NewBlockchain(difficulty)

	// 生成两个交易者身份的密钥对，也就是对应了钱包地址
	_, senderPublicKey := generateKeyPair()
	_, receiverPublicKey := generateKeyPair()

	//公钥作为钱包的地址，标记转账时哪个钱包地址->另外一个钱包地址

	t1 := Transaction{to: senderPublicKey, from: receiverPublicKey, amount: 100}
	t2 := Transaction{to: senderPublicKey, from: receiverPublicKey, amount: 99}

	//使用发送者的密钥对里(其实只用到了私钥)来进行签名
	// t1.Sign(senderPrivateKey)
	// t2.Sign(senderPrivateKey)

	//尝试添加交易记录到chain的交易池子transactionPool里，等待"挖出来"的block来保存这些交易记录
	myChain.addTransction2Pool(t1)
	myChain.addTransction2Pool(t2)

	//准备矿工的身份
	_, minerPublicKey := generateKeyPair()

	//挖矿
	fmt.Println("正在挖矿...")
	myChain.mineTransctionFromPool(minerPublicKey)
	fmt.Println("挖完矿了")

	fmt.Println("hh")
}
