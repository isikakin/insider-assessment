package circuitbreaker

import (
	"math/rand"
	"time"
)

func retry(attempts int, sleep time.Duration, f func() (interface{}, error)) (interface{}, error) {

	var (
		resp interface{}
		err  error
	)

	if resp, err = f(); err != nil {
		if _, ok := err.(stop); ok {
			// Return the original error for later checking
			return resp, nil
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return resp, err
	}

	return resp, nil
}

type stop struct {
	error
}
