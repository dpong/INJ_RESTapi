package injapi

import (
	"context"
	"fmt"
	"sync"

	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type INJTradeData struct {
	Price     string
	Qty       string
	Side      string
	OrderType string
	Timestamp string
}

type StreamMarketTradesBranch struct {
	client       *PublicClient
	cancel       *context.CancelFunc
	marketId     string
	tradeChan    chan INJTradeData
	tradesBranch struct {
		Trades []INJTradeData
		sync.Mutex
	}
	logger *logrus.Logger
}

// use default const type , ex: BTCUSDTperpId
func (GC *PublicClient) DerivativeTradeStream(marketId string, logger *logrus.Logger) *StreamMarketTradesBranch {
	o := new(StreamMarketTradesBranch)
	ctx, cancel := context.WithCancel(context.Background())
	o.client = GC
	o.cancel = &cancel
	o.marketId = marketId
	o.tradeChan = make(chan INJTradeData, 100)
	o.logger = logger
	go o.maintainSession(ctx)
	return o
}

func (o *StreamMarketTradesBranch) GetTrades() []INJTradeData {
	o.tradesBranch.Lock()
	defer o.tradesBranch.Unlock()
	trades := o.tradesBranch.Trades
	o.tradesBranch.Trades = []INJTradeData{}
	return trades
}

func (o *StreamMarketTradesBranch) Close() {
	(*o.cancel)()
	o.tradesBranch.Lock()
	defer o.tradesBranch.Unlock()
	o.tradesBranch.Trades = []INJTradeData{}
}

func (o *StreamMarketTradesBranch) maintainSession(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := o.maintain(ctx); err == nil {
				return
			} else {
				o.logger.Warningf("reconnect INJ %s trade stream with err: %s\n", DerivativeIdToSymbol[o.marketId], err.Error())
			}
		}
	}
}

func (o *StreamMarketTradesBranch) maintain(ctx context.Context) error {
	client, err := o.client.DerivativeExchangeClient.StreamTrades(ctx, &derivativeExchangePB.StreamTradesRequest{
		MarketId: o.marketId,
	})
	if err != nil {
		return err
	}
	o.logger.Infof("Connected derivative %s trade stream.", DerivativeIdToSymbol[o.marketId])
	for {
		msg, err := client.Recv()
		if err != nil {
			return err
		}
		if msg.GetTrade().GetTradeExecutionType() == "limitMatchNewOrder" || msg.GetTrade().GetTradeExecutionType() == "market" {
			o.insertTradeData(msg)
		}
	}
}

func (o *StreamMarketTradesBranch) insertTradeData(msg *derivativeExchangePB.StreamTradesResponse) {
	o.tradesBranch.Lock()
	defer o.tradesBranch.Unlock()
	p, _ := decimal.NewFromString(msg.Trade.GetPositionDelta().GetExecutionPrice())
	p = p.Div(decimal.New(1, 6))
	o.tradesBranch.Trades = append(o.tradesBranch.Trades, INJTradeData{
		Price:     p.String(),
		Qty:       msg.Trade.GetPositionDelta().GetExecutionQuantity(),
		Side:      msg.Trade.GetPositionDelta().GetTradeDirection(),
		OrderType: msg.GetTrade().GetTradeExecutionType(),
		Timestamp: fmt.Sprint(msg.GetTimestamp()),
	})
}
