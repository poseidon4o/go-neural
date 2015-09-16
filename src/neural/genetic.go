package neural

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const RAND_BUFFER_SIZE int = 300000

var values []float64
var generator rand.Rand
var idx int

var BUFFER_REFILS float64 = 0

func init() {
	idx = 0
	fmt.Println("Genetic: pre-generating random values...")
	values = make([]float64, RAND_BUFFER_SIZE, RAND_BUFFER_SIZE)
	generator = *rand.New(rand.NewSource(time.Now().UnixNano()))

	for c := 0; c < RAND_BUFFER_SIZE; c++ {
		values[c] = generator.Float64()
	}

	go func() {
		step := 1 / float64(RAND_BUFFER_SIZE)
		for {
			for c := 0; c < RAND_BUFFER_SIZE; c++ {
				values[c] = generator.Float64()
				if c%10 == 0 {
					time.Sleep(time.Microsecond)
				}
				BUFFER_REFILS += step
			}
		}
	}()
}

func Rand() float64 {
	idx++
	myIdx := idx
	return values[myIdx%RAND_BUFFER_SIZE]
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

func Signof(val float64) float64 {
	//return float64(int(val > 0) - int(val < 0))
	if val > 0 {
		return 1.
	} else if val < 0 {
		return -1.
	}
	return 0.
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

// Differs from Cross as this version does not generate random child
// in the case when mother and father are very different. When encountering
// big differences it will prefer the mother's genes
func Cross2(mother, father *Net) *Net {
	if father.Size() != father.Size() {
		panic("Cannot cross Nets with different sizes")
	}

	child := NewNet(mother.Size())

	for c := 0; c < child.Size(); c++ {
		for r := 0; r < child.Size(); r++ {
			if father.HasSynapse(c, r) != mother.HasSynapse(c, r) {
				continue
			}
			if father.HasSynapse(c, r) {
				fVal := *father.Synapse(c, r)
				mVal := *mother.Synapse(c, r)

				if Signof(fVal) != Signof(mVal) || math.Abs(math.Abs(fVal)-math.Abs(mVal)) > 0.3 {
					*child.Synapse(c, r) = mVal
				} else {
					*child.Synapse(c, r) = (fVal + mVal) / 2.
				}
			}
		}
	}

	return child
}
