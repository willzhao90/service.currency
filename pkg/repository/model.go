package repository

import (
	mgo "gitlab.com/sdce/exlib/mongo"
)

const (
	//DatabaseName ...
	DatabaseName = "exchange"
	//CollectionCurrency ...
	CollectionCurrency = "currency"
	CollectionQuote    = "market_feed"
)

type Repository struct {
	Name     string
	Currency *Currency
}

//CreateRepository ...
func CreateRepository(db *mgo.Database) (*Repository, error) {
	currency, err := NewCurrency(db)
	if err != nil {
		return nil, err
	}
	return &Repository{
		Name:     DatabaseName,
		Currency: currency,
	}, nil
}

//===============================================================================================================
// tobe moved to proto
type QuoteDetails struct {
	Price            float64 `json:"price" bson:"price"`
	Volume24H        float64 `json:"volume_24h" bson:"volume24h"`
	PercentChange1H  float64 `json:"percent_change_1h" bson:"percentChange1h"`
	PercentChange24H float64 `json:"percent_change_24h" bson:"percentChange24h"`
	PercentChange7D  float64 `json:"percent_change_7d" bson:"percentChange7d"`
	MarketCap        float64 `json:"market_cap" bson:"marketCap"`
	LastUpdated      string  `json:"last_updated" bson:"lastUpdated"`
}

type QuoteDBModel struct {
	Name    string        `json:"name" bson:"name"`
	Symbol  string        `json:"symbol" bson:"symbol"`
	Time    string        `json:"time" bson:"time"`
	Quote   *QuoteDetails `json:"quote" bson:"quoteDetails"`
	Created int64         `json:"created" bson:"created"`
}
