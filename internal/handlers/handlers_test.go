package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mauricio-mds/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"general's quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	// {"major's suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	// {"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"post search availability", "/search-availability", "POST", []postData{
	// 	{"start", "2024-01-01"}, {"end", "2024-02-01"},
	// }, http.StatusOK},
	// {"post search availability JSON", "/search-availability-json", "POST", []postData{
	// 	{"start", "2024-01-01"}, {"end", "2024-02-01"},
	// }, http.StatusOK},
	// {"post make reservation", "/make-reservation", "POST", []postData{
	// 	{"first_name", "John"}, {"last_name", "Doe"}, {"email", "test@email.com"}, {"phone", "8675309"},
	// }, http.StatusOK},	
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	srv := httptest.NewTLSServer(routes)
	defer srv.Close()

	for _, test := range tests {
		if test.method == "GET" {
			res, err := srv.Client().Get(srv.URL + test.url)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != test.expectedStatusCode {
				t.Fatalf("for %s, expected %d, but got %d", test.name, test.expectedStatusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}
			res, err := srv.Client().PostForm(srv.URL + test.url, values)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != test.expectedStatusCode {
				t.Fatalf("for %s, expected %d, but got %d", test.name, test.expectedStatusCode, res.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{ID: 1, RoomName: "General's Quarters",},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func getCtx(req *http.Request) context.Context{
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}