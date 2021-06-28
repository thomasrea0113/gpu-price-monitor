package domain

type RequestMessage struct {
	ConfigOverrides *Config `json:"configOverrides"`
}

type PriceCheckJob struct {
	Site    Site
	Product Product
}

type PriceCheckResponse struct {
	ProductName string
	Models      []struct {
		Name              string
		QuantityAvailable int
		price             float32
	}
}
