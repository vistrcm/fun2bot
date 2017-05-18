package randomizer

import (
	"math"
	"math/rand"
	"time"
)

// round float number to the nearest integer
func round(f float64) int {
	return int(math.Floor(f + .5))
}

// afterMinutes returns time after random number of minutes around 24 hours.
// it is using normal distribution with mean = 24 and StdDev = 6. Assumes hours. Translate it to the minutes.
func GetNextTime() time.Time {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	desiredStdDev := 360.0 // manually tuned deviation parameter to make distribution wide enough
	desiredMean := 1440.0  // 24 hours is 1440 minutes.

	n := round(r.NormFloat64()*desiredStdDev + desiredMean) // calculate randomized time delay
	return time.Now().Add(time.Duration(n) * time.Minute)   // return duration
}
