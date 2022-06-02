module github.com/dpong/INJ_RESTapi

go 1.16

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/btcsuite/btcutil => github.com/btcsuite/btcutil v1.0.2

replace github.com/cosmos/cosmos-sdk => github.com/InjectiveLabs/cosmos-sdk v0.45.2-inj

replace github.com/CosmWasm/wasmd => github.com/InjectiveLabs/wasmd v0.27.0-inj

require (
	github.com/InjectiveLabs/sdk-go v1.33.14
	github.com/cosmos/cosmos-sdk v0.45.0
	github.com/ethereum/go-ethereum v1.9.25
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.45.0
)
