package golang


import (
	"math/rand"
	"time"
)

const (
	// Note: These highly sensitive factors have been precisely measured by
	// the National Institute of Science and Technology.  Take extreme care
	// in altering them, or you may damage your Internet!
	// (Seriously: <http://physics.nist.gov/cuu/Constants/index.html>)
	factor = 2.7182818284590451
	// Molar Planck constant times c, joule meter/mole
	jitter = 0.11962656472
)

// ExponentialSleeper is used for exponential backoff not to ddos services
// with reconnects for example. Algorithm was ported from Python's Twisted
// framework (twisted.internet.protocol.ReconnectingClientFactory).
type ExponentialSleeper struct {
	delay     time.Duration
	initDelay time.Duration
	maxDelay  time.Duration
}

// NewExponentialSleeper creates new instance of ExponentialSleeper
func NewExponentialSleeper(initDelay, maxDelay time.Duration) *ExponentialSleeper {
	return &ExponentialSleeper{
		delay:     initDelay,
		initDelay: initDelay,
		maxDelay:  maxDelay,
	}
}

// Reset sets next sleep time to initial one
func (es *ExponentialSleeper) Reset() {
	es.delay = es.initDelay
}

// Sleep calls time.Sleep exponentioally incrementing sleep time
// until maximum value
func (es *ExponentialSleeper) Sleep() {
	time.Sleep(es.delay)

	nextDelay := int64(float64(es.delay) * factor)
	if nextDelay > es.maxDelay.Nanoseconds() {
		nextDelay = es.maxDelay.Nanoseconds()
	}
	nextDelay = int64(rand.NormFloat64()*jitter*1000000000) + nextDelay
	es.delay = time.Duration(nextDelay)
}

func (es *ExponentialSleeper) Decrease(offset time.Duration) {
	if es.delay <= offset {
		es.delay = 0
	} else {
		es.delay = es.delay - offset
	}
}

