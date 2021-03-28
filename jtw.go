package main

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"net/http"
)

/* Set up a global string for our secret */
var mySigningKey = []byte("Jormungandr42LepersAtrophy")

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func getClaims(r *http.Request) jwt.MapClaims {
	user := r.Context().Value("user")
	return user.(*jwt.Token).Claims.(jwt.MapClaims)
}
