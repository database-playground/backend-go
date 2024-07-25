package dbrunner

import (
	"testing"

	"github.com/samber/lo"
)

func TestInput_Hash(t *testing.T) {
	a := Input{
		Init:  "INIT",
		Query: "QUERY",
	}

	b := Input{
		Init:  "INIT",
		Query: "QUERY",
	}

	c := Input{
		Init:  "INIT",
		Query: "QUERY2AA",
	}

	d := Input{
		Init:  "init",
		Query: "query",
	}

	if a.Hash() != b.Hash() {
		t.Errorf("a.Hash() [%s] != b.Hash() [%s]", a.Hash(), b.Hash())
	}

	if a.Hash() == c.Hash() {
		t.Errorf("a.Hash() [%s] == c.Hash() [%s]", a.Hash(), c.Hash())
	}

	if a.Hash() == d.Hash() {
		t.Errorf("a.Hash() [%s] == d.Hash() [%s]", a.Hash(), d.Hash())
	}

	t.Logf("a.Hash() = %s", a.Hash())
	t.Logf("b.Hash() = %s", b.Hash())
	t.Logf("c.Hash() = %s", c.Hash())
	t.Logf("d.Hash() = %s", d.Hash())
}

func TestOutput_Hash(t *testing.T) {
	a := Output{
		Header: []string{"COL1", "COL2"},
		Data: [][]*string{
			{
				nil,
				lo.ToPtr("Hello!"),
			},
		},
	}

	b := Output{
		Header: []string{"COL1", "COL2"},
		Data: [][]*string{
			{
				nil,
				lo.ToPtr("Hello!"),
			},
		},
	}

	c := Output{
		Header: []string{"COL1", "COL2"},
		Data: [][]*string{
			{
				nil,
				lo.ToPtr("Hello!"),
			},
			{
				lo.ToPtr("Hello!"),
				nil,
			},
		},
	}

	d := Output{
		Header: nil,
		Data:   nil,
	}

	if lo.Must(a.Hash()) != lo.Must(b.Hash()) {
		t.Errorf("a.Hash() [%s] != b.Hash() [%s]", lo.Must(a.Hash()), lo.Must(b.Hash()))
	}

	if lo.Must(a.Hash()) == lo.Must(c.Hash()) {
		t.Errorf("a.Hash() [%s] == c.Hash() [%s]", lo.Must(a.Hash()), lo.Must(c.Hash()))
	}

	if lo.Must(a.Hash()) == lo.Must(d.Hash()) {
		t.Errorf("a.Hash() [%s] == d.Hash() [%s]", lo.Must(a.Hash()), lo.Must(d.Hash()))
	}

	t.Logf("a.Hash() = %s", lo.Must(a.Hash()))
	t.Logf("b.Hash() = %s", lo.Must(b.Hash()))
	t.Logf("c.Hash() = %s", lo.Must(c.Hash()))
	t.Logf("d.Hash() = %s", lo.Must(d.Hash()))
}
