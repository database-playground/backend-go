package dbrunner

type Input struct {
	Init  string `json:"init"`
	Query string `json:"query"`
}

func (i Input) Normalize() (Input, error) {
	// Init should not be normalized and usually constrant.
	normlizedQuery, err := FormatSQL(i.Query)
	if err != nil {
		return Input{}, err
	}

	return Input{
		Init:  i.Init,
		Query: normlizedQuery,
	}, nil
}

type Output struct {
	Result [][]struct {
		Column string
		Value  *string
	} `json:"result"`
}
