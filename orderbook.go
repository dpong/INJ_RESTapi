package injapi

import (
	"context"
	"fmt"

	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	"github.com/shopspring/decimal"
)

//askPrice, askQty, bidPrice, bidQty, timeStamp, error
func (GC *PublicClient) GetDerivativeToB(symbol string, quoteDecimals int32) (string, string, string, string, string, error) {
	orderbook, err := GC.DerivativeExchangeClient.Orderbook(context.Background(), &derivativeExchangePB.OrderbookRequest{
		MarketId: symbol,
	})
	if err != nil {
		return "", "", "", "", "", err
	}
	tobBid := orderbook.GetOrderbook().GetBuys()[0]
	tobAsk := orderbook.GetOrderbook().GetSells()[0]
	askPrice, _ := decimal.NewFromString(tobAsk.GetPrice())
	bidPrice, _ := decimal.NewFromString(tobBid.GetPrice())
	askPrice = askPrice.Div(decimal.New(1, quoteDecimals))
	bidPrice = bidPrice.Div(decimal.New(1, quoteDecimals))
	return askPrice.String(), tobAsk.GetQuantity(), bidPrice.String(), tobBid.GetQuantity(), fmt.Sprint(tobAsk.GetTimestamp()), nil
}
