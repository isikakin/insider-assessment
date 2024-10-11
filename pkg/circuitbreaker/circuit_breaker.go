package circuitbreaker

import (
	customHttp "assesment/pkg/http"
	"assesment/pkg/log"
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	"net/http"
	"time"
)

var breakers = make(map[string]*ClientCircuitBreakerProxy)

type ClientCircuitBreakerProxy struct {
	client customHttp.ClientDecorator
	logger log.Logger
	gb     *gobreaker.CircuitBreaker // downloaded lib structure
}

// shouldBeSwitchedToOpen checks if the circuit breaker should
// switch to the Open state
func shouldBeSwitchedToOpen(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 10 && failureRatio >= 0.6
}

func (c *ClientCircuitBreakerProxy) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// We call the Execute method and wrap our client's call
	resp, err := c.gb.Execute(func() (interface{}, error) {
		return retry(3, time.Second, func() (interface{}, error) {
			requestCtx, cancel := context.WithTimeout(ctx, time.Minute*1)
			resp, err := c.client.Do(requestCtx, req)

			defer cancel()

			if err != nil {
				// This error will result in a retry
				c.logger.Exception(ctx, fmt.Sprintf("%s retry err.", c.gb.Name()), err)
				return nil, err
			}
			defer resp.Body.Close()

			s := resp.StatusCode
			switch {
			case s >= 500:
				// Retry
				c.logger.Exception(ctx, fmt.Sprintf("%s interval server error.", c.gb.Name()), err)
				return resp, fmt.Errorf("server error: %v", s)
			case s >= http.StatusBadRequest:
				// Don't retry, it was client's fault
				c.logger.Exception(ctx, fmt.Sprintf("%s client err.", c.gb.Name()), err)
				return resp, stop{fmt.Errorf("client error: %v", s)}
			default:
				c.logger.Info(ctx, fmt.Sprintf("%s request success", c.gb.Name()))
				return resp, nil
			}
		})
	})

	if resp == nil {
		return nil, err
	}

	return resp.(*http.Response), err
}

func NewClientCircuitBreakerProxy(ctx context.Context, name string, logger log.Logger, client customHttp.ClientDecorator) *ClientCircuitBreakerProxy {

	if existBreaker, ok := breakers[name]; ok {
		return existBreaker
	}

	// We need circuit breaker configuration
	cfg := gobreaker.Settings{
		Name: name,
		// When to flush counters int the Closed state
		Interval: 5 * time.Second,
		// Time to switch from Open to Half-open
		Timeout: 10 * time.Second,
		// Function with check when to switch from Closed to Open
		ReadyToTrip: shouldBeSwitchedToOpen,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Handler for every state change. We'll use for debugging purpose
			logger.Info(ctx, fmt.Sprintf("%s breaker state changed from %s to %s", name, from.String(), to.String()))
		},
	}

	circuitBreakerProxy := &ClientCircuitBreakerProxy{
		client: client,
		logger: logger,
		gb:     gobreaker.NewCircuitBreaker(cfg),
	}

	breakers[name] = circuitBreakerProxy

	return circuitBreakerProxy
}
