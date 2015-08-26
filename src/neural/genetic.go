package neural

import (
	"math/rand"
	"time"
)

var gen = rand.New(rand.NewSource(time.Now().UnixNano()))

var values []float64 = nil
var idx int

func Seed(vals []float64) {
	values = vals
	idx = 0
}

func Rand() float64 {
	if values == nil {
		return gen.Float64()
	}
	if idx == len(values) {
		idx = 0
	}
	next := values[idx]
	idx++
	return next
}

func RandMax(max float64) float64 {
	return Rand() * max
}

func Chance(percent float64) bool {
	return percent > Rand()
}

func (n *Net) Randomize() {
	for c := 0; c < n.size; c++ {
		for r := 0; r < n.size; r++ {
			if n.HasSynapse(c, r) {
				*n.Synapse(c, r) = RandMax(2) - 1.0
			}
		}
	}
}

func (n *Net) MutateWithMagnitude(rate, magnitude float64) {
	for c := 0; c < n.size; c++ {
		for r := 0; r < n.size; r++ {
			if n.HasSynapse(c, r) && Chance(rate) {
				if Chance(0.5) {
					n.synapses[c][r] += (Rand() - 0.5) * magnitude
				} else {
					n.synapses[c][r] = RandMax(2) - 1.0
				}
			}
		}
	}
}

func (n *Net) Mutate(rate float64) {
	n.MutateWithMagnitude(rate, 1)
}

func Cross(mother, father *Net) *Net {
	if father.Size() != father.Size() {
		panic("Cannot cross Nets with different sizes")
		// return nil, errors.New()
	}

	parents := [2]*Net{father, mother}
	idx := 0

	child := NewNet(father.Size())

	for c := 0; c < child.Size(); c++ {
		for r := 0; r < child.Size(); r++ {
			if father.HasSynapse(c, r) != mother.HasSynapse(c, r) {
				continue
				// return nil, errors.New("Cannot cross Nets with missmatching synapses")
			}
			if father.HasSynapse(c, r) {
				*child.Synapse(c, r) = *parents[idx].Synapse(c, r)
				idx = (idx + 1) % 2
			}
		}
	}

	return child
}
