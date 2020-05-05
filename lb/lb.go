package lb


type Balancer interface {
	Next(key string, hosts []*ServerInstance) *ServerInstance
}