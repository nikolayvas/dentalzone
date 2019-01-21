package main

import (
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"

	sw "dental_hub/dentalzone"
	//migrations "dental_hub/migrations"
)

func main() {

	//migrations.MigrateAccessDb("")

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Printf("Server started")

	router := sw.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"x-xsrf-token", "authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":4001", handlers.CORS(headersOk, originsOk, methodsOk)(router)))
	//log.Fatal(http.ListenAndServeTLS(":443", "./ssl/gowebappssl.westeurope.cloudapp.azure.com.crt", "./ssl/gowebappssl.westeurope.cloudapp.azure.com.key", handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}
