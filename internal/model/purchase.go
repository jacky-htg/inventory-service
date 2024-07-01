package model

import (
	"context"
	"inventory-service/pb/inventories"
	"inventory-service/pb/purchases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Purchase struct {
	PurchaseClient purchases.PurchaseServiceClient
	Id             string
}

func (u *Purchase) Outstanding(ctx context.Context) (*inventories.OutstandingResponse, error) {
	var output inventories.OutstandingResponse
	out, err := u.PurchaseClient.GetOutstandingPurchaseDetails(ctx, &purchases.OutstandingPurchaseRequest{Id: u.Id})
	if s, ok := status.FromError(err); !ok {
		if s.Code() == codes.Unknown {
			err = status.Errorf(codes.Internal, "Error when calling purchase.GetOutstandingPurchaseDetails service: %s", err)
		}

		return &output, err
	}
	for _, v := range out.Detail {
		output.Detail = append(output.Detail, &inventories.OutstandingDetail{
			ProductId: v.ProductId,
			Quantity:  uint32(v.Quantity),
		})
	}

	return &output, nil
}
