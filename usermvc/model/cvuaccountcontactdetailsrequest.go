package model

import "database/sql"

type Contactdetails struct {
	BillingAdvancedFilter  bool           `json:"billing_deep_filter"`
	ShippingAdvancedFilter bool           `json:"shipping_deep_filter"`
	FilterParam            string         `json:"deep_filter_args"`
	Create                 bool           `json:"create"`
	View                   bool           `json:"view"`
	Update                 bool           `json:"update"`
	UpdateConBilling       bool           `json:"contactbilling_update"`
	UpdateConShipping      bool           `json:"contactshipping_update"`
	CreateConBilling       bool           `json:"contactbilling_create"`
	CreateConShipping      bool           `json:"contactshipping_create"`
	AccountID              string         `json:"accountid"`
	Accountname            string         `json:"accountname"`
	Salutationid           int            `json:"salutationid"`
	Firstname              string         `json:"firstname"`
	Lastname               string         `json:"lastname"`
	Email                  string         `json:"email"`
	Position               string         `json:"position"`
	Phone                  string         `json:"phone"`
	Mobile                 string         `json:"mobile"`
	ContactOwner           string         `json:"contactowner"`
	ContactID              string         `json:"contactid"`
	CurrentSalesNo         string         `json:"current_sales_no"`
	PastSalesNo            string         `json:"past_sales_no"`
	BillingInfo            []BillingInfo  `json:"billinginfo"`
	ShippingInfo           []ShippingInfo `json:"shippinginfo"`
	Loggedinuserid         string         `json:"loggedinuserid"`
}
type BillingInfo struct {
	B_BillingID  string `json:"billing_id"`
	B_Street     string `json:"billing_street"`
	B_City       string `json:"billing_city"`
	B_State      string `json:"billing_state"`
	B_PostalCode string `json:"billing_postalcode"`
	B_Country    string `json:"billing_country"`
	B_Primary    bool   `json:"billing_primary"`
}
type ShippingInfo struct {
	S_ShippingID string `json:"shipping_id"`
	S_Street     string `json:"shipping_street"`
	S_City       string `json:"shipping_city"`
	S_State      string `json:"shipping_state"`
	S_PostalCode string `json:"shipping_postalcode"`
	S_Country    string `json:"shipping_country"`
	S_Primary    bool   `json:"shipping_primary"`
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
