package main

import (
	"testing"
	"fmt"
)

// TestRandomInRange tests randomInRange function.
func TestRandomInRange(t *testing.T) {

	fmt.Println("TestRandomInRange")

	tests := []struct {
		min            int
		max            int
		expectedResult int
	}{
		// Edge cases.
		{0, 0, 0},
		{1000, 1000, 1000},
		// Normal case.
		{500, 1000, -1},
		// Wrong parameters.
		{5000, 1000, 0},
		{1000, 500, 0},
	}

	for _, testData := range tests {
		givenResult := randomInRange(testData.min, testData.max)

		// Normal parameters, random results.
		if testData.expectedResult == -1 {
			if testData.min > givenResult || givenResult > testData.max {
				t.Errorf("randomInRange(%d, %d): Given: %d", testData.min, testData.max, givenResult)
			}
			continue
		}

		// Edge and wrong parameters, fixed result.
		if givenResult != testData.expectedResult {
			t.Errorf("randomInRange(%d, %d): Given: %d, Expected: %d", testData.min, testData.max, givenResult, testData.expectedResult)
		}
	}
}
