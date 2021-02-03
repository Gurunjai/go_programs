package main

import "fmt"

func main() {
	lis := []int{2, -5, 3, 4, -1, 2, 1, 3, -4, -3}
	// lis := []int{2, -1, 3, 4, -10, 2, 1, 3, 4, -3}
	largestSumSubArray(lis)
}

func largestSumSubArray(in []int) {
	if len(in) <= 0 {
		panic("zero length array!!!!!")
	}
	maxScore := in[0]
	prevScore := in[0]
	dex := 0
	st := 0

	for i := 1; i < len(in); i++ {
		newScore := prevScore + in[i]
		// prevScore = newScore

		if newScore <= 0 {
			// +1 - new score hit a negative, start from next index
			// +1 - Accumulate the index of 0
			st, dex = i + 1, i + 1
			maxScore = 0			
		}
		
		if maxScore < newScore {
			maxScore = newScore
			// +1 - to adjust the index 0 on the element
			dex = i
		}
		prevScore = newScore
	}

	fmt.Printf("Input Array: %+v\n", in)
	fmt.Printf("\tMax Score: %v\n \tStart Index: %v \t End Index: %v\n",
		maxScore, st, dex)
}