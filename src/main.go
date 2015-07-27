package main

import (
	neural "./neural"
	"fmt"
	"math"
	// "math/rand"
	"sort"
)

var xorInp [][]int = [][]int{
	{-1, -1, -1},
	{-1, 1, 1},
	{1, -1, 1},
	{1, 1, -1},
}

func testNet(net *neural.Net) {
	for c := range xorInp {
		net.Stimulate(0, float64(xorInp[c][0]))
		net.Stimulate(1, float64(xorInp[c][1]))

		net.Step()

		fmt.Printf("INPUT 1: %d\tINPUT 2: %d\tEXPECTED: %d\tOUTPUT: %f\n", xorInp[c][0], xorInp[c][1], xorInp[c][2], net.ValueOf(6))
	}
}

type NetGrade struct {
	net   *neural.Net
	grade float64
}

type NetGrades []NetGrade

func (grades NetGrades) Len() int {
	return len(grades)
}

func (grades NetGrades) Less(c, r int) bool {
	return grades[c].grade > grades[r].grade
}

func (grades NetGrades) Swap(c, r int) {
	grades[c], grades[r] = grades[r], grades[c]
}

func gradeNets(nets []*neural.Net) NetGrades {
	grades := make(NetGrades, len(nets), len(nets))

	reps := 20

	var maxDeviation float64 = float64(len(xorInp)) * float64(reps) * 2.0
	for idx, net := range nets {
		var deviation float64 = 0

		for r := 0; r < reps; r++ {
			for _, inp := range xorInp {
				net.Stimulate(0, float64(inp[0]))
				net.Stimulate(1, float64(inp[1]))

				net.Step()

				deviation += math.Abs((float64(inp[2]) + 1.0) - (net.ValueOf(6) + 1.0))
			}

			grades[idx] = NetGrade{
				net:   net,
				grade: deviation / maxDeviation,
			}
		}
	}

	return grades
}

func mutateNets(grades NetGrades) {
	sort.Sort(grades)

	randNet := func() *neural.Net {
		return grades[int(neural.RandMax(float64(len(grades))))].net
	}

	best := grades[0].net

	for c := 0; c < len(grades)/2; c++ {
		grades[c].net = neural.Cross(best, randNet())

		if neural.Chance(0.05) {
			grades[c].net.Mutate(0.33)
		}
	}

	for c := 0; c <= len(grades)/5; c++ {
		grades[c+len(grades)/2].net.Mutate(0.5)
	}
}

func xorNets(cnt int) *neural.Net {
	nets := make([]*neural.Net, cnt, cnt)
	for c := range nets {
		nets[c] = neural.NewNet(7)

		// input 0 - to hidden
		*nets[c].Synapse(0, 2) = 0.0
		*nets[c].Synapse(0, 3) = 0.0
		*nets[c].Synapse(0, 4) = 0.0
		*nets[c].Synapse(0, 5) = 0.0

		// input 1 - to hidden
		*nets[c].Synapse(1, 2) = 0.0
		*nets[c].Synapse(1, 3) = 0.0
		*nets[c].Synapse(1, 4) = 0.0
		*nets[c].Synapse(1, 5) = 0.0

		// hidden to output
		*nets[c].Synapse(2, 6) = 0.0
		*nets[c].Synapse(3, 6) = 0.0
		*nets[c].Synapse(4, 6) = 0.0
		*nets[c].Synapse(5, 6) = 0.0

		nets[c].Randomize()
	}
	lastBest := 0.
	var maxErr float64 = 0.01
	for {
		grades := gradeNets(nets)

		bestNet := grades[0].net
		bestGrade := grades[0].grade

		for _, grade := range grades {
			if grade.grade < bestGrade {
				bestGrade = grade.grade
				bestNet = grade.net
			}
		}
		// testNet(bestNet)
		if lastBest != bestGrade {
			fmt.Println(bestGrade)
			testNet(bestNet)
			fmt.Println("------------------------------------")
			lastBest = bestGrade
		}
		if bestGrade < maxErr {
			return bestNet
		}

		mutateNets(grades)
	}

	return nil
}

func main() {
	testNet(xorNets(20))
}
