package main

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var muxLambda *gorillamux.GorillaMuxAdapter

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.Proxy(request)
}

func init() {
	log.Println("Mux cold start")
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/list", bucketListHandler)
	router.HandleFunc("/{bucket}", bucketUploadHandler).Methods("POST", "PUT")
	router.HandleFunc("/{bucket}", bucketDownloadHandler).Methods("GET")
	router.HandleFunc("/", handler)

	muxLambda = gorillamux.New(router)
}

func main() {
	lambda.Start(Handler)
}
