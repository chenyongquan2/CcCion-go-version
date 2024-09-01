package server

import (
	"CcCoin-go-version/internal/blockchain"
	"encoding/json"
	"errors"
	"net/http"
)

type BlockchainServer struct {
	blockchain blockchain.Blockchain
	http.Handler
}

func NewBlockchainServer(blockchain blockchain.Blockchain) *BlockchainServer {
	p := new(BlockchainServer)
	p.blockchain = blockchain

	router := http.NewServeMux()
	router.Handle("/transction/", http.HandlerFunc(p.transactionHandler))
	router.Handle("/mine/", http.HandlerFunc(p.mineHandler))

	p.Handler = router
	return p
}

func (p *BlockchainServer) addTransaction(w http.ResponseWriter, r *http.Request) error {
	//Todo:理论上SenderPrivateKey不应该每次都通过网络传递来的，应该存在server的数据库，这里为了简便，先这么搞着
	// 解析交易数据
	var txData struct {
		SenderPublicKey   string  `json:"SenderPublicKey"`
		SenderPrivateKey  string  `json:"SenderPrivateKey"`
		ReceiverPublicKey string  `json:"ReceiverPublicKey"`
		Amount            float64 `json:"Amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&txData)
	if err != nil {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return err
	}
	if txData.SenderPublicKey == "" || txData.SenderPrivateKey == "" || txData.ReceiverPublicKey == "" || txData.Amount == 0 {
		err = errors.New("missing required transaction fields")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	// 使用NewTransaction方法创建Transaction对象
	tx, err := blockchain.NewTransaction(txData.SenderPublicKey, txData.SenderPrivateKey, txData.ReceiverPublicKey, txData.Amount)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusBadRequest)
		return err
	}

	// 验证交易
	if !tx.IsValid() {
		http.Error(w, "Invalid transaction", http.StatusBadRequest)
		return errors.New("invalid transaction")
	}

	// 添加交易到交易池
	err = p.blockchain.AddTransction2Pool(tx)
	if err != nil {
		http.Error(w, "Failed to add transaction to pool", http.StatusInternalServerError)
		return err
	}

	// 返回成功响应
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction added successfully"})

	return nil
}

func (p *BlockchainServer) startMineTask(w http.ResponseWriter, r *http.Request) error {
	// 解析请求数据
	var mineData struct {
		MinerPublicKey string `json:"MinerPublicKey"`
	}
	err := json.NewDecoder(r.Body).Decode(&mineData)
	if err != nil {
		http.Error(w, "Invalid mine data", http.StatusBadRequest)
		return err
	}

	err = p.blockchain.MineTransctionFromPool(mineData.MinerPublicKey)
	if err != nil {
		http.Error(w, "mine data failed", http.StatusBadRequest)
		return err
	}

	return nil
}

func (p *BlockchainServer) transactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := p.addTransaction(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		// 处理获取交易的逻辑
		// 例如：返回所有交易或特定交易的信息
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (p *BlockchainServer) mineHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := p.startMineTask(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		// 处理获取交易的逻辑
		// 例如：返回所有交易或特定交易的信息
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
