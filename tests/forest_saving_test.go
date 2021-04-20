package test

import (
	"fmt"
	"os"
	"testing"

	randomforest "github.com/MandelV/randomForest/v2"
	"github.com/MandelV/randomForest/v2/internal"
	"github.com/MandelV/randomForest/v2/tests/generator"
)

func TestSaving(t *testing.T) {
	n := 100
	features := 20
	classes := 4
	trees := 1000
	forest := randomforest.Forest{}
	data, res := generator.CreateDataset(n, features, classes)
	forestData := randomforest.ForestData{X: data, Class: res}
	forest.Data = forestData
	forest.Train(trees)

	if fileName, err := forest.Save("saved/", false); err != nil {
		t.Error(err)
	} else {
		println(fileName)

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			t.Error(err)
		} else {

			if errRemove := os.Remove(fileName); errRemove != nil {
				t.Error(errRemove)
			}
		}
	}
}

func TestToBytes(t *testing.T) {
	n := 100
	features := 20
	classes := 4
	trees := 1000
	forest := randomforest.Forest{}
	data, res := generator.CreateDataset(n, features, classes)
	forestData := randomforest.ForestData{X: data, Class: res}
	forest.Data = forestData
	forest.Train(trees)

	if bytes, err := internal.ToBytes(&forest, false); err != nil {
		t.Error(err)
	} else {
		if len(bytes) == 0 {
			t.Fail()
		}
	}
}

func TestLoading(t *testing.T) {

	var forest *randomforest.Forest = nil
	var errForest error

	if forest, errForest = randomforest.LoadForest("saved/test.bin"); errForest != nil {
		t.Error(errForest)
		return
	}

	results := forest.Vote([]float64{0.1, 0.1, 0.1, 0.1})
	if len(results) != 4 {
		t.Error("Error with vote")
	}
}

func TestByteToForest(t *testing.T) {

	var forestByte []byte
	var errToByte error
	if forest, errForest := randomforest.LoadForest("saved/test.bin"); errForest != nil {
		t.Error(errForest)
		return
	} else {
		if forestByte, errToByte = internal.ToBytes(&forest, true); errToByte != nil {
			t.Error(errToByte)
			return
		}

		forest, err := randomforest.BytesToForest(forestByte)
		if err != nil {
			t.Error(err)
			return
		}

		results := forest.Vote([]float64{0.1, 0.1, 0.1, 0.1})

		fmt.Println(results)
		if len(results) != 4 {
			t.Error("Error with vote")
		}
	}

}
