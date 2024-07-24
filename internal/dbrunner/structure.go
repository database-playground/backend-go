package dbrunner

type Input struct {
	Init  string `json:"init"`
	Query string `json:"query"`
}

type Output struct {
	Result [][]struct {
		Column string
		Value  string
	} `json:"result"`
}
