package injapi

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	accountsPB "github.com/InjectiveLabs/sdk-go/exchange/accounts_rpc/pb"
	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type PublicClient struct {
	AccountClient            accountsPB.InjectiveAccountsRPCClient
	SpotExchangeClient       spotExchangePB.InjectiveSpotExchangeRPCClient
	DerivativeExchangeClient derivativeExchangePB.InjectiveDerivativeExchangeRPCClient
}

// sentry 0~3
func NewPublicClient(sentryNum int) *PublicClient {

	exchangeGRPC := fmt.Sprintf("tcp://sentry%d.injective.network:9910", sentryNum)

	// bohao exchange client init
	exchangeWaitCtx, cancelWait := context.WithTimeout(context.Background(), time.Minute)
	exchangeConn, err := grpcDialEndpoint(exchangeGRPC)
	if err != nil {
		log.Fatal(err, "endpoint:", exchangeGRPC, "failed to connect to API, is injective-exchange running?")
	}
	waitForService(exchangeWaitCtx, exchangeConn)
	cancelWait()

	// set up api clients
	accountsClient := accountsPB.NewInjectiveAccountsRPCClient(exchangeConn)
	spotExchangeClient := spotExchangePB.NewInjectiveSpotExchangeRPCClient(exchangeConn)
	derivativeExchangeClient := derivativeExchangePB.NewInjectiveDerivativeExchangeRPCClient(exchangeConn)
	return &PublicClient{
		AccountClient:            accountsClient,
		SpotExchangeClient:       spotExchangeClient,
		DerivativeExchangeClient: derivativeExchangeClient,
	}
}

func grpcDialEndpoint(protoAddr string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(protoAddr, grpc.WithInsecure(), grpc.WithContextDialer(dialerFunc))
	if err != nil {
		err := errors.Wrapf(err, "failed to connect to the gRPC: %s", protoAddr)
		return nil, err
	}

	return conn, nil
}

func waitForService(ctx context.Context, conn *grpc.ClientConn) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("Service wait timed out. Please run injective exchange service:\n\nmake install && injective-exchange")
		default:
			state := conn.GetState()

			if state != connectivity.Ready {
				time.Sleep(time.Second)
				continue
			}

			return nil
		}
	}
}

func DefaultSubaccount(acc cosmtypes.AccAddress) common.Hash {
	return common.BytesToHash(common.RightPadBytes(acc.Bytes(), 32))
}

func dialerFunc(ctx context.Context, protoAddr string) (net.Conn, error) {
	proto, address := protocolAndAddress(protoAddr)
	conn, err := net.Dial(proto, address)
	return conn, err
}

func protocolAndAddress(listenAddr string) (string, string) {
	protocol, address := "tcp", listenAddr
	parts := strings.SplitN(address, "://", 2)
	if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	}
	return protocol, address
}
