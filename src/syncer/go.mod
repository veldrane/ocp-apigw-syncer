module syncer

go 1.20

require ( 
	goa.design/goa/v3 v3.13.2
	github.com/nginx v0.0.0
	bitbucket.org/veldrane/golibs/ocp4cli v0.0.0
)
require (
	github.com/AnatolyRugalev/goregen v0.1.0 // indirect
	github.com/dimfeld/httppath v0.0.0-20170720192232-ee938bf73598 // indirect
	github.com/go-chi/chi/v5 v5.0.10 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/manveru/faker v0.0.0-20171103152722-9fbc68a78c4d // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/nginx => ./local/nginx
	bitbucket.org/veldrane/golibs/ocp4cli => ./local/ocp4cli
)
