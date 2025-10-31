package probe

import "time"

type Sample struct {
	timestamp int64
	min       float64
	max       float64
	avg       float64
	stddev    float64
	loss      float64
}

func MakeSample(min float64, max float64, avg float64, stddev float64, loss float64) Sample {
	timestamp := time.Now().Unix()
	return Sample{timestamp, min, max, avg, stddev, loss}
}

func (s Sample) Timestamp() int64 {
	return s.timestamp
}

func (s Sample) Min() float64 {
	return s.min
}

func (s Sample) Max() float64 {
	return s.max
}

func (s Sample) Avg() float64 {
	return s.avg
}

func (s Sample) StdDev() float64 {
	return s.stddev
}

func (s Sample) Loss() float64 {
	return s.loss
}
