package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/chimoney/chimoney-go"
	"github.com/chimoney/chimoney-go/modules/payouts"
)

func main() {
	client := chimoney.New(
		chimoney.WithAPIKey("your-api-key-here"),
		chimoney.WithSandbox(true),
	)

	ctx := context.Background()

	txns, err := client.Account.GetAllTransactions(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	var txnsList []interface{}
	if err := json.Unmarshal(txns.Data, &txnsList); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of transactions: %d\n", len(txnsList))


	assets, err := client.Info.GetSupportedAssets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var assetsResp struct {
		BenefitsList []interface{} `json:"benefitsList"`
	}
	if err := json.Unmarshal(assets.Data, &assetsResp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of supported assets: %d\n", len(assetsResp.BenefitsList))

	countries, err := client.Info.GetAirtimeCountries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var countriesList []interface{}
	if err := json.Unmarshal(countries.Data, &countriesList); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of airtime countries: %d\n", len(countriesList))

	codes, err := client.Info.GetMobileMoneyCodes(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var codesList []interface{}
	if err := json.Unmarshal(codes.Data, &codesList); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of mobile money codes: %d\n", len(codesList))

	banks, err := client.Info.GetBanks(ctx, "NG")
	if err != nil {
		log.Fatal(err)
	}
	var banksResp []interface{}
	if err := json.Unmarshal(banks.Data, &banksResp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of banks in Nigeria: %d\n", len(banksResp))

	walletResp, err := client.Wallet.List(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	var walletList []interface{}
	if err := json.Unmarshal(walletResp.Data, &walletList); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of wallets: %d\n", len(walletList))

	chimoneyPayload := []payouts.ChimoneyPayload{
		{
			Email:      "heypleasant@gmail.com",
			ValueInUSD: 10,
		},
	}

	payoutResp, err := client.Payouts.Chimoney(ctx, chimoneyPayload, "")
	if err != nil {
		log.Fatal(err)
	}
	var payoutData struct {
		ID          string  `json:"id"`
		Status      string  `json:"status"`
		ValueInUSD  float64 `json:"valueInUsd"`
		ChimoneyID  string  `json:"chimoneyId"`
		Email       string  `json:"email"`
	}

	if err := json.Unmarshal(payoutResp.Data, &payoutData); err != nil {
		log.Fatal(err)
	}

	prettyJSON, err := json.MarshalIndent(payoutData, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Payout response:\n%s\n", string(prettyJSON))

}
