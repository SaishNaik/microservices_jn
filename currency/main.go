package main

import (
	"os"
	"net"
	protos "github.com/SaishNaik/microservices_jn/currency/protos/currency"
	"github.com/SaishNaik/microservices_jn/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	rates "github.com/SaishNaik/microservices_jn/currency/data" 
	// "fmt"
	"github.com/hashicorp/go-hclog"
)

func main(){
	log := hclog.Default()
	
	gs := grpc.NewServer()

	rates,err := rates.NewRates(log)
	if err != nil{
		log.Error("Unable to generate rates","error",err)
		os.Exit(1)
	}
	cs := server.NewCurrency(rates,log)

	protos.RegisterCurrencyServer(gs,cs)
	reflection.Register(gs)

	l,err := net.Listen("tcp",":9092")
	if err != nil {
		log.Error("Unable to listen","error",err)
		os.Exit(1)
	} 

	gs.Serve(l)
	
}