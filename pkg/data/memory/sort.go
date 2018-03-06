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
