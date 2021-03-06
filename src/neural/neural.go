package neural

import (
	"fmt"
	"io"
	"math"
	"os"
)

type Neuron struct {
	fire float64
}

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(x))
}

func NewNeuron() *Neuron {
	return &Neuron{0}
}

const InvalidSynapse float64 = math.MaxFloat64

func isValidSynapse(syn float64) bool {
	return syn != InvalidSynapse
}

type Net struct {
	neurons  []Neuron
	synapses [][]float64
	size     int
}

func NewNet(size int) *Net {
	syn := make([][]float64, size, size)
	for c := range syn {
		syn[c] = make([]float64, size, size)
		for r := range syn[c] {
			syn[c][r] = InvalidSynapse
		}
	}

	return &Net{
		neurons:  make([]Neuron, size, size),
		synapses: syn,
		size:     size,
	}
}

func (n *Net) Size() int {
	return n.size
}

func (n *Net) Stimulate(idx int, value float64) {
	n.neurons[idx].fire = value
}

func (n *Net) HasSynapse(from, to int) bool {
	return isValidSynapse(n.synapses[from][to])
}

func (n *Net) ValueOf(idx int) float64 {
	return n.neurons[idx].fire
}

func (n *Net) Synapse(from, to int) *float64 {
	return &n.synapses[from][to]
}

func (n *Net) ReadFrom(r io.Reader) {
	con := 0
	fmt.Fscanf(r, "%d %d", &n.size, &con)
	n = NewNet(n.size)
	for c := 0; c < con; c++ {
		from, to := 0, 0
		var value float64
		fmt.Fscanf(r, "%d -> %d = %f", &from, &to, &value)
		n.synapses[from][to] = value
	}
}

func (n *Net) WriteTo(w io.Writer) {
	connections := 0
	for c := range n.synapses {
		for r := range n.synapses[c] {
			if n.HasSynapse(c, r) {
				connections++
			}
		}
	}

	fmt.Fprintf(w, "%d %d\n", n.size, connections)
	for c := range n.synapses {
		for r := range n.synapses[c] {
			if n.HasSynapse(c, r) {
				fmt.Fprintf(w, "%d -> %d = %f\n", c, r, n.synapses[c][r])
			}
		}
	}
}

func (n *Net) Print() {
	n.WriteTo(os.Stdout)
}

func (n *Net) Clear() {
	for _, nrn := range n.neurons {
		nrn.fire = 0
	}
}

func (n *Net) Step() {
	for c := range n.synapses {
		for r := range n.synapses[c] {
			if n.HasSynapse(c, r) {
				n.neurons[r].fire += n.synapses[c][r] * n.neurons[c].fire
			}
		}
	}

	for c := range n.neurons {
		n.neurons[c].fire = 2.0 * (Sigmoid(n.neurons[c].fire) - 0.5)
	}
}
