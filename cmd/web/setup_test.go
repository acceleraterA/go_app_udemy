package main

import (
	"net/http"
	"os"
	"testing"
)

//run before the tests run

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

// create a type myHandler that has a method ServeHTTP
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
