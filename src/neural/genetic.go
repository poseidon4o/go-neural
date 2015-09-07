package neural

import (
	"math"
	"math/rand"
	"time"
)

const RAND_NUMBERS_SIZE int = 300000
const RAND_BUFFER_COUNT int = 10
const RAND_BUFFER_SIZE int = RAND_NUMBERS_SIZE / RAND_BUFFER_COUNT

var gen [RAND_BUFFER_COUNT]*rand.Rand
var genData = [RAND_BUFFER_COUNT]chan float64{
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
	make(chan float64, RAND_BUFFER_SIZE),
}

var ChanRand int = 0
var GlobRand int = 0

func init() {
	for c := 0; c < RAND_BUFFER_COUNT; c++ {
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
