FROM golang:1.11 as builder

RUN mkdir /go/src/dental_hub
WORKDIR /go/src/dental_hub
ADD . .
RUN go get -v -u github.com/auth0/go-jwt-middleware 
RUN go get -v -u github.com/denisenkom/go-mssqldb
RUN go get -v -u github.com/go-sql-driver/mysql
RUN go get -v -u github.com/gorilla/handlers
RUN go get -v -u github.com/gorilla/mux
RUN go get -v -u gopkg.in/gomail.v2
RUN go get -v -u github.com/mongodb/mongo-go-driver/mongo
RUN go get -v -u github.com/satori/go.uuid

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/dental_hub /dental_hub/
WORKDIR /dental_hub
CMD ["./main"]

EXPOSE 4001

# ----deploy to the docker hub
# docker build -t nikolyvas/golang_dentalzone:1.0.1 -f golang.dockerfile .
# docker push nikolyvas/golang_dentalzone:1.0.1

# ----if we need to create/start it locally
# docker build -t nikolyvas/golang_dentalzone -f golang.dockerfile .
# docker run -p 4001:4001 nikolyvas/golang_dentalzone

# List all containers(only IDs)
# docker ps -aq

# List all containers
# docker container ls --all

# List all running containers
# docker container ls

# Stop all running containers
# docker stop $(docker ps -aq)

# Remove all containers
# docker rm $(docker ps -aq)

# Remove all images
# docker rmi $(docker images -q)
