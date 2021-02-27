package randomforest

import (
	"math/rand"
	"os"
	"testing"
)

func TestSaving(t *testing.T) {
	xData := [][]float64{}
	yData := []int{}
	for i := 0; i < 1000; i++ {
		x := []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		y := int(x[0] + x[1] + x[2] + x[3])
		xData = append(xData, x)
		yData = append(yData, y)
	}
	forest := &Forest{}

	forestData := ForestData{X: xData, Class: yData}
	forest.Data = forestData
	forest.Train(1000)

	if fileName, err := forest.Save("saved/"); err != nil {
		t.Error(err)
	} else {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			t.Error(err)
		} else {

			if errRemove := os.Remove(fileName); errRemove != nil {
				t.Error(errRemove)
			}
		}
	}

}

func TestLoading(t *testing.T) {

	var forest *Forest = nil
	var errForest error

	if forest, errForest = Load("saved/forestTest.bin"); errForest != nil {
		t.Error(errForest)
		return
	}

	results := forest.Vote([]float64{0.1, 0.1, 0.1, 0.1})

	if len(results) != 4 {
		t.Error("Error with vote")
	}
}
