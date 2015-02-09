// Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package types

import (
	"math"
	"os"
	"testing"
)

func TestStringifyInterface(t *testing.T) {

	testcases := []struct {
		name   string
		input  interface{}
		output string
	}{{
		name:   "string",
		input:  "foobar",
		output: "foobar",
	}, {
		name:   "uint8",
		input:  uint8(43),
		output: "43",
	}, {
		name:   "uint16",
		input:  uint16(43),
		output: "43",
	}, {
		name:   "uint32",
		input:  uint32(43),
		output: "43",
	}, {
		name:   "uint64",
		input:  uint64(43),
		output: "43",
	}, {
		name:   "max uint64",
		input:  uint64(math.MaxUint64),
		output: "18446744073709551615",
	}, {
		name:   "int8",
		input:  int8(-32),
		output: "-32",
	}, {
		name:   "int16",
		input:  int16(-32),
		output: "-32",
	}, {
		name:   "int32",
		input:  int32(-32),
		output: "-32",
	}, {
		name:   "int64",
		input:  int64(-32),
		output: "-32",
	}, {
		name:   "true",
		input:  true,
		output: "true",
	}, {
		name:   "false",
		input:  false,
		output: "false",
	}, {
		name:   "float32",
		input:  float32(2.345),
		output: "f1075188859",
	}, {
		name:   "float64",
		input:  float64(-34.6543),
		output: "f-4593298060402564373",
	}, {
		name: "map[string]interface{}",
		input: map[string]interface{}{
			"b": uint32(43),
			"a": "foobar",
			"ex": map[string]interface{}{
				"d": "barfoo",
				"c": uint32(45),
			},
		},
		output: "foobar_43_45_barfoo",
	}}

	tcaseFilter := os.Getenv("TESTCASE")
	for _, tcase := range testcases {
		if tcaseFilter != "" && tcase.name != tcaseFilter {
			continue
		}

		result, err := StringifyInterface(tcase.input)
		if err != nil {
			t.Errorf("Test %s: Error returned: %s", tcase.name, err.Error())
		} else if tcase.output != result {
			t.Errorf("Test %s: Result is different\nReceived: %s\nExpected: %s",
				tcase.name, result, tcase.output)
		}
	}
}