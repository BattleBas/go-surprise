package matching

import (
	"reflect"
	"sort"
	"testing"
)

func TestCreateMatches(t *testing.T) {

	tests := []struct {
		name        string
		givers      Group
		receivers   Group
		matches     Matches
		expectError bool
	}{
		{"no-valid-pairs",
			Group{
				[]Person{
					{Name: "Bob"},
					{Name: "Bob"},
				},
			},
			Group{
				[]Person{
					{Name: "Bob"},
					{Name: "Bob"},
				},
			},
			Matches{},
			true,
		},
		{"valid-pairs",
			Group{
				[]Person{
					{Name: "Bob"},
					{Name: "Eric"},
				},
			},
			Group{
				[]Person{
					{Name: "Bob"},
					{Name: "Eric"},
				},
			},
			Matches{
				[]Pair{
					{
						Giver:    Person{Name: "Bob"},
						Reciever: Person{Name: "Eric"},
					},
					{
						Giver:    Person{Name: "Eric"},
						Reciever: Person{Name: "Bob"},
					},
				},
			},
			false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {

			matches, err := CreateMatches(test.givers, test.receivers)
			if err != nil && !test.expectError {
				t.Fatalf("did not expect an error: %v", err)
			}
			if err == nil && test.expectError {
				t.Fatalf("was expecting an error: %s", err)
			}

			// CreateMatches shuffles the test.givers slice
			sort.Slice(matches.Pairs, func(i, j int) bool {
				return matches.Pairs[i].Giver.Name < matches.Pairs[j].Giver.Name
			})
			if !reflect.DeepEqual(test.matches, matches) {
				t.Fatalf("wrong matches")
			}
		})
	}

}
