package main

import "testing"

func TestTwoSum(t *testing.T) {
	input := []int{2, 7, 11, 15}
	target := 9
	want := []int{0, 1}

	got := twoSum(input, target)

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Mistmatch at index: %v\n\t Got: %v\n\t Want: %v\n", i, got[i], want[i])
		}
	}
}
