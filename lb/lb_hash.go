package lb

import "hash/crc32"

var _ Balancer = new(HashBalancer)

type HashBalancer struct {

}

func (h *HashBalancer) Next(key string, hosts[]*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}

	count := crc32.ChecksumIEEE([]byte(key))

	index := int(count) % len(hosts)

	instance := hosts[index]

	return instance
}
