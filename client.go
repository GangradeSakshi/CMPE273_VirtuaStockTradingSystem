package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"os"
	"encoding/json"
)

type Args struct {
	StockSymbolAndPercentage string	`json:"StockSymbolAndPercentage"`
	Budget float64	`json:"Budget"`
}

type StockResult struct {
	TradeId uint32 `json:"tradeid"`
	Stocks string `json:"Stocks"`
	UnvestedAmount float64  `json:"UnvestedAmount"`
}

type Porfolio struct {
	tradeId string
}

var args Args

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		log.Fatal(1)
	}

	service := os.Args[1]
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("Error Dialing:", err)
	}

	var option string
	fmt.Println("Enter buy or Porfolio")
	fmt.Scanln(&option)

	if option == "buy" {
		fmt.Println ("Getting stock details")
		contents := []byte(os.Args[2])

		err = json.Unmarshal(contents, &args)

		var reply StockResult
		err = client.Call("Arith.GetStockDetails", args, &reply)
		if err != nil {
		log.Fatal("arith error:", err)
		}
	}
}
