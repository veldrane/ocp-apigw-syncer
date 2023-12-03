package nginx

import (
	"sync"
)

type NginxInstance struct {
	Address string
	Port    string
}

type Config struct {
	Namespace   string
	Deployment  string
	HostHeader  string
	HostDomain  string
	HttpPath    string
	HttpsPort   string
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
