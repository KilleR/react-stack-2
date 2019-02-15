package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"io"
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
	svc := newS3Session()

	listObjectResult, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("react-stack-data"),
	})

	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	fmt.Fprintln(w, listObjectResult)
}

func bucketUploadHandler(w http.ResponseWriter, r *http.Request) {
	sess := newAWSSession()

	uploader := s3manager.NewUploader(sess)

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	res, err := uploader.Upload(&s3manager.UploadInput{
		Body:   r.Body,
		Bucket: aws.String(bucket),
	})

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, res)
}

func bucketDownloadHandler(w http.ResponseWriter, r *http.Request) {
	sess := newAWSSession()

	downloader := s3manager.NewDownloader(sess)

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	buf := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		fmt.Fprintln(w, err)
	}

	io.Copy(w, bytes.NewReader(buf.Bytes()))
}
