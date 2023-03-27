package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	//define the input var myHandler type
	var myH myHandler
	h := NoSurf(&myH)

	//check if the return value is type http.Handler
	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
