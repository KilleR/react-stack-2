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
	// default response
	//index, err := ioutil.ReadFile("public/index.html")
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//}
	//
	//return events.APIGatewayProxyResponse{
	//	StatusCode: 200,
	//	Body:       string(index),
	//	Headers: map[string]string{
	//		"Content-Type": "text/html",
	//	},
	//}, nil
	return muxLambda.Proxy(request)
}

func init() {
	log.Println("Mux cold start")
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/", handler)

	muxLambda = gorillamux.New(router)
}

func main() {
	lambda.Start(Handler)
}
