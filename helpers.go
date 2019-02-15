package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func newAWSSession() *session.Session {

	sess := session.Must(session.NewSession())
	return sess

}

func newS3Session() *s3.S3 {
	sess := newAWSSession()

	svc := s3.New(sess)
	return svc
}
