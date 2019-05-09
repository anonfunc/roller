package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/anonfunc/roller/fn"
)

func main() {
	mux := fn.Mux()
	switch {
	case os.Getenv("AWS_EXECUTION_ENV") != "":
		lambda.Start(httpadapter.New(mux).ProxyWithContext)
	default:
		fmt.Println("http://localhost:3000/roll/2d6")
		log.Fatal(http.ListenAndServe(":3000", mux))
	}
}
