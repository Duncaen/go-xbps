package version

import (
	"testing"
)

var parseTests = []struct {
	suc bool
	str string
	res Version
}{
	{true, "1.2.3_1", Version{[]int{1, Dot, 2, Dot, 3}, 1}},
	{true, "1.2.3_", Version{[]int{1, Dot, 2, Dot, 3}, 0}},
	{true, "1_1", Version{[]int{1}, 1}},
	{true, "1_1_1", Version{[]int{1}, 1}},
	{true, "1_2_1", Version{[]int{1}, 1}},
	{true, "1_1_3", Version{[]int{1}, 3}},
	{true, "1", Version{[]int{1}, 0}},
	{true, "1_", Version{[]int{1}, 0}},
	{true, "_", Version{[]int{}, 0}},
	{true, "_1", Version{[]int{}, 1}},
	{true, "", Version{[]int{}, 0}},
	{true, "abc_1", Version{[]int{0, 1, 0, 2, 0, 3}, 1}},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		if ver := Parse(tt.str); !ver.Equal(tt.res) {
			t.Logf("got %v, expected %v", ver, tt.res)
		}
	}
}
