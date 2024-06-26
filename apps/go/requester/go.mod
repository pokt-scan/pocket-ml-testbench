module requester

go 1.22.2

require (
	github.com/pokt-foundation/pocket-go v0.17.0
	github.com/puzpuzpuz/xsync/v3 v3.1.0
	github.com/rs/zerolog v1.32.0
	github.com/stretchr/testify v1.9.0
	go.mongodb.org/mongo-driver v1.15.0
	go.temporal.io/api v1.32.0
	go.temporal.io/sdk v1.26.1
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	packages/logger v0.0.0-00010101000000-000000000000
	packages/mongodb v0.0.0-00010101000000-000000000000
	packages/pocket_rpc v0.0.0-00010101000000-000000000000
	packages/utils v0.0.0-00010101000000-000000000000
)

replace packages/logger => ./../../../packages/go/logger

replace packages/mongodb => ./../../../packages/go/mongodb

replace packages/utils => ./../../../packages/go/utils

replace packages/pocket_rpc => ./../../../packages/go/pocket_rpc

require (
	github.com/alitto/pond v1.8.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.7.0-rc.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pokt-foundation/utils-go v0.7.0 // indirect
	github.com/puzpuzpuz/xsync v1.5.2 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
