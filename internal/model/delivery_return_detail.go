package model

import "inventory-service/pb/inventories"

// DeliveryReturnDetail struct
type DeliveryReturnDetail struct {
	Pb               inventories.DeliveryReturnDetail
	PbDeliveryReturn inventories.DeliveryReturn
}
