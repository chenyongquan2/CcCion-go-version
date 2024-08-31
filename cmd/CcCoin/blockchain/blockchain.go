// package blockchain
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	MinerRewardFromAddress = ""
)

type Transaction struct {
	//from和to表示交易者的钱包地址，amount表示交易的金额
	to        string
	from      string
	amount    float64
	signature string
}

func (t *Transaction) computeHash() []byte {
	data := fmt.Sprintf("%v%v%v", t.from, t.to, t.amount)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// Sign 使用私钥对交易数据的哈希值进行签名
func (t *Transaction) Sign(privateKey string) error {
	var err error
	t.signature, err = signMessage(privateKey, string(t.computeHash()))
	return err
}

func (t *Transaction) isValid() (res bool) {
	//当this.from === ''，说明该转账是由区块链发起的矿工奖励，无需校验签名的合法性
	if t.from == MinerRewardFromAddress {
		return true
	}

	var err error
	res, err = verifySignature(t.from, string(t.computeHash()), t.signature)
	if err != nil {
		fmt.Println("verify signature failed,err:", err)
		return false
	}
	return
}

// 区块，用来存储交易信息
type Block struct {
	transactions []Transaction //这个区块所存储的交易信息
	prevHash     string        //前一个区块的hash
	hash         string        //hash是一个区块的指纹
	nonce        int           //随机数
	timestamp    uint64        //时间戳
}

func NewBlock(transactions []Transaction, prevHash string) Block {
	block := Block{
		transactions: transactions,
		prevHash:     prevHash,
	}
	block.hash = hex.EncodeToString(block.computeHash())
	block.nonce = 1
	block.timestamp = uint64(time.Now().Unix())
	return block
}

func (block *Block) computeHash() []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(block.prevHash),
			[]byte(fmt.Sprintf("%v", block.transactions)),
			[]byte(strconv.FormatUint(block.timestamp, 10)),
			[]byte(strconv.Itoa(block.nonce)),
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (block *Block) getAnswer(difficulty int) string {
	//开头前n位为0的hash
	return strings.Repeat("0", difficulty)
}

func (block *Block) validateBlockTransations() bool {
	for _, t := range block.transactions {
		if !t.isValid() {
			fmt.Println("invalid transaction found in transations, 发现异常交易")
			return false
		}
	}
	return true
}

// 计算符号区块难度要求的hash
// 为什么需要引入难度要求?为了控制每10min会有一个区块被挖矿挖出来，需要动态调整这个难度要求
func (block *Block) mine(difficulty int) {
	//开挖之前，应该要检查一下即将要挖来存储的transctions的合法性,避免浪费算力
	bOk := block.validateBlockTransations()
	if !bOk {
		fmt.Println("invalid transaction found in transations, stop mining!")
		return
	}

	ans := block.getAnswer(difficulty)
	for {
		hashRes := block.computeHash()
		//fmt.Println(hashRes)
		if hex.EncodeToString(hashRes)[:difficulty] != ans {
			//改变随机数，继续尝试
			block.nonce++
		} else {
			block.hash = hex.EncodeToString(hashRes)
			fmt.Printf("挖矿结束, nonce:%d,difficulty:%d,hash:%s\n", block.nonce, difficulty, block.hash)
			break
		}
	}

}

// 区块的链表
// 区块链是一个transations转账记录的池子，需要一个miner reword
type Blockchain struct {
	blocks          []Block       //保存的所有区块
	difficulty      int           //复杂度
	transationsPool []Transaction //交易池子
	minerReward     float64       //矿工奖励
}

func NewBlockchain(difficulty int) Blockchain {
	blockchain := Blockchain{
		blocks:          []Block{},
		difficulty:      difficulty,
		transationsPool: []Transaction{},
		minerReward:     50,
	}
	//每当这个puzzle被发出来后，矿工会从transctionPool池子里面去一部分收益最高的transction(因为每一个block的大小是有限的，能容纳的transction数目是有限的)，以这些transction为基础去新建这个block
	//也就意味着这个block的新建应该是发生在链上的，发生在哪一个步骤呢，发生在挖transction这个操作里
	blockchain.blocks = append(blockchain.blocks, blockchain.bingBang())
	return blockchain
}

// 生成祖先区块/创世区块(Genesis Block)
// 创世区块是区块链中第一个被创建的区块
// 隐喻了区块链网络的诞生,就像宇宙大爆炸(Big Bang)一样,创世区块标志着区块链网络的开始。
func (blockchain *Blockchain) bingBang() Block {
	genesisBlock := NewBlock([]Transaction{}, "0")
	return genesisBlock
}

func (blockchain *Blockchain) getLatestBlock() Block {
	return blockchain.blocks[len(blockchain.blocks)-1]
}

// 添加待存储的transction到transction pool里面，供后续挖出来的block来存储这些transction交易记录
func (blockchain *Blockchain) addTransction2Pool(transaction Transaction) error {
	// 添加transaction到transationsPool之前，先校验一下transation的合法性
	if !transaction.isValid() {
		return errors.New("invalid transaction,reject it")
	}
	blockchain.transationsPool = append(blockchain.transationsPool, transaction)
	fmt.Println("valid transaction has been pushed to transationsPool")
	return nil
}

// 从chain的待存储的transationsPool里面挑选收益最高的transations来存储到新生成的block
// 也就是说生成block的过程应该是chain来负责了，而不是像上面方法一样是外面传进来的
func (blockchain *Blockchain) mineTransctionFromPool(minerRewardAddress string) {
	///生成矿工奖励的transction,放到transationsPool里面
	minerRewardTransction := Transaction{
		from:   MinerRewardFromAddress,
		to:     minerRewardAddress,
		amount: blockchain.minerReward,
	}
	err := blockchain.addTransction2Pool(minerRewardTransction)
	if err != nil {
		fmt.Println("add minerRewardTransction to transationsPool failed")
		return
	}

	//从transationsPool挑选收益最高的transations来存储到新生成的block
	newBlock := NewBlock(blockchain.transationsPool, blockchain.getLatestBlock().hash)
	newBlock.mine(blockchain.difficulty)
	blockchain.blocks = append(blockchain.blocks, newBlock)
	blockchain.transationsPool = []Transaction{} //reset transationsPool
}

// 验证区块的合法性
func (blockchain *Blockchain) isValidChain() bool {
	if len(blockchain.blocks) == 1 {
		//通过区块的hash值，验证内容和hash值有无被篡改
		if blockchain.blocks[0].hash != string(blockchain.blocks[0].computeHash()) {
			fmt.Println("祖先区块被篡改了!")
			return false
		}
		return true
	}

	for i := 1; i < len(blockchain.blocks); i++ {
		block := blockchain.blocks[i]
		//检验当前数据是否有无被篡改
		if block.hash != string(block.computeHash()) {
			fmt.Printf("区块 %d 被篡改了!\n", i)
			return false
		}
		//通过prevHash来判断是否断链
		prevBlockHash := blockchain.blocks[i-1].hash
		if block.prevHash != prevBlockHash {
			fmt.Printf("区块 %d 断联了!\n", i)
			return false
		}

		//还需要验证 链里面的每一个区块是否被篡改了
		if !block.validateBlockTransations() {
			fmt.Printf("发现链里面有非法交易,异常block idx: %d\n", i)
			return false
		}
	}

	return true
}
