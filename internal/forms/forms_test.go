package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form should be invalid")
	}

	postdata := url.Values{}
	postdata.Add("a", "a")
	postdata.Add("b", "b")
	postdata.Add("c", "c")

	r.PostForm = postdata
	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form should be valid")
	}

}

func TestForm_Has(t *testing.T) {
	postdata := url.Values{}
	form := New(postdata)
	if form.Has("a") {
		t.Error("form shows has field when it does not")
	}

	postdata = url.Values{}
	postdata.Add("a", "a")

	form = New(postdata)

	//check the request associate with the form
	if !form.Has("a") {
		t.Error("form doesn't have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postdata := url.Values{}
	form := New(postdata)
	form.MinLength("a", 1)
	if form.Valid() {
		t.Error("form shows minlength for non-existent field")
	}

	postdata = url.Values{}
	postdata.Add("a", "abc")

	form = New(postdata)
	if !form.MinLength("a", 1) {
		t.Error("field length larger than minlength but got false")
	}
	iserror := form.Errors.Get("a")
	if iserror != "" {
		t.Error("field doesn't has error msg but return error")
	}
	if form.MinLength("a", 5) {
		t.Error("field length smaller than minlength but got true")
	}
	iserror = form.Errors.Get("a")
	if iserror == "" {
		t.Error("field has error msg but return no error")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postdata := url.Values{}
	form := New(postdata)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}
	postdata = url.Values{}
	postdata.Add("email", "a")
	form = New(postdata)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("should not be a email but got email type")
	}
	postdata.Set("email", "abc@gmail.com")
	form = New(postdata)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("should be a email but got not email type")
	}
}
