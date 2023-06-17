package compare

import (
	"fmt"
	"testing"
)

func Test_Array(t *testing.T) {
	var cases = []struct {
		a1     []string
		a2     []string
		result bool
	}{
		{nil, nil, true},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"a", "c"}, []string{"a", "b"}, false},
		{[]string{"a", "b"}, []string{"a", "b", "b"}, false},
		{[]string{"b", "a", "d", "c"}, []string{"a", "b", "c", "d"}, true},
		{[]string{"foo", "bar", "baz"}, []string{"bar", "baz", "foo"}, true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("array case %d", i), func(t *testing.T) {
			result := Array(c.a1, c.a2)
			if c.result != result {
				t.Errorf("exptected a1: %v compared to a2: %v to be: %v but was :%v", c.a1, c.a2, c.result, result)
			}
		})
	}

}
