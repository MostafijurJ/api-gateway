package proxy

import (
	"sync/atomic"
)

type roundRobin struct{ idx uint64 }

func (r *roundRobin) Next(n int) int {
	if n <= 0 {
		return 0
	}
	i := atomic.AddUint64(&r.idx, 1)
	return int(i % uint64(n))
}

type leastConnections struct{}

func (l *leastConnections) Next(conns []int) int {
	if len(conns) == 0 {
		return 0
	}
	minIdx, minVal := 0, conns[0]
	for i := 1; i < len(conns); i++ {
		if conns[i] < minVal {
			minIdx, minVal = i, conns[i]
		}
	}
	return minIdx
}
