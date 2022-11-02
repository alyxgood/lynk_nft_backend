package utils

import (
	"testing"
)

func TestDecode(t *testing.T) {

	var input = "0000000000000000000000000000000000000000000000000000000000000000"
	var outputParameters []string
	outputParameters = append(outputParameters, "bool")

	data, err := Decode(outputParameters, input)
	if err != nil {
		println(err)
	} else {
		println(data)
	}
}
