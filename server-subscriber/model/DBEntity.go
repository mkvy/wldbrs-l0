package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type DataItem struct {
	ID        int
	OrderData OrderData
}

type OrderData struct {
	OrderUid    string `json:"order_uid" validate:"required,min=1"`
	TrackNumber string `json:"track_number" validate:"required"`
	Entry       string `json:"entry" validate:"required"`
	Delivery    struct {
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required"`
		Zip     string `json:"zip" validate:"required"`
		City    string `json:"city" validate:"required"`
		Address string `json:"address" validate:"required"`
		Region  string `json:"region" validate:"required"`
		Email   string `json:"email" validate:"required"`
	} `json:"delivery" validate:"required"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`
		RequestId    string `json:"request_id" validate:"omitempty,required"`
		Currency     string `json:"currency" validate:"required"`
		Provider     string `json:"provider" validate:"required"`
		Amount       int    `json:"amount" validate:"omitempty,required"`
		PaymentDt    int    `json:"payment_dt" validate:"omitempty,required"`
		Bank         string `json:"bank" validate:"required"`
		DeliveryCost int    `json:"delivery_cost" validate:"omitempty,required"`
		GoodsTotal   int    `json:"goods_total" validate:"omitempty,required"`
		CustomFee    int    `json:"custom_fee" validate:"omitempty,required"`
	} `json:"payment" validate:"required"`
	Items []struct {
		ChrtId      int    `json:"chrt_id" validate:"required"`
		TrackNumber string `json:"track_number" validate:"required"`
		Price       int    `json:"price" validate:"omitempty,required"`
		Rid         string `json:"rid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Sale        int    `json:"sale" validate:"omitempty,required"`
		Size        string `json:"size" validate:"required"`
		TotalPrice  int    `json:"total_price" validate:"omitempty,required"`
		NmId        int    `json:"nm_id" validate:"omitempty,required"`
		Brand       string `json:"brand" validate:"required"`
		Status      int    `json:"status" validate:"omitempty,required"`
	} `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature" validate:"omitempty,required"`
	CustomerId        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"omitempty,required"`
	Shardkey          string    `json:"shardkey" validate:"omitempty,required"`
	SmId              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"omitempty,required"`
}

func (od OrderData) Value() (driver.Value, error) {
	return json.Marshal(od)
}
func (od *OrderData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion to []byte failed")
	}
	return json.Unmarshal(b, &od)
}
