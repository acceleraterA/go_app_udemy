package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/acceleraterA/go_app_udemy/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	r := getRoutes()
	ts := httptest.NewTLSServer(r)
	//when current function is finished, defer got executed
	defer ts.Close()
	for _, e := range theTests {

		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//new req
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	//add header to ctx and ass ctx to req
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	//test Reservation func by casting it into handlerFunc
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session, doesn't have session.put (reset everthing)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	//modify the roomid in reservation to an non-existent number
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2023-05-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-05-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jianci")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//new req
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// test for invalid start date
	reqBody = "start_date="
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-05-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jianci")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	reqBody = "start_date=2023-05-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jianci")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid roomID
	reqBody = "start_date=2023-05-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-05-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jianci")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid roomID: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// test for invalid form
	reqBody = "start_date=2023-05-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-05-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid roomID: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert Reservations to database
	reqBody = "start_date=2023-05-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-05-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jct")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for not inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert restriction into database
	reqBody = "start_date=2023-05-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-05-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jct")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Tan")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=123@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=11")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for not inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	//case 1 - rooms are not available

	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	//create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	//set hte request header
	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	//get respose recorder
	rr := httptest.NewRecorder()
	//make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	var j jsonResponse

	err := json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("faild to parse json")
	}
	// since we specified a start date > 2049-12-31, we expect no availability
	if j.Ok {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	//case 2 - rooms are available

	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//create request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	//set hte request header
	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	//get respose recorder
	rr = httptest.NewRecorder()
	//make handler handlerfunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("faild to parse json")
	}
	// since we specified a start date > 2049-12-31, we expect no availability
	if !j.Ok {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}

	//case 3 - no request body
	//create request
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	//set hte request header
	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	//get respose recorder
	rr = httptest.NewRecorder()
	//make handler handlerfunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("faild to parse json")
	}
	// response for err parseform
	if j.Ok || j.Message != "Internal server error" {
		t.Error("Got availability when request body was empty")
	}

	//case 4 - database error
	//create request
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	//set hte request header
	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	//get respose recorder
	rr = httptest.NewRecorder()
	//make handler handlerfunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("faild to parse json")
	}
	// response for err parseform
	if j.Ok || j.Message != "Error connecting to database" {
		t.Error("Got availability when simulating database error")
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	//create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	//add header to ctx and ass ctx to req
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	//set hte request header
	req.Header.Set("content-Type", "application/x-www-form-urlencoded")
	//get respose recorder
	rr := httptest.NewRecorder()
	//make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)
}

// bookRoomTests is the data for the BookRoom handler tests
var bookRoomTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{
		name:               "database-works",
		url:                "/book-room?s=2050-01-01&e=2050-01-02&id=1",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "database-fails",
		url:                "/book-room?s=2040-01-01&e=2040-01-02&id=3",
		expectedStatusCode: http.StatusSeeOther,
	},
}

// TestBookRoom tests the BookRoom handler
func TestBookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	for _, e := range bookRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		session.Put(ctx, "reservation", reservation)

		handler := http.HandlerFunc(Repo.BookRoom)

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("%s failed: returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`"action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	//range through all the tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")
		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		if e.expectedLocation != "" {

			//get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())

			}
			// checking for expected values in HTML
			if e.expectedHTML != "" {
				//read the response body into a string
				html := rr.Body.String()
				if !strings.Contains(html, e.expectedHTML) {
					t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
				}

			}
		}
	}
}

// getCtx returns the ctx with header
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
