package synclib

import (
	"net/http"
	"sync"
)

type NginxInstance struct {
	Address string
	Port    string
}

type Config struct {
	Namespace      string `yaml:"namespace"`
	Deployment     string `yaml:"deployment"`
	HostHeader     string `yaml:"host"`
	HostDomain     string `yaml:"domain"`
	HttpPath       string `yaml:"path"`
	HttpsPort      string `yaml:"port"`
	SyncTimeout    int    `yaml:"sync_timeout"`
	ConnTimeout    int    `yaml:"connection_timeout"`
	Retries        int    `yaml:"retries"`
	Deadline       int    `yaml:"request_deadline"`
	MaxKeepAlives  int    `yaml:"max_keepalives"`
	HostKeepAlives int    `yaml:"host_keepalives"`
}

type NginxInstancies struct {
	Lock sync.RWMutex
	Pods map[string]NginxInstance
}

type CheckPayload struct {
	authToken  *string
	origin     *string
	httpClient *http.Client
}
