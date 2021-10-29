package rmq

import (
	"encoding/json"

	daas "github.com/touno-io/goasa"
)

type RMQInbound struct {
	Host    string `json:"host"`
	Queue   string `json:"queue"`
	Channel string `json:"channel"`
}
type RMQOutbound struct {
	Host    string       `json:"host"`
	Publish []RMQPublish `json:"publish"`
}
type RMQPublish struct {
	Queue   string `json:"queue"`
	Channel string `json:"channel"`
	Route   string `json:"route"`
}

func ParseRMQInbound(raw []byte) RMQInbound {
	var data RMQInbound
	if err := json.Unmarshal(raw, &data); err != nil {
		daas.Fatal("JSON ParseRMQInbound::", err)
	}
	return data
}

func ParseRMQOutbound(raw []byte) []RMQOutbound {
	var data []RMQOutbound
	if err := json.Unmarshal(raw, &data); err != nil {
		daas.Fatal("JSON ParseRMQOutbound::", err)
	}
	return data
}

type Inbound struct {
	BU      string           `json:"bu"`
	Data    []InboundBarcode `json:"data"`
	TraceID *string          `json:"trace_uuid"`
}

type InboundBarcode struct {
	Barcode   string `json:"barcode"`
	StoreCode string `json:"loc"`
	SKU       string `json:"sku"`
}

type Outbound struct {
	Ack     bool        `json:"ack"`
	Store   int         `json:"store"`
	SKU     int         `json:"sku"`
	Missing *[]string   `json:"missing,omitempty"`
	Result  interface{} `json:"result"`
}

func ParseInbound(raw *[]byte) Inbound {
	var data Inbound
	if err := json.Unmarshal(*raw, &data); err != nil {
		daas.Fatal("JSON ParseInbound::", err)
	}
	return data
}
