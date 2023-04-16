package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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
	//{"rs", "/reservation-summary", "GET", http.StatusOK},

	// {"sap", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "e2023-03-23"},
	// 	{key: "end", value: "e2023-03-24"},
	// }, http.StatusOK},
	// {"sajp", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "e2023-03-23"},
	// 	{key: "end", value: "e2023-03-24"},
	// }, http.StatusOK},
	// {"mrp", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "John@gmail.com"},
	// 	{key: "phone", value: "1234543545"},
	// }, http.StatusOK},
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

func TestRepository_ReservationSummary(t *testing.T) {

}

// getCtx returns the ctx with header
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
