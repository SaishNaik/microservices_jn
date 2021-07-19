package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	protos "github.com/SaishNaik/microservices_jn/currency/protos/currency" 
	rates "github.com/SaishNaik/microservices_jn/currency/data" 
)

type Currency struct{
	rates *rates.ExchangeRates
	log hclog.Logger
}


func NewCurrency(r *rates.ExchangeRates,l hclog.Logger) (*Currency){
	return &Currency{r,l}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
	c.log.Info("coming to currency")
	c.log.Info("Handle GetRate","base",rr.GetBase(),"destination",rr.GetDestination())
	
	rate,err := c.rates.GetRate(rr.GetBase().String(),rr.GetDestination().String());
	if err != nil{
		return nil,err
	}
	return &protos.RateResponse{Rate:rate},nil
}