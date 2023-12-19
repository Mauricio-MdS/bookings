package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/mauricio-mds/bookings/internal/config"
	"github.com/mauricio-mds/bookings/internal/driver"
	"github.com/mauricio-mds/bookings/internal/forms"
	"github.com/mauricio-mds/bookings/internal/helpers"
	"github.com/mauricio-mds/bookings/internal/models"
	"github.com/mauricio-mds/bookings/internal/render"
	"github.com/mauricio-mds/bookings/internal/repository"
	"github.com/mauricio-mds/bookings/internal/repository/dbrebo"
)

// Repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// dateLayout is the format expected from the forms
const dateLayout = "2006-01-02"

// Repo is the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrebo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new test repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB: dbrebo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is page handler for Home Page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the page handler for the About Page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// ChooseRoom displays available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("could not get reservation from session"))
		return
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Reservation is page handler for Reservation Page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("could not get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format(dateLayout)
	ed := res.EndDate.Format(dateLayout)
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]any)
	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

// PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone= r.Form.Get("phone")
	reservation.Email= r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals is page handler for Generals Page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is page handler for Majors Page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is page handler for Availability Page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is page handler for Post at the Availability Page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	startDate, err := time.Parse(dateLayout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(dateLayout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]any)
	data["rooms"] = rooms

	reservation := models.Reservation{StartDate: startDate, EndDate: endDate}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{Data: data})
}

type jsonResponse struct{
	OK bool `json:"ok"`
	Message string `json:"message"`
	RoomID string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

// AvailabilityJSON is page handler for Availability send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	startDate, _ := time.Parse(dateLayout, sd)
	endDate, _ := time.Parse(dateLayout, ed)

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse{
		OK: available,
		Message: "",
		RoomID: strconv.Itoa(roomID),
		StartDate: sd,
		EndDate: ed,
	}

	out, err := json.MarshalIndent(resp, "","     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	startDate, _ := time.Parse(dateLayout, sd)
	endDate, _ := time.Parse(dateLayout, ed)
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res := models.Reservation{
		RoomID: roomID,
		StartDate: startDate,
		EndDate: endDate,
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Contact is page handler for Contact Page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	m.App.Session.Remove(r.Context(), "reservation")
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data := map[string]any{"reservation": reservation}
	stringMap := map[string]string{
		"start_date": reservation.StartDate.Format(dateLayout),
		"end_date": reservation.EndDate.Format(dateLayout),
	}

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{ Data: data, StringMap: stringMap})
}