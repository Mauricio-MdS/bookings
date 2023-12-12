package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mauricio-mds/bookings/internal/config"
	"github.com/mauricio-mds/bookings/internal/models"
)

type mockWriter struct{}

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	//what I am going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	pathToTemplates = "./../../templates"

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

func (tw *mockWriter) Header() http.Header{
	return http.Header{}
}

func (tw *mockWriter) Write(b []byte) (int, error){
	return len(b), nil
}

func (tw *mockWriter) WriteHeader(i int) {
	
}