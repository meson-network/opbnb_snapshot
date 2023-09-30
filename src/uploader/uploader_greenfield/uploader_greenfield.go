package uploader_greenfield

import (
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
)



func genGreenfieldClient(privateKey string,rpcAddr string, chainId string) (client.IClient, error) {
	// greenfield client
	// import account
	account, err := types.NewAccountFromPrivateKey("meson", privateKey)
	if err != nil {
		return nil,err
	}

	// create client
	cli, err := client.New(chainId, rpcAddr, client.Option{DefaultAccount: account})
	if err != nil {
		return nil,err
	}
	return cli, nil
}