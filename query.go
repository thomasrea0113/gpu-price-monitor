package monitor

import "github.com/thomasrea0113/gpu-price-monitor/domain"

func QueryBestBuy(j domain.PriceCheckJob, keyword string) []domain.Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]domain.Model, 0)
}

func QueryWalMart(j domain.PriceCheckJob, keyword string) []domain.Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]domain.Model, 0)
}

func QueryNewegg(j domain.PriceCheckJob, keyword string) []domain.Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]domain.Model, 0)
}

func QueryMicroCenter(j domain.PriceCheckJob, keyword string) []domain.Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]domain.Model, 0)
}
