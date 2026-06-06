package payment

import (
	"context"

	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

type Client struct {
	client payment_v1.PaymentServiceClient
}

func NewClient(client payment_v1.PaymentServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) PayOrder(
	ctx context.Context,
	orderUUID string,
	userUUID string,
	method sharedmodel.PaymentMethod,
) (string, error) {
	resp, err := c.client.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: payment_v1.PaymentMethod(method),
	})
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}
