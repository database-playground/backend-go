package dbrunner

import (
	"crypto/sha256"
	"encoding/ascii85"
	"encoding/gob"
)

type Input struct {
	Init  string `json:"init"`
	Query string `json:"query"`
}

// Normalize returns a normalized input.
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

// Hash returns a hash of the input.
//
// It is encouraged to run [Input.Normalize] before hashing.
func (i Input) Hash() string {
	hash := sha256.New()

	hash.Write([]byte(i.Init))
	hash.Write([]byte(i.Query))

	hashed := hash.Sum(nil)

	output := make([]byte, ascii85.MaxEncodedLen(len(hashed)))
	ascii85.Encode(output, hashed)
	return string(output)
}

type Output struct {
	Result [][]struct {
		Column string
		Value  *string
	} `json:"result"`
}

func (o Output) Hash() (string, error) {
	hasher := sha256.New()
	encoder := gob.NewEncoder(hasher)

	err := encoder.Encode(o.Result)
	if err != nil {
		return "", err
	}

	hashed := hasher.Sum(nil)
	output := make([]byte, ascii85.MaxEncodedLen(len(hashed)))

	ascii85.Encode(output, hashed[:])
	return string(output), nil
}
