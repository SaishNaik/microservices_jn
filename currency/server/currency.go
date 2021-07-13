package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	protos "github.com/SaishNaik/microservices_jn/currency/protos/currency" 
)

type Currency struct{
	log hclog.Logger
}


func NewCurrency(l hclog.Logger) (*Currency){
	return &Currency{log:l}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
	c.log.Info("Handle GetRate","base",rr.GetBase(),"destination",rr.GetDestination())
	
	return &protos.RateResponse{Rate:0.5},nil
}