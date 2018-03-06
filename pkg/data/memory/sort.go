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
	"errors"
	"fmt"
)

func sortArrayOld(arr []int, old, new int) ([]int, error) {
	max := len(arr)
	if !(old < max && old >= 0 &&
		new < max && new >= 0) {
		return arr, errors.New(fmt.Sprintln(
			"invalid sorting from", old, "to", new, "max", max))
	}
	el := arr[old]
	arr = append(arr[:old], arr[old+1:]...)
	in := make([]int, len(arr))
	copy(in, arr)
	return append(append(arr[:new], el), in[new:]...), nil
}

func findInArray(in []int, search int) (int, bool) {
	for i := range in {
		if in[i] == search {
			return i, true
		}
	}
	return 0, false
}

func deleteFromArray(in []int, idx int) []int {
	out := []int{}
	for i, el := range in {
		if i != idx {
			out = append(out, el)
		}
	}
	return out
}

func sortArray(in []int, old, new int) ([]int, error) {
	if old == new {
		return in, nil
	}
	max := len(in)
	if !(old < max && old >= 0 &&
		new < max && new >= 0) {
		return in, errors.New(fmt.Sprintln(
			"invalid sorting from", old, "to", new, "max", max))
	}
	out := make([]int, len(in))
	i, j := 0, 0
	for j < max {
		if j == new {
			out[j] = in[old]
			j++
			continue
		}
		if i != old {
			out[j] = in[i]
			j++
		}
		i++
	}
	return out, nil
}
