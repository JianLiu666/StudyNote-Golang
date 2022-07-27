package simplecache

import (
	"math"
	"time"
)

func CalcScoreByAccessTime(weight int, lastAccessTime time.Time) float64 {
	return float64(weight) / math.Log(float64(time.Now().UnixMilli()-lastAccessTime.UnixMilli()+1))
}

func CalcScoreByWeight(weight int) float64 {
	return float64(weight / 100)
}
