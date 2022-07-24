package version

import (
	"testing"
)

var parseTests = []struct {
	suc bool
	str string
	res Version
}{
	{true, "1.2.3_1", Version{[]int{1, dot, 2, dot, 3}, 1}},
	{true, "1.2.3_", Version{[]int{1, dot, 2, dot, 3}, 0}},
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
		if ver := Parse(tt.str); ver.Cmp(tt.res) != 0 {
			t.Errorf("got %v, expected %v", ver, tt.res)
		}
	}
}

var cmpTests = []struct {
	r int
	a string
	b string
}{
	{0, "", ""},
	{-1, "1", "2"},
	{1, "2", "1"},
	{1, "2", "1"},

	// tests lifted from xbps
	{0, "1.0", "1.0"},
	{-1, "1.0", "1.0_1"},
	{-1, "2.0rc2", "2.0rc3"},
	{1, "2.0rc3", "2.0rc2"},
	{-1, "129", "129_1"},
	{0, "21", "21_0"},
	{1, "21", "2.1"},
	{1, "1.0.1", "1.0_1"},
	{-1, "2.0rc3", "2.0"},
}

func TestCmp(t *testing.T) {
	for _, tt := range cmpTests {
		if r := Cmp(tt.a, tt.b); r != tt.r {
			t.Errorf("Cmp(%q, %q): got %d, expected %d", tt.a, tt.b, r, tt.r)
		}
	}
}
