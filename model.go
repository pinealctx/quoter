package quoter

type SubscribeRequest struct {
	Method string `json:"method"`
	ID     string `json:"id"`
	Data   struct {
		CryptoIDs []int `json:"cryptoIds"`
	} `json:"data"`
}

func NewSubscribeRequest(cryptoIDs []int) *SubscribeRequest {
	return &SubscribeRequest{
		Method: "subscribe",
		ID:     "price",
		Data: struct {
			CryptoIDs []int `json:"cryptoIds"`
		}{
			CryptoIDs: cryptoIDs,
		},
	}
}

type SubscribeResponse struct {
	ID   string `json:"id"`
	Data struct {
		Timestamp int64 `json:"t"`
		CR        struct {
			ID    int     `json:"id"`
			Price float64 `json:"p"`
		} `json:"cr"`
	} `json:"d"`
	S string `json:"s"`
}

type LatestResponse struct {
	Data   []*LatestData `json:"data"`
	Status *Status       `json:"status"`
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code,string"`
	ErrorMessage string `json:"error_message"`
}

type LatestData struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Symbol string   `json:"symbol"`
	Slug   string   `json:"slug"`
	Quotes []*Quote `json:"quotes"`
}

type Quote struct {
	Name  int     `json:"name,string"`
	Price float64 `json:"price"`
}
