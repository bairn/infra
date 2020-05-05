package lb

import (
	"math/rand"
	"sync/atomic"
)

var _ Balancer = new(RoundRobinBalancer)

type RoundRobinBalancer struct{
	ct uint32
}

func (r *RoundRobinBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}

	count := atomic.AddUint32(&r.ct, 1)

	index := int(count) % len(hosts)

	instance := hosts[index]

	return instance
}

type RandomBalancer struct {

}

func (r *RandomBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}

	count := rand.Int31()
	index := int(count) % len(hosts)
	instance := hosts[index]
	return instance
}