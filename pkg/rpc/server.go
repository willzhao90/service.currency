package rpc

import (
	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"

	"fmt"

	pb "gitlab.com/sdce/protogo"
	//"gitlab.com/sdce/service/currency/pkg/config"

	"gitlab.com/sdce/service/currency/pkg/repository"
)

//Server ...
type Server struct {
	url      string
	currency *repository.Currency
}

//DoCreateCurrency ...
func (r *Server) DoCreateCurrency(ctx context.Context, req *pb.CreateCurrencyRequest) (*pb.CreateCurrencyResponse, error) {
	log.Debugln("creating new currency")

	if r == nil {
		log.Fatalln("RPCServer is null")
	}

	id, err := r.currency.CreateCurrency(ctx, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to create new currency in db: %v", err.Error())
	}

	log.Infoln("Currency created:" + id.String())
	res := &pb.CreateCurrencyResponse{
		Id: id,
	}
	return res, nil
}

//DoGetCurrency ...
func (r *Server) DoGetCurrency(ctx context.Context, req *pb.GetCurrencyRequest) (*pb.GetCurrencyResponse, error) {
	log.Debugln("get currency data")
	if r == nil {
		log.Fatalln("RPCServer nullptr")
	}
	c, err := r.currency.GetCurrency(ctx, req.Id)
	res := &pb.GetCurrencyResponse{
		Currency: c,
	}
	return res, err
}

//DoGetCurrency ...
func (r *Server) DoGetCurrencyBySymbol(ctx context.Context, req *pb.GetCurrencyBySymbolRequest) (*pb.GetCurrencyBySymbolResponse, error) {
	log.Debugln("get currency data")
	if r == nil {
		log.Fatalln("RPCServer nullptr")
	}
	c, err := r.currency.GetCurrencyBySymbol(ctx, req.Symbol)
	res := &pb.GetCurrencyBySymbolResponse{
		Currency: c,
	}
	return res, err
}

//DoUpdateCurrency ...
func (r *Server) DoUpdateCurrency(ctx context.Context, req *pb.UpdateCurrencyRequest) (*pb.UpdateCurrencyResponse, error) {
	log.Debugln("update currency data")
	if r == nil {
		log.Fatalln("RPCServer nullptr")
	}
	id, err := r.currency.UpdateCurrency(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &pb.UpdateCurrencyResponse{Id: id}
	return res, nil
}

//DoSearchCurrencies ...
func (r *Server) DoSearchCurrencies(ctx context.Context, req *pb.SearchCurrenciesRequest) (*pb.SearchCurrenciesResponse, error) {
	log.Debugln("Search currencies")
	if r == nil {
		log.Fatalln("RPCServer nullptr")
	}

	res, err := r.currency.SearchCurrencies(ctx, nil, req.Status, req.Symbol, req.Name)
	return &pb.SearchCurrenciesResponse{
		Currencies: res,
	}, err
}

func (r *Server) DoSearchMarketQuotes(ctx context.Context, req *pb.SearchMarketQuotesRequest) (*pb.SearchMarketQuotesResponse, error) {
	quotes, err := r.currency.SearchMarketQuotes(ctx, req.GetCurrencies())

	if err != nil {
		log.Errorf("failed to fetch market quote: %v", req.GetCurrencies())
		return nil, err
	}

	log.Errorf("market quotes received : %v", quotes)

	return &pb.SearchMarketQuotesResponse{Quotes: quotes}, nil
}

//Notice:
//////////Deprecated: meta methods below don't need to be implemented until meta structure is confirmed.

//DoCreateCurrencyMeta ...
func (r *Server) DoCreateCurrencyMeta(context.Context, *pb.CreateCurrencyMetaRequest) (*pb.CreateCurrencyMetaResponse, error) {
	log.Debugln("creating currency meta")
	return nil, nil
}

//DoGetCurrencyMeta ...
func (r *Server) DoGetCurrencyMeta(context.Context, *pb.GetCurrencyMetaRequest) (*pb.GetCurrencyMetaResponse, error) {
	log.Debugln("get currency meta")
	return nil, nil
}

//DoUpdateCurrencyMeta ...
func (r *Server) DoUpdateCurrencyMeta(context.Context, *pb.UpdateCurrencyMetaRequest) (*pb.UpdateCurrencyMetaResponse, error) {
	log.Debugln("update currency meta")
	return nil, nil
}

//NewRPCServer ...
//if config is nil, defPort will be used instead.
func NewRPCServer(config interface{}, defaultPort string, r *repository.Repository) *Server {
	return &Server{
		url:      defaultPort,
		currency: r.Currency,
	}
}
