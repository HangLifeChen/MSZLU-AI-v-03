package admin

type GrowthTrendResponse struct {
	Months []string `json:"months"`
	Data   []int64  `json:"data"`
}
