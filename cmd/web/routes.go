package main

import (
	"net/http"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(NoSurf)
	r.Use(SessionLoad)

	r.Get("/", handlers.Repo.Home)
	r.Get("/about", handlers.Repo.About)
	r.Get("/generals-quarters", handlers.Repo.Generals)
	r.Get("/majors-suite", handlers.Repo.Majors)
	r.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	r.Get("/book-room", handlers.Repo.BookRoom)

	r.Get("/search-availability", handlers.Repo.Availability)
	r.Post("/search-availability", handlers.Repo.PostAvailability)
	r.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	r.Get("/contact", handlers.Repo.Contact)
	r.Get("/make-reservation", handlers.Repo.Reservation)
	r.Post("/make-reservation", handlers.Repo.PostReservation)
	r.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	r.Get("/user/login", handlers.Repo.ShowLogin)
	r.Post("/user/login", handlers.Repo.PostShowLogin)
	r.Get("/user/logout", handlers.Repo.Logout)
	//redirect to secure page for admin user
	r.Route("/admin", func(r chi.Router) {
		//r.Use(Auth)
		r.Get("/dashboard", handlers.Repo.AdminDashboard)

		r.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		r.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		r.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		r.Get("/process-reservation/{src}/{id}", handlers.Repo.AdminProcessReservation)
		r.Get("/delete-reservation/{src}/{id}", handlers.Repo.AdminDeleteReservation)
		r.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
		r.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
		r.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)
	})
	fileServer := http.FileServer(http.Dir("./static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return r
}
