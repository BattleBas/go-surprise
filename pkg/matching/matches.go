package matching

import (
	"errors"
	"math/rand"
	"time"
)

// Group represents a slice of people participating in Surprise
type Group struct {
	People []Person
}

// Person represents an individual participating in Surprise
type Person struct {
	Name    string
	Email   string
	Invalid []string
}

// isInvalid checks if a given name can be assigned to this person
func (p Person) isInvalid(name string) bool {
	if p.Name == name {
		return true
	}

	for _, n := range p.Invalid {
		if n == name {
			return true
		}
	}
	return false
}

// Matches represents a slice of matched participants
type Matches struct {
	Pairs []Pair
}

// Pair represents a giver and a reciever that have been matched together
type Pair struct {
	Giver    Person
	Reciever Person
}

// CreateMatches assigns each giver with a valid reciever
func CreateMatches(g Group, r Group) (Matches, error) {
	if len(g.People) != len(r.People) {
		return Matches{}, errors.New("there should be an equal number of givers and receivers")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(r.People), func(i, j int) {
		r.People[i], r.People[j] = r.People[j], r.People[i]
	})

	matches := Matches{}
	for i := 0; i < len(g.People); i++ {
		matches.Pairs = append(matches.Pairs, Pair{Giver: g.People[i], Reciever: r.People[i]})
	}

	for {
		valid, i := validateMatches(matches)
		if !valid {
			err := swapMatch(&matches, i)
			if err != nil {
				return Matches{}, err
			}
			continue
		}
		break
	}

	return matches, nil
}

// validateMatches iterates over possible matches to see if every giver has been assigned a valid reciever
func validateMatches(m Matches) (bool, int) {

	for i, pair := range m.Pairs {
		if pair.Giver.isInvalid(pair.Reciever.Name) {
			return false, i
		}
	}

	return true, -1
}

// swapMatch looks for another person to swap recievers with to create two valid matches
func swapMatch(m *Matches, swapIndex int) error {

	swap := m.Pairs[swapIndex]
	for i, pair := range m.Pairs {
		if i == swapIndex {
			continue
		}

		if swap.Giver.isInvalid(pair.Reciever.Name) || pair.Giver.isInvalid(swap.Reciever.Name) {
			continue
		}

		tmp := swap.Reciever
		m.Pairs[swapIndex].Reciever = pair.Reciever
		m.Pairs[i].Reciever = tmp
		return nil
	}

	return errors.New("No valid swap")
}
