package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
		fmt.Fprint(w, "Failed to write output")
		return
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "test handler!")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	result := newAPIResponse(w)

	defer result.Write()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	// decode the JSON body of the request
	var loginReq loginStruct
	err := decoder.Decode(&loginReq)
	if err != nil {
		result.Status = http.StatusInternalServerError
		return
	}

	if loginReq.UserName != "MotherNight" || loginReq.Password != "BrotherCrow" {
		result.Status = http.StatusUnauthorized
		return
	}

	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["userId"] = loginReq.UserName
	claims["isAdmin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	/* Finally, write the token to the browser window */
	result.Result["token"] = tokenString
}

func bucketListHandler(w http.ResponseWriter, r *http.Request) {
	svc := newS3Session()

	vars := mux.Vars(r)
	bucket := vars["bucket"]
	bucketPrefix := bucket + "/"

	listObjectResult, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("react-stack-data"),
		Prefix: aws.String(bucketPrefix),
	})

	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	var output []string
	for _, v := range listObjectResult.Contents {
		// strip the prefix
		bucketFile := strings.Replace(*v.Key, bucketPrefix, "", 1)
		output = append(output, bucketFile)
	}

	outBytes, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(outBytes))
}

func bucketUploadHandler(w http.ResponseWriter, r *http.Request) {
	sess := newAWSSession()

	S3 := s3.New(sess)

	//uploader := s3manager.NewUploader(sess)

	vars := mux.Vars(r)
	bucket := vars["bucket"]
	child := vars["child"]

	if child != "" {
		bucket = bucket + "/" + child
	}

	uploadInput := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(r.Body),
		Bucket: aws.String("react-stack-data"),
		Key:    aws.String(bucket),
	}

	res, err := S3.PutObject(uploadInput)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, res)
	return

	//uploadConfig := &s3manager.UploadInput{
	//	Body:   r.Body,
	//	Bucket: aws.String("react-stack-data"),
	//	Key:    aws.String(bucket),
	//
	//}
	//contentType := r.Header.Get("content-type")
	//if contentType != "" {
	//	uploadConfig.ContentType = aws.String(contentType)
	//}
	//res, err := uploader.Upload(uploadConfig)
	//
	//if err != nil {
	//	fmt.Fprintln(w, err)
	//}
	//
	//fmt.Fprintln(w, res)
}

func bucketDownloadHandler(w http.ResponseWriter, r *http.Request) {
	sess := newAWSSession()

	S3 := s3.New(sess)

	//downloader := s3manager.NewDownloader(sess)

	vars := mux.Vars(r)
	bucket := vars["bucket"]
	child := vars["child"]

	if child != "" {
		bucket = bucket + "/" + child
	}

	res, err := S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("react-stack-data"),
		Key:    aws.String(bucket),
	})

	if err != nil {
		fmt.Fprintln(w, err)
	}

	io.Copy(w, res.Body)
	return

	//buf := &aws.WriteAtBuffer{}
	//_, err := downloader.Download(buf, &s3.GetObjectInput{
	//	Bucket: aws.String("react-stack-data"),
	//	Key:    aws.String(bucket),
	//
	//})
	//
	//if err != nil {
	//	fmt.Fprintln(w, err)
	//}
	//
	//io.Copy(w, bytes.NewReader(buf.Bytes()))
}
