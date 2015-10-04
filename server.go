package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"net/rpc"
	"strconv"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/rpc/jsonrpc"
	"time"
)
type Args struct {
	stockSymbolAndPercentage string	`json:"stockSymbolAndPercentage"`
	budget float64	`json:"budget"`
}

type StockResult struct {
	TradeId uint32 `json:"tradeid"`
	Stocks string `json:"Stocks"`
	unvestedAmount float64  `json:"unvestedAmount"`
}

type StockData struct {
    List struct{
	Meta struct{
		Count int	`json:"count"`
   		Start int 	`json:"start"`
		Type string `json:"type"`
	} `json:"meta"`
	Resources []struct{
		Resource struct{
			Classname string `json:"classname"`
			Fields struct{
				Name string    `json:"name"`
				Price string   `json:"price"`
				Symbol string  `json:"symbol"`
				Ts string      `json:"ts"`
				Type string    `json:"type"`
				UTCtime string `json:"utctime"`
				Volume string  `json:"volume"`
			}`json:"fields"`
		}`json:"resource"`
	}`json:"resources"`
    }`json:"list"`
}

type Arith int

var stockData StockData
var stockResult StockResult

func (t *Arith) GetPortfolioDetails(transId int64 ,reply *StockResult) error {
	fmt.Print(transId)
	fmt.Print(stockResult.TradeId)
	return nil
}

func (t *Arith) GetStockDetails(args *Args, reply *StockResult) error {
		stockAndPercent := args.stockSymbolAndPercentage
		
		stockAndPercent = strings.Replace(stockAndPercent,":",",",strings.Count(stockAndPercent,":"))
		stockAndPercent = strings.Replace(stockAndPercent,"%","",strings.Count(stockAndPercent,"%"))
		list := strings.Split(stockAndPercent,",")

		stockSymbol := ""
		stockPercent := ""
		stocksBudget := args.budget


		for i :=0; i<len(list); i++ {
			if i%2 == 0 {
				stockSymbol = stockSymbol+list[i]+","
			} else {
				stockPercent = stockPercent + list[i] + ","
			}
		}

		fmt.Println(stockSymbol)
		fmt.Println(stockPercent)

		ComputeStock(stockSymbol,stockPercent,stocksBudget)
		*reply = stockResult
		return nil


		
}

const(
	timeout = time.Duration(time.Second*100)
)

func ComputeStock(stockSymbol string, stockPercent string, stocksBudget float64) {
	
	client := http.Client{Timeout: timeout}
	httpurl := fmt.Sprintf("http://finance.yahoo.com/webservice/v1/symbols/"+stockSymbol+"/quote?format=json")

	req, err := client.Get(httpurl)

	if err != nil {
		fmt.Printf("Unable to fetch stocks: %s", err)
		os.Exit(1)
	}
	defer req.Body.Close()
	contents, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Unable to access json response : %s", err)
		os.Exit(1)
	}

	err = json.Unmarshal(contents, &stockData)
	if err != nil {
		fmt.Errorf("Stock cant be parsed")
	}

	var stockPriceString = ""

	numberOfStocks:= stockData.List.Meta.Count
	for i :=0; i<numberOfStocks;i++ {
		stockPriceString = stockPriceString + stockData.List.Resources[i].Resource.Fields.Price + ","
	}


	stockPrice := strings.Split(stockPriceString,",")
	percent := strings.Split(stockPercent,",")
	stocks := strings.Split(stockSymbol,",")


	var stockPriceF64 float64
	var percentF64 float64
	var quantity int
	var quantityStr = ""
	var sum float64
	var amountRemaining float64
	count := len(stocks) 
	for i :=0;i<count;i++ {
		stockPriceF64,_ = strconv.ParseFloat(stockPrice[i],64)
		percentF64,_ = strconv.ParseFloat(percent[i],64)
		quantity = int((stocksBudget*percentF64)/(100.00*stockPriceF64))
		sum = sum + (float64(quantity)*stockPriceF64)
		quantityStr = quantityStr + strconv.Itoa(quantity) + ","
	}

	stocksNumber := strings.Split(quantityStr,",")
	stockDetails := ""

	if(sum < stocksBudget) {
		for i :=0; i<count-1;i++ {
			stockPriceF64,_ = strconv.ParseFloat(stockPrice[i],64)
			stockDetails = stockDetails + stocks[i]+":"+stocksNumber[i]+"$"+stockPrice[i]
			fmt.Print(stockDetails) 
			if(i!=count-2) {
				stockDetails += ","
			}	
		}
		amountRemaining = stocksBudget - sum
		fmt.Printf("\nUnvested Amount: $%.2f\n",amountRemaining)
		sinceTime := time.Unix(0,11/11/2014)
		stockResult.TradeId = uint32(time.Since(sinceTime))
		stockResult.Stocks = stockDetails	
		stockResult.unvestedAmount = amountRemaining

	} else {
		fmt.Print("Amount exceeds the budget")
	}

	

}

func main() {
	fmt.Printf("I am here")
    arith := new(Arith)
    rpc.Register(arith)

    tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        jsonrpc.ServeConn(conn)
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}