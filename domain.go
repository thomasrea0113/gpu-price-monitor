package monitor

type RequestMessage struct {
	ConfigOverrides *Config `json:"configOverrides"`
}

type PriceCheckJob struct {
	SiteName        string
	ProductName     string
	PriceThreshhold int
	Url             string
	PageContent     *string
}

// a struct to hold quantity/price information for a specific model of a product
type Model struct {
	Url    string
	Number string
	Price  float32
	Error  error
}

//
type PriceCheckResponse struct {
	Job    PriceCheckJob
	Error  error
	Models []Model
}

type JobContext struct {
}
