package neural

import (
	"math/rand"
	"time"
)

var gen [10]*rand.Rand

var genData = [10]chan float64{
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
	make(chan float64, 1000),
}

var ChanRand int = 0
var GlobRand int = 0

func init() {
	for c := 0; c < 10; c++ {
		gen[c] = rand.New(rand.NewSource(time.Now().UnixNano()))
		go func(r int) {
			for {
				genData[r] <- gen[r].Float64()
			}
		}(c)
	}
}

func Rand() float64 {
	chRead := false
	var f float64
	select {
	case f = <-genData[0]:
		chRead = true
		break
	case f = <-genData[1]:
		chRead = true
		break
	case f = <-genData[2]:
		chRead = true
		break
	case f = <-genData[3]:
		chRead = true
		break
	case f = <-genData[4]:
		chRead = true
		break
	case f = <-genData[5]:
		chRead = true
		break
	case f = <-genData[6]:
		chRead = true
		break
	case f = <-genData[7]:
		chRead = true
		break
	case f = <-genData[8]:
		chRead = true
		break
	case f = <-genData[9]:
		chRead = true
		break
	default:
		f = rand.Float64()
	}
	if chRead {
		ChanRand++
	} else {
		GlobRand++
	}
	return f
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
