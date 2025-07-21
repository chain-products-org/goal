package limiter

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TokenBucketLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	stop   chan struct{}

	// statistics
	mu              sync.Mutex
	totalRequests   int64
	blockedRequests int64
	lastResetTime   time.Time
}

// Start will start the token bucket limiter
func (l *TokenBucketLimiter) Start() {
	go func() {
		for {
			select {
			case <-l.ticker.C:
				select {
				case l.tokens <- struct{}{}:
					// add token
				default:
					// drop token
				}
			case <-l.stop:
				return
			}
		}
	}()
}

// registerExitHandler will register the exit handler
func (l *TokenBucketLimiter) RegisterExitHandler() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		l.Stop()
		os.Exit(0)
	}()
}

// Stop will stop the token bucket limiter
func (l *TokenBucketLimiter) Stop() {
	l.ticker.Stop()
	close(l.stop)
}

// TryTake will try to take a token, non-blocking
func (l *TokenBucketLimiter) TryTake() bool {
	l.mu.Lock()
	l.totalRequests++
	l.mu.Unlock()

	select {
	case <-l.tokens:
		return true
	default:
		l.mu.Lock()
		l.blockedRequests++
		l.mu.Unlock()
		return false
	}
}

// Take will take a token, block until a token is available
func (l *TokenBucketLimiter) Take() {
	<-l.tokens
}

// TakeWithTimeout will take a token, block until a token is available or timeout
func (l *TokenBucketLimiter) TakeWithTimeout(timeout time.Duration) bool {
	select {
	case <-l.tokens:
		return true
	case <-time.After(timeout):
		return false
	}
}

// GetStats will get the statistics
func (l *TokenBucketLimiter) GetStats() (total, blocked int64, successRate float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.totalRequests > 0 {
		successRate = float64(l.totalRequests-l.blockedRequests) / float64(l.totalRequests) * 100
	}

	return l.totalRequests, l.blockedRequests, successRate
}

// ResetStats will reset the statistics
func (l *TokenBucketLimiter) ResetStats() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.totalRequests = 0
	l.blockedRequests = 0
	l.lastResetTime = time.Now()
}

// NewTokenBucketLimiter will create a new token bucket limiter
//
// - capacity: the capacity of the bucket
// - rate: the number of tokens generated per window
// - window: the time window
func NewTokenBucketLimiter(capacity int, rate int, window time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		tokens: make(chan struct{}, capacity),
		ticker: time.NewTicker(window / time.Duration(rate)),
		stop:   make(chan struct{}),
	}
}

// CreateAndStartTokenBucketLimiter will create a new token bucket limiter and start it
func CreateAndStartTokenBucketLimiter(capacity int, rate int, window time.Duration) *TokenBucketLimiter {
	limiter := NewTokenBucketLimiter(capacity, rate, window)
	limiter.Start()
	limiter.RegisterExitHandler()
	return limiter
}
