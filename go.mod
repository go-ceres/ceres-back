module github.com/go-ceres/ceres

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/getkin/kin-openapi v0.81.0
	github.com/go-ceres/cli/v2 v2.2.2
	github.com/go-ceres/go-ceres v0.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/google/wire v0.5.0
	github.com/gookit/gcli/v3 v3.0.0
	github.com/jhump/protoreflect v1.10.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20210701191553-46259e63a0a9
	google.golang.org/protobuf v1.27.1
	vitess.io/vitess v0.12.0
)

replace github.com/go-ceres/go-ceres v0.0.0 => ../go-ceres
