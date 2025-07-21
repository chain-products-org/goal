package limiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketStart(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		// max 2 tokens, only consume 1 token every 2 seconds
		limiter := CreateAndStartTokenBucketLimiter(2, 1, time.Second*2)
		defer limiter.Stop()

		time.Sleep(time.Second * 5)
		assert.Equal(t, 2, len(limiter.tokens))
	})
}

func TestTokenBucketTryTake(t *testing.T) {
	limiter := CreateAndStartTokenBucketLimiter(2, 1, time.Second*2)
	defer limiter.Stop()

	time.Sleep(time.Second * 5)

	if !limiter.TryTake() {
		t.Error("should be able to take the first token")
	}
	if !limiter.TryTake() {
		t.Error("should be able to take the second token")
	}
	if limiter.TryTake() {
		t.Error("bucket is empty, should not be able to take the third token")
	}

	time.Sleep(time.Second * 5)
	if !limiter.TryTake() {
		t.Error("should be able to take the new token after 5 seconds")
	}
	if !limiter.TryTake() {
		t.Error("bucket is empty, should not be able to take the third token")
	}
	if limiter.TryTake() {
		t.Error("bucket is empty, should not be able to take the third token")
	}
}

func TestTokenBucketTake(t *testing.T) {
	limiter := CreateAndStartTokenBucketLimiter(1, 1, time.Second)

	time.Sleep(time.Second * 3)

	limiter.Take() // ok
	t.Log("take the first token")
	limiter.Take() // ok
	t.Log("take the second token")

	go func() {
		time.AfterFunc(time.Second*4, func() {
			t.Log("stop limiter after 4 seconds")
			limiter.Stop()
		})
	}()
	t.Log("taking the third token, should be blocked")
	limiter.Take() // blocked
}

func TestTokenBucketTakeWithTimeout(t *testing.T) {
	limiter := CreateAndStartTokenBucketLimiter(1, 1, time.Second*2)
	defer limiter.Stop()

	start := time.Now()
	limiter.TakeWithTimeout(time.Second * 3) // ok
	t.Log("take the first token")
	limiter.TakeWithTimeout(time.Second * 2) // ok
	t.Log("take the second token")
	used := time.Since(start)
	t.Logf("take 2 tokens, used %s", used)
	// need 4 seconds to take 2 tokens
	if used < time.Second*4 {
		t.Errorf("expect 4 seconds, actual %s", used)
	}
}

func TestTokenBucketLimiterConcurrent(t *testing.T) {
	// create a token bucket for testing: capacity 1, 1 token per second
	limiter := CreateAndStartTokenBucketLimiter(1, 1, time.Second)
	defer limiter.Stop()

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	time.Sleep(time.Second * 2)

	// start 10 concurrent requests
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.TryTake() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// since the capacity is 1, only 1 request should succeed
	if successCount != 1 {
		t.Errorf("1 token per second, expect 1 success request, actual %d", successCount)
	}
}

func TestTokenBucketLimiterStats(t *testing.T) {
	limiter := CreateAndStartTokenBucketLimiter(1, 1, time.Second)
	defer limiter.Stop()

	// reset stats
	limiter.ResetStats()

	time.Sleep(time.Second * 2)

	// try to take tokens
	limiter.TryTake() // success
	limiter.TryTake() // failed
	limiter.TryTake() // failed

	total, blocked, successRate := limiter.GetStats()

	if total != 3 {
		t.Errorf("expect 3 total requests, actual %d", total)
	}
	if blocked != 2 {
		t.Errorf("expect 2 blocked requests, actual %d", blocked)
	}
	if successRate != 33.33333333333333 {
		t.Errorf("expect 33.33%% success rate, actual %f%%", successRate)
	}
}
