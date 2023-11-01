package syncer

type NginxInstance struct {
	Hostname string
	Address  string
	Port     string
}

func New() []NginxInstance {
	res := make([]NginxInstance, 0)
	return res
}
