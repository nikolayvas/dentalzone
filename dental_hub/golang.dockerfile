FROM golang:1.9 as builder
RUN mkdir /go/src/dental_hub
WORKDIR /go/src/dental_hub
ADD . .
RUN go get -v -u github.com/auth0/go-jwt-middleware 
RUN go get -v -u github.com/denisenkom/go-mssqldb
RUN go get -v -u github.com/go-sql-driver/mysql
RUN go get -v -u github.com/gorilla/handlers
RUN go get -v -u github.com/gorilla/mux
RUN go get -v -u gopkg.in/gomail.v2

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /go/src/dental_hub /dental_hub/
WORKDIR /dental_hub
CMD ["./main"]

EXPOSE 4001

# docker build -t golang_dentalzone -f golang.dockerfile .
# docker run -p 4001:4001 golang_dentalzone