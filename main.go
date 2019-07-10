package main

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var muxLambda *gorillamux.GorillaMuxAdapter

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	proxyRes, err := muxLambda.Proxy(request)
	proxyRes.MultiValueHeaders["Access-Control-Allow-Origin"] = []string{"*"}
	proxyRes.MultiValueHeaders["Access-Control-Allow-Headers"] = []string{"*"}
	proxyRes.MultiValueHeaders["Access-Control-Allow-Methods"] = []string{"*"}
	return proxyRes, err
}

func init() {
	log.Println("Mux cold start")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", loginHandler).Methods("POST", "PUT")
	router.Handle("/{bucket}", jwtMiddleware.Handler(http.HandlerFunc(bucketUploadHandler))).Methods("POST", "PUT")
	router.HandleFunc("/{bucket}", bucketDownloadHandler).Methods("GET")
	router.HandleFunc("/{bucket}/list", bucketListHandler).Methods("GET")
	router.Handle("/{bucket}/{child}", jwtMiddleware.Handler(http.HandlerFunc(bucketUploadHandler))).Methods("POST", "PUT")
	router.HandleFunc("/{bucket}/{child}", bucketDownloadHandler).Methods("GET")
	router.HandleFunc("/", handler)

	muxLambda = gorillamux.New(router)
}

func main() {
	lambda.Start(Handler)
}
