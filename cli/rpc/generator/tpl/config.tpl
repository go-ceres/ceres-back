[ceres.logger.default]
	Debug=true
[ceres.client.grpc.{{.serviceName}}]
    Debug=true
[ceres.server.grpc.{{.serviceName}}]
	Host="0.0.0.0"
	Port=5201{{if .registry}}
[ceres.registry.{{.registry}}.{{.serviceName}}]
    endpoints=["127.0.0.1:2379"]{{end}}{{if .orm}}
[ceres.store.{{.orm.name}}.{{.serviceName}}]
    dns="www.baidu.com"{{end}}
