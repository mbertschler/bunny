// Copyright 2018 Martin Bertschler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memory

import (
	"reflect"
	"testing"
)

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

func TestArraySortingOld(t *testing.T) {
	for i, test := range cases {
		in := make([]int, len(test.Input))
		copy(in, test.Input)
		out, err := sortArrayOld(in, test.OldIndex, test.NewIndex)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(out, test.Output) {
			t.Error("old sorter", i, out, "should be", test.Output)
		}
	}
}

func TestArraySorting(t *testing.T) {
	for i, test := range cases {
		in := make([]int, len(test.Input))
		copy(in, test.Input)
		out, err := sortArray(in, test.OldIndex, test.NewIndex)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(out, test.Output) {
			t.Error("new sorter", i, out, "should be", test.Output)
		}
	}
}

func BenchmarkSliceSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range cases {
			sortArrayOld(test.Input, test.OldIndex, test.NewIndex)
		}
	}
}

func BenchmarkCounterSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range cases {
			sortArray(test.Input, test.OldIndex, test.NewIndex)
		}
	}
}
