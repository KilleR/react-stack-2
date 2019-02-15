package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to read index: "+err.Error())
		return
	}

	_, err = fmt.Fprint(w, string(index))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Faiiled to write output")
		return
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "test handler!")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func bucketListHandler(w http.ResponseWriter, r *http.Request) {
	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Fprintln(w, aerr.Error())
		} else {
			fmt.Fprintln(w, err.Error())
		}
	}

	fmt.Fprintln(w, result)
}
