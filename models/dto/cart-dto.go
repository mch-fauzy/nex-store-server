package dto

import (
	"github.com/gofrs/uuid"
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/utils/failure"
)

type UserCartAddItemRequest struct {
	ProductId int    `json:"productId"`
	Quantity  int    `json:"quantity"`
	UserId    string `json:"-"`
	Email     string `json:"-"`
}

func (r UserCartAddItemRequest) Validate() error {
	if r.ProductId <= 0 {
		return failure.BadRequest("Product id is required and greater than 0")
	}

	if r.Quantity <= 0 {
		return failure.BadRequest("Quantity is required and greater than 0")
	}

	if r.UserId == "" {
		return failure.BadRequest("User id is required")
	}

	if r.Email == "" {
		return failure.BadRequest("Email is required")
	}

	return nil
}

func (r UserCartAddItemRequest) ToModel() models.UserCart {
	userId, _ := uuid.FromString(r.UserId)
	return models.UserCart{
		UserId:    userId,
		ProductId: r.ProductId,
		Quantity:  r.Quantity,
		CreatedBy: r.Email,
		UpdatedBy: r.Email,
	}
}

func (r UserCartAddItemRequest) IncreaseCartItem(cart *models.UserCart) {
	cart.Quantity += r.Quantity
	cart.UpdatedBy = r.Email
}

func (r UserCartAddItemRequest) UpdateCartQuantity(cart models.UserCart) models.UserCart {
	return models.UserCart{
		Quantity:  cart.Quantity + r.Quantity,
		UpdatedBy: r.Email,
	}
}

type UserCartGetListRequest struct {
	UserId string `json:"-"`
}

type UserCartGetListResponse struct {
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float32 `json:"totalPrice"`
}
