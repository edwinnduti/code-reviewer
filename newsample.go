package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {

	url := os.Getenv("PAYMENT_URL")
	if url == "" {
		fmt.Println("environment variable not set. Please set it to your API key.")
		return
	}
	method := "POST"

	payload := strings.NewReader(`{
    "ShopperAmount": "1",
    "CustomerMSISDN": "254743143907",
    "ShopperCurrency": "KES",
    "Country": "KEN",
    "OriginatorCountry": "KEN",
    "OrganisationShortCode": "444408",
    "TransactionReference": "116453",
    "ThirdPartyConversationID": "1q12114yj109",
    "PurchasedItemsDesc": "Harry Potter and the OpenAPI",
    "VodaPartnerVirtualShortCode": "135456",
    "RealOrgName": "bokutestUSDcurrency",
    "ProcessingCurrency": "USD",
    "AmountInSettlementCurrency":"1.00",
    "SettlementCurrency":"USD",
    "MerchantCategoryCode": "5330",
    "InterchangeTariff": "0",
    "Language": "EN",
    "ForeignExchangeOption": "2",
    "ResultPageURL": "https://www.4cit.group/",
    "ResultURLDestination": "https://webhook.site/b385e50f-3f1e-443b-9349-6c89ed701bb8",
    "AutoReturn": false,
    "CustomSmsReceipt": "false",
    "ApiVersion": "3.1"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
