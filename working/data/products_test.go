package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Mahantesh",
		Price: 2.00,
		//SKU:   "abc-avw-afg",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
