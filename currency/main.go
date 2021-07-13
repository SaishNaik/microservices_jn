package main

import (
	"os"
	"net"
	protos "github.com/SaishNaik/microservices_jn/currency/protos/currency"
	"github.com/SaishNaik/microservices_jn/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// "fmt"
	"github.com/hashicorp/go-hclog"
)

func main(){
	log := hclog.Default()
	
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs,cs)
	reflection.Register(gs)

	l,err := net.Listen("tcp",":9092")
	if err != nil {
		log.Error("Unable to listen","error",err)
		os.Exit(1)
	} 

	gs.Serve(l)
	
}