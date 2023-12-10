package render

import (
	"net/http"
	"testing"

	"github.com/mauricio-mds/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	res, err := getSession()
	failIfError(err, t)
	session.Put(res.Context(), "flash", "test-flash")
	result := AddDefaultData(&td, res)
	if result.Flash != "test-flash" {
		t.Errorf("expected flash value of test-flash and got %s", result.Flash)
	}
}

func TestCreateTemplateCache(t *testing.T) {
	_, err := CreateTemplateCache()
	failIfError(err, t)
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestRenderTemplate(t *testing.T) {
	tc, err := CreateTemplateCache()
	failIfError(err, t)
	app.TemplateCache = tc
	r, err := getSession()
	failIfError(err, t)
	ww := mockWriter{}
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	failIfError(err, t)
	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func failIfError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	res, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		return nil, err
	}
	context := res.Context()
	context, _ = session.Load(context, res.Header.Get("X-Session"))
	res = res.WithContext(context)
	return res, nil
}