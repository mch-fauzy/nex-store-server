package dto

import "github.com/nexmedis-be-technical-test/utils/failure"

type TransactionTopUpBalanceByUserIdRequest struct {
	TopupAmount int    `json:"topupAmount"`
	UserId      string `json:"-"`
	Email       string `json:"-"`
}

func (r TransactionTopUpBalanceByUserIdRequest) Validate() error {
	if r.TopupAmount <= 10000 {
		return failure.BadRequest("Topup amount is required and greater than 10000")
	}

	if r.UserId == "" {
		return failure.BadRequest("User id is required")
	}

	if r.Email == "" {
		return failure.BadRequest("Email is required")
	}

	return nil
}

type TransactionTopUpBalanceByUserIdResponse struct {
	Balance float32 `json:"balance"`
}

type TransactionWithdrawBalanceByUserIdRequest struct {
	WithdrawAmount int    `json:"withdrawAmount"`
	UserId         string `json:"-"`
	Email          string `json:"-"`
}

func (r TransactionWithdrawBalanceByUserIdRequest) Validate() error {
	if r.WithdrawAmount < 0 {
		return failure.BadRequest("Withdraw amount is required and greater than 0")
	}

	if r.UserId == "" {
		return failure.BadRequest("User id is required")
	}

	if r.Email == "" {
		return failure.BadRequest("Email is required")
	}

	return nil
}

type TransactionWithdrawBalanceByUserIdResponse struct {
	Balance float32 `json:"balance"`
}

type TransactionPurchaseCartRequest struct {
	UserId string `json:"-"`
	Email  string `json:"-"`
}

func (r TransactionPurchaseCartRequest) Validate() error {
	if r.UserId == "" {
		return failure.BadRequest("User id is required")
	}

	if r.Email == "" {
		return failure.BadRequest("Email is required")
	}

	return nil
}

type TransactionPurchaseCartResponse struct {
	InvoiceNumber string  `json:"invoiceNumber"`
	TotalAmount   float32 `json:"totalAmount"`
}
