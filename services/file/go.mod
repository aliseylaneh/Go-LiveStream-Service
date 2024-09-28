module vpeer_file

go 1.21.3

require (
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	safir/libs/appconfigs v0.0.0-00010101000000-000000000000
	safir/libs/appstates v0.0.0-00010101000000-000000000000
	safir/libs/idgen v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/rs/xid v1.5.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
)

replace safir/libs/appconfigs => ../../libs/appconfigs

replace safir/libs/appstates => ../../libs/appstates

replace safir/libs/idgen => ../../libs/idgen
