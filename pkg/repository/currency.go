package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
	"gitlab.com/sdce/exlib/exutil"
	"gitlab.com/sdce/exlib/mongo"
	pb "gitlab.com/sdce/protogo"
	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

//Currency ...
type Currency struct {
	C *mgo.Collection
	Q *mgo.Collection
}

//NewCurrency create and initialize the currency collection
func NewCurrency(db *mongo.Database) (*Currency, error) {
	currency := &Currency{
		C: db.CreateCollection(CollectionCurrency), // currency collection
		Q: db.CreateCollection(CollectionQuote),    // market quote collection
	}
	err := mongo.UniqueFields(context.Background(), currency.C, []string{
		"symbol", "name"})
	if err != nil {
		return nil, err
	}
	return currency, nil
}

//CreateCurrency ...
func (s *Currency) CreateCurrency(ctx context.Context, cr *pb.Currency) (*pb.UUID, error) {
	cId := exutil.NewUUID()
	cr.Id = cId
	_, err := s.C.InsertOne(ctx, cr)
	if err != nil {
		return nil, err
	}
	return cId, nil
}

//GetCurrency ...
func (s *Currency) GetCurrency(ctx context.Context, id *pb.UUID) (*pb.Currency, error) {
	cur := &pb.Currency{}
	err := s.C.FindOne(ctx, bson.M{"_id": id}).Decode(cur)
	return cur, err
}

//GetCurrencyBySymbol ...
func (s *Currency) GetCurrencyBySymbol(ctx context.Context, symbol string) (*pb.Currency, error) {
	cur := &pb.Currency{}
	err := s.C.FindOne(ctx, bson.M{"symbol": symbol}).Decode(cur)
	return cur, err
}

//UpdateCurrency ...
func (s *Currency) UpdateCurrency(ctx context.Context, req *pb.UpdateCurrencyRequest) (*pb.UUID, error) {
	if req == nil || req.GetCurrency() == nil {
		return nil, errors.New("currency object nil")
	}

	uobj, err := exutil.ApplyFieldMaskToBson(req.GetCurrency(), req.GetUpdateMask())
	if err != nil {
		return nil, err
	}
	err = s.C.FindOneAndUpdate(ctx, bson.M{"_id": req.Currency.Id}, bson.M{"$set": uobj}).Err()
	if err != nil {
		return nil, err
	}
	return req.Currency.Id, nil
}

//SearchCurrencies ...
func (s *Currency) SearchCurrencies(ctx context.Context, tl []pb.Currency_CurrencyType, sl []pb.Currency_CurrencyStatus, symbol string, name string) ([]*pb.Currency, error) {
	var currencies []*pb.Currency
	opts := &options.FindOptions{
		Sort: bson.M{"sortOrder": -1},
	}
	filter := bson.M{}
	if name != "" {
		filter["name.en"] = name
	}
	if symbol != "" {
		filter["symbol"] = symbol
	}
	if len(sl) != 0 {
		filter["status"] = bson.M{"$in": sl}
	}
	if len(tl) != 0 {
		filter["type"] =  bson.M{"$in": tl}
	}
	cur, err := s.C.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var currency pb.Currency
		err = cur.Decode(&currency)
		if err != nil {
			log.Info("Cannot decode currency")
			return nil, errors.New("Cannot decode currency")
		}
		currencies = append(currencies, &currency)
	}
	return currencies, err
}

////////Deprected: methods below are not confirmed, may be removed in the future.

//CreateCurrencyMeta ...
func (s *Currency) CreateCurrencyMeta(ctx context.Context, cr *pb.CurrencyMetaInfo) (*pb.UUID, error) {
	//todo
	return nil, nil
}

//GetCurrencyMeta ...
func (s *Currency) GetCurrencyMeta(ctx context.Context, id *pb.UUID) (*pb.CurrencyMetaInfo, error) {
	//todo
	return nil, nil
}

//UpdateCurrencyMeta ...
func (s *Currency) UpdateCurrencyMeta(ctx context.Context, cr *pb.CurrencyMetaInfo) (*pb.UUID, error) {
	//todo
	return nil, nil
}

func (s *Currency) SearchMarketQuotes(ctx context.Context, coins []string) (map[string]*pb.MarketQuote, error) {
	var feeds = make(map[string]*pb.MarketQuote)
	// use no filter for now
	filter := bson.M{}
	if coins != nil {
		filter["symbol"] = bson.M{"$in": coins}
	}
	cur, err := s.Q.Find(ctx, filter)
	if err != nil {
		log.Errorf("failed to find quotes from DB %v", err)
		return feeds, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var quote QuoteDBModel
		cur.Decode(&quote)
		feeds[quote.Symbol] = &pb.MarketQuote{
			Name:   quote.Name,
			Symbol: quote.Symbol,
			Time:   quote.Time,
			Price:  fmt.Sprintf("%.2f", quote.Quote.Price),
		}
	}
	return feeds, err
}
