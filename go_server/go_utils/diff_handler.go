package handlers

import (
	"fmt"
	"math"
	"log"
)

func smallDiff(before, after []float32) ([]float32, float32) {
    if len(before) != len(after) {
		log.Fatal("Small Diff:Lengths of before and after are not equal")
    }
    
    vectDiff := make([]float32, len(after))
    var absDiff float32
    
    for i := range after {
        vectDiff[i] = after[i] - before[i]
        absDiff += vectDiff[i] * vectDiff[i]
    }
    
    absDiff = float32(math.Sqrt(float64(absDiff)))
	log.Println("Small Diff is working alright")
    
    return vectDiff, absDiff
}

// here we need to know how the git works under the hood?
// why? because we want to compare two very big files. things migth
//be added to it, things might be removed.
// lets re-invent the wheel but use the git's docs as a cheat sheet.
func useSmalls(allBefore, allAfter [][]float32) ([]float32) {

}
