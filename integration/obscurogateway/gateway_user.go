package faucet

import (
	"context"
	"fmt"

	"github.com/ten-protocol/go-ten/tools/walletextension/lib"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/go/wallet"
)

// GatewayUser TODO (@ziga) refactor GatewayUser and integrate it with OGlib.
// GatewayUser is a struct that includes everything a gateway user has and uses (userID, wallets, http & ws addresses and client )
type GatewayUser struct {
	Wallets    []wallet.Wallet
	HTTPClient *ethclient.Client
	WSClient   *ethclient.Client
	tgClient   *lib.TGLib
}

func NewUser(wallets []wallet.Wallet, serverAddressHTTP string, serverAddressWS string) (*GatewayUser, error) {
	ogClient := lib.NewTenGatewayLibrary(serverAddressHTTP, serverAddressWS)

	// automatically join
	err := ogClient.Join()
	if err != nil {
		return nil, err
	}

	// create clients
	httpClient, err := ethclient.Dial(serverAddressHTTP + "/v1/" + "?token=" + ogClient.UserID())
	if err != nil {
		return nil, err
	}
	wsClient, err := ethclient.Dial(serverAddressWS + "/v1/" + "?token=" + ogClient.UserID())
	if err != nil {
		return nil, err
	}

	return &GatewayUser{
		Wallets:    wallets,
		HTTPClient: httpClient,
		WSClient:   wsClient,
		tgClient:   ogClient,
	}, nil
}

func (u GatewayUser) RegisterAccounts() error {
	for _, w := range u.Wallets {
		err := u.tgClient.RegisterAccount(w.PrivateKey(), w.Address())
		if err != nil {
			return err
		}
		fmt.Printf("Successfully registered address %s for user: %s.\n", w.Address().Hex(), u.tgClient.UserID())
	}

	return nil
}

func (u GatewayUser) PrintUserAccountsBalances() error {
	for _, w := range u.Wallets {
		balance, err := u.HTTPClient.BalanceAt(context.Background(), w.Address(), nil)
		if err != nil {
			return err
		}
		fmt.Println("Balance for account ", w.Address().Hex(), " - ", balance.String())
	}
	return nil
}
