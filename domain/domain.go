package domain

type RequestMessage struct {
	ConfigOverrides *Config `json:"configOverrides"`
}

type PriceCheckJob struct {
	Site    Site
	Product Product
}

// a struct to hold quantity/price information for a specific model of a product
type Model struct {
	Number            string
	QuantityAvailable int
	Price             float32
	Error             error
}

//
type PriceCheckResponse struct {
	Job    PriceCheckJob
	Error  error
	Models []Model
}
