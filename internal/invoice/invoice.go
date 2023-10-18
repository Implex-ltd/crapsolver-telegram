package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rizalgowandy/coinbase-commerce-go"
	"github.com/rizalgowandy/coinbase-commerce-go/pkg/api"
	"github.com/rizalgowandy/coinbase-commerce-go/pkg/entity"
	"github.com/rizalgowandy/coinbase-commerce-go/pkg/enum"
)

const (
	api_key = "1f691521-78b2-4bb5-b655-76b23a1190e0"
)

func CalculateOrder(x int) float64 {
	coutParRequete := 1
	coutParMilleRequetes := 2.0

	coutEnCredits := x * coutParRequete
	coutEnDollars := (float64(coutEnCredits) / 1000) * coutParMilleRequetes

	return coutEnDollars
}

func WaitForOrder(orderId, orderCode string) bool {
	t := time.NewTicker(time.Minute * 5)
	st := time.Now()
	client, err := coinbase.NewClient(api.Config{
		Key:   api_key,
		Debug: false,
	})

	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case <-t.C:
			showResp, err := client.ShowCharge(context.Background(), &entity.ShowChargeReq{
				ChargeCode: orderCode,
				ChargeID:   orderId,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			status := showResp.Data.Timeline[len(showResp.Data.Timeline)-1].Status
			fmt.Println(status, orderId, orderCode)

			switch status {
			case "NEW":
				continue
			case "COMPLETED":
				return true
			case "CANCELED", "REFUNDED", "PENDING REFUND", "EXPIRED", "UNRESOLVED", "RESOLVED":
				return false
			}
		default:
			if time.Since(st).Hours() > 2 {
				return false
			}
		}
	}
}

func NewInvoice(name string, amount int) (*entity.CreateChargeResp, error) {
	client, err := coinbase.NewClient(api.Config{
		Key:   api_key,
		Debug: false,
	})

	if err != nil {
		return nil, err
	}

	return client.CreateCharge(context.Background(), &entity.CreateChargeReq{
		Name:        fmt.Sprintf("Crapsolver - %d credit refill", amount),
		Description: fmt.Sprintf("API Credit - 0,002/c (%d x 0,002)", amount),
		LocalPrice: entity.CreateChargePrice{
			Amount:   fmt.Sprintf("%.2f", CalculateOrder(amount)),
			Currency: "USD",
		},
		PricingType: enum.PricingTypeFixedPrice,
		Metadata: map[string]string{
			"CustomerID":   uuid.NewString(),
			"CustomerName": name,
		},
		RedirectURL: "https://t.me/crapsolver_bot",
		CancelURL:   "https://t.me/crapsolver_bot",
	})
}
