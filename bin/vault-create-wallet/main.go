package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brave-intl/bat-go/utils/altcurrency"
	"github.com/brave-intl/bat-go/utils/vaultsigner"
	"github.com/brave-intl/bat-go/wallet"
	"github.com/brave-intl/bat-go/wallet/provider/uphold"
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		log.Printf("Create a new wallet backed by vault.\n\n")
		log.Printf("Usage:\n\n")
		log.Printf("        %s WALLET_NAME\n\n", os.Args[0])
		log.Printf("  If a vault keypair exists with name WALLET_NAME, it will be used.\n")
		log.Printf("  Otherwise a new vault keypair with that name will be generated.\n\n")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Printf("ERROR: Must pass a single argument to name generated wallet / keypair\n\n")
		flag.Usage()
		os.Exit(1)
	}

	name := args[0]

	var info wallet.Info
	info.Provider = "uphold"
	info.ProviderID = ""
	{
		tmp := altcurrency.BAT
		info.AltCurrency = &tmp
	}

	client, err := vaultsigner.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	signer, err := vaultsigner.New(client, name)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Generated keypair with public key: %s\n", signer)

	wallet := &uphold.Wallet{Info: info, PrivKey: signer, PubKey: signer}
	err = wallet.Register(name)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Success, registered new keypair and wallet \"%s\"\n", name)
	fmt.Printf("Uphold card ID %s", wallet.Info.ProviderID)
	_, err = client.Logical().Write("wallets/"+name, map[string]interface{}{
		"providerId": wallet.Info.ProviderID,
	})
	if err != nil {
		log.Fatalln(err)
	}
}