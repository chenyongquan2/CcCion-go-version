package main

import (
	"CcCoin-go-version/internal/blockchain"
	"CcCoin-go-version/internal/server"
	"log"
	"net/http"
)

func main() {
	difficulty := 3

	blockchain := blockchain.NewBlockchain(difficulty)
	server := server.NewBlockchainServer(blockchain)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
