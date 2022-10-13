package coveragetesting

import "testing"

type Test struct {
	in  int
	out string
}

var tests = []Test{
	{
		in:  -1,
		out: "negative",
	},
	{
		in:  0,
		out: "zero",
	},
	{
		in:  5,
		out: "small",
	},
	{
		in:  10,
		out: "big",
	},
	{
		in:  100,
		out: "huge",
	},
	{
		in:  1000,
		out: "enormous",
	},
}

func TestSize(t *testing.T) {
	for i, test := range tests {
		size := Size(test.in)
		if size != test.out {
			t.Errorf("#%d: Size(%d)=%s; want %s", i, test.in, size, test.out)
		}
	}
}
