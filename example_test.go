package solve_test

import (
	"github.com/bertbaron/solve"
	"fmt"
)

type state struct {
	// the state of the vector
	vector [5]byte
	// the cost to get to this state
	cost   int
	// index of the element that was swapped with its right neigbour
	index  int
}

func (s state) Id() interface{} {
	return s.vector
}

func (s state) Expand() []solve.State {
	n := len(s.vector) - 1
	steps := make([]solve.State, n, n)
	for i := 0; i < n; i++ {
		copy := s.vector
		copy[i], copy[i + 1] = copy[i + 1], copy[i]
		steps[i] = state{copy, s.cost + 1, i}
	}
	return steps
}

func (s state) IsGoal() bool {
	for i := 1; i < len(s.vector); i++ {
		if s.vector[i - 1] > s.vector[i] {
			return false
		}
	}
	return true
}

func (s state) Cost() float64 {
	return float64(s.cost)
}

func (s state) Heuristic() float64 {
	return 0
}

// Finds the minumum number of swaps of neighbouring elements required to
// sort a vector
func Example() {
	s := state{[...]byte{3, 2, 5, 4, 1}, 0, -1}
	result := solve.NewSolver(s).
		Algorithm(solve.IDAstar).
		Constraint(solve.NO_LOOP).
		Solve()
	for _, st := range result.Solution {
		fmt.Printf("%v\n", st.(state).vector)
	}
	// Output:
	// [3 2 5 4 1]
	// [3 2 5 1 4]
	// [3 2 1 5 4]
	// [3 2 1 4 5]
	// [3 1 2 4 5]
	// [1 3 2 4 5]
	// [1 2 3 4 5]
}
