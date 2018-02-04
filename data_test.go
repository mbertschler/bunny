package main

import (
	"reflect"
	"testing"
)

func TestArraySorting(t *testing.T) {
	var cases = []struct {
		Input    []int
		OldIndex int
		NewIndex int
		Output   []int
	}{
		{
			Input:    []int{1, 2, 3, 4, 5, 6},
			OldIndex: 1,
			NewIndex: 0,
			Output:   []int{2, 1, 3, 4, 5, 6},
		},
		{
			Input:    []int{1, 2, 3, 4, 5, 6},
			OldIndex: 0,
			NewIndex: 1,
			Output:   []int{2, 1, 3, 4, 5, 6},
		},
		{
			Input:    []int{1, 2, 3, 4, 5, 6},
			OldIndex: 3,
			NewIndex: 1,
			Output:   []int{1, 4, 2, 3, 5, 6},
		},
		{
			Input:    []int{1, 2, 3, 4, 5, 6},
			OldIndex: 1,
			NewIndex: 3,
			Output:   []int{1, 3, 4, 2, 5, 6},
		},
		{
			Input:    []int{1, 2, 3, 4, 5, 6},
			OldIndex: 0,
			NewIndex: 5,
			Output:   []int{2, 3, 4, 5, 6, 1},
		},
		{
			Input:    []int{1, 2, 3, 4, 5, 6},
			OldIndex: 5,
			NewIndex: 0,
			Output:   []int{6, 1, 2, 3, 4, 5},
		},
	}
	for i, test := range cases {
		out := sortArray(test.Input, test.OldIndex, test.NewIndex)
		if !reflect.DeepEqual(out, test.Output) {
			t.Error(i, "arrays are not equal:", out, "should be", test.Output)
		}
	}
}
