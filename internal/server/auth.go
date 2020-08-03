package server

import (
	"fmt"
	"net/http"
)

var admintoken, usertoken string

const AuthorizationHeader = "Authorization"
const UnauthorizedErrorMessage = "Отказано в доступе"

func authAdmin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(AuthorizationHeader) == admintoken {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, UnauthorizedErrorMessage, http.StatusUnauthorized)
		}
	})
}

func authUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(AuthorizationHeader) == usertoken {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, UnauthorizedErrorMessage, http.StatusUnauthorized)
		}
	})
}

func DefineParams(adminparam, userparam string) {
	admintoken = adminparam
	usertoken = userparam
	fmt.Printf("Params defined %s\t%s\n", admintoken, usertoken)
}

/*
func init() {
	var parser = flags.NewParser(&envs, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("got auth params: %v\n", envs)
}
*/
