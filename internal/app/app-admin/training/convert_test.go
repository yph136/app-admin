package training

import (
	"fmt"
	"testing"
)

func TestConvertToArgs(t *testing.T) {
	args := make(map[string]string)
	args["test"] = "test1"
	args["test1"] = "test2"
	fmt.Println(convertToArgs(args))
}
