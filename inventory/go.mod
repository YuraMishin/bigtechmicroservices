module github.com/YuraMishin/bigtechmicroservices/inventory

go 1.24.4

replace github.com/YuraMishin/bigtechmicroservices/shared => ../shared

require (
	github.com/YuraMishin/bigtechmicroservices/shared v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
)
