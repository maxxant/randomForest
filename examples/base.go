package main

import (
  "fmt"
  "math/rand"

  randomForest "github.com/malaschitz/randomForest"
)

func main() {
  xData := [][]float64{}
  yData := []int{}
  for i := 0; i < 1000; i++ {
    x := []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
    y := int(x[0] + x[1] + x[2] + x[3])
    xData = append(xData, x)
    yData = append(yData, y)
    // fmt.Println(x, y)
  }

  forest := randomForest.Forest{}
  forest.Data = randomForest.ForestData{X: xData, Class: yData}
  forest.Train(1000)
  //test
  fmt.Println("randomForest")
  fmt.Println("Vote ", forest.Vote([]float64{0.1, 0.1, 0.1, 0.1}))
  // fmt.Println("VoteW", forest.WeightVote([]float64{0.1, 0.1, 0.1, 0.1}))
  fmt.Println("Vote ", forest.Vote([]float64{0.9, 0.9, 0.9, 0.9}))
  // fmt.Println("VoteW", forest.WeightVote([]float64{0.9, 0.9, 0.9, 0.9}))
  forest.PrintFeatureImportance()

  if false {
    fmt.Println("deepForest")
    dForest := forest.BuildDeepForest()
    dForest.Train(100, 100, 100)
    fmt.Println("Vote", dForest.Vote([]float64{0.1, 0.1, 0.1, 0.1}))
    // fmt.Println("Vote", dForest.Vote([]float64{0.5, 0.4, 0.6, 0.8}))
    fmt.Println("Vote", dForest.Vote([]float64{0.9, 0.9, 0.9, 0.9}))
  }
}
