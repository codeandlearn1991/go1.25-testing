package concurrency_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/codeandlearn1991/go1.25-testing/concurrency"
)

func TestSleepLegacy(t *testing.T) {
	start := time.Now()
	concurrency.Sleep()
	t.Log(time.Since(start))
	t.Fail()
}

func TestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		concurrency.Sleep()
		t.Log(time.Since(start))
		t.Fail()
	})
}

func TestWait(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := false
		go func() {
			done = true
		}()
		synctest.Wait()
		t.Log(done)
		t.Fail()
	})
}
