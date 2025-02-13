# MICROSERVICES - VIDEO CATALOG SERVICE

Video Catalog Service for the [Microservices](https://github.com/SagarMaheshwary/microservices) project.

### OVERVIEW

- Golang
- ZeroLog
- gRPC – Acts as both the main server and client for the User service
- PostgreSQL – Stores video metadata
- GORM – ORM (Object-Relational Mapper) for PostgreSQL
- RabbitMQ – Enables asynchronous communication with the [encode service](https://github.com/SagarMaheshwary/microservices-encode-service) for storing video metadata after encoding
- Prometheus Client – Exports default and custom metrics for Prometheus server monitoring
- AWS S3 & CloudFront – Retrieves video thumbnails from S3 and generates CloudFront URLs for video streaming

### SETUP

Follow the instructions in the [README](https://github.com/SagarMaheshwary/microservices?tab=readme-ov-file#setup) of the main microservices repository to run this service along with others using Docker Compose.

### APIs (gRPC)

Proto files are located in the **internal/proto** directory.

| SERVICE                                                        | RPC      | BODY          | METADATA | DESCRIPTION                                                                                      |
| -------------------------------------------------------------- | -------- | ------------- | -------- | ------------------------------------------------------------------------------------------------ |
| VideoCatalogService                                            | FindAll  | -             | -        | List videos                                                                                      |
| VideoCatalogService                                            | FindById | {"id": "int"} | -        | Get specified video details as well as DASH manifest url from cloudfront for streaming the video |
| [Health](https://google.golang.org/grpc/health/grpc_health_v1) | Check    | -             | -        | Service health check                                                                             |

### APIs (REST)

| API      | METHOD | BODY | Headers | Description                 |
| -------- | ------ | ---- | ------- | --------------------------- |
| /metrics | GET    | -    | -       | Prometheus metrics endpoint |

### RABBITMQ MESSAGES

#### Received Messages (Consumed from the Queue)

| MESSAGE NAME           | RECEIVED FROM                                                                     | DESCRIPTION                                               |
| ---------------------- | --------------------------------------------------------------------------------- | --------------------------------------------------------- |
| VideoEncodingCompleted | [Encode Service](https://github.com/SagarMaheshwary/microservices-encode-service) | Handles storing video metadata after encoding is complete |
