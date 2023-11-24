package nginx

import (
	"sync"
)

type NginxInstance struct {
	Address string
	Port    string
}

type RequestConfig struct {
	HostHeader  string
	HostDomain  string
	HttpPath    string
	SyncTimeout int
	Retries     int
}

type NginxInstancies struct {
	Lock sync.RWMutex
	Pods map[string]NginxInstance
}

type CheckPayload struct {
	token  *string
	origin *string
}
