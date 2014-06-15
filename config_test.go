package main

import "testing"

func TestParseLine(t *testing.T) {
	inputs := []string{
		"",
		"#",
		"1",
		"e",
		"e e",
		// 5
		"e e e e",
		"er er er er",
		" e",
		" e e",
		"	e",
		// 10
		"	e	e",
		"	e	e#adfasfd",
		"	e	ee#adfasfd",
		"	ee	e#adfasfd",
	}

	expectedTriplets := []Triplet{
		Triplet{},
		Triplet{},
		Triplet{"1", "", ""},
		Triplet{"e", "", ""},
		Triplet{"e", "e", ""},
		// 5
		Triplet{"e", "e e e", ""},
		Triplet{"er", "er er er", ""},
		Triplet{"e", "", ""},
		Triplet{"e", "e", ""},
		Triplet{"e", "", ""},
		// 10
		Triplet{"e", "e", ""},
		Triplet{"e", "e", "adfasfd"},
		Triplet{"e", "ee", "adfasfd"},
		Triplet{"ee", "e", "adfasfd"},
	}

	if len(inputs) != len(expectedTriplets) {
		t.Fatal("2921883114 invalid test harness, expected ", len(inputs), " == ", len(expectedTriplets))
	}

	for i, input := range inputs {
		candidate := parseLine(i, input)
		if candidate.IsEqual(expectedTriplets[i]) == false {
			t.Error(i, "Candidate != Expected:\n", candidate, "\n", expectedTriplets[i])
		}
	}
}
