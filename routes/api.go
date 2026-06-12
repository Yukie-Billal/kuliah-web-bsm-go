package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"kuliah-web-bsm-go/app/http/controllers"
	"kuliah-web-bsm-go/app/http/middleware"
)

func Api() {
	authController := controllers.NewAuthController()
	customerBookingController := controllers.NewCustomerBookingController()

	adminRoleController := controllers.NewAdminRoleController()
	adminUserController := controllers.NewAdminUserController()
	adminCustomerController := controllers.NewAdminCustomerController()
	adminStudioController := controllers.NewAdminStudioController()
	adminAlatMusikController := controllers.NewAdminAlatMusikController()
	adminBookingController := controllers.NewAdminBookingController()

	// Public Auth Routes
	facades.Route().Post("/api/auth/login", authController.Login)
	facades.Route().Post("/api/auth/register", authController.Register)
	facades.Route().Post("/api/auth/google", authController.GoogleCallback)

	// Protected Customer Routes
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		router.Get("/api/profile", authController.Profile)
		router.Get("/api/studios", customerBookingController.Studios)
		router.Get("/api/studios/{id}/slots", customerBookingController.Slots)
		router.Post("/api/bookings/confirm", customerBookingController.Confirm)
		router.Post("/api/bookings/store", customerBookingController.Store)
		router.Get("/api/bookings/saya", customerBookingController.MyBookings)
		router.Get("/api/bookings/tiket/{bookingId}", customerBookingController.Ticket)
		router.Delete("/api/bookings/{id}", customerBookingController.Cancel)
	})

	// Protected Admin Routes
	facades.Route().Middleware(middleware.Admin()).Group(func(router route.Router) {
		// Roles CRUD
		router.Get("/api/admin/roles", adminRoleController.Index)
		router.Get("/api/admin/roles/{id}", adminRoleController.Show)
		router.Post("/api/admin/roles", adminRoleController.Store)
		router.Put("/api/admin/roles/{id}", adminRoleController.Update)
		router.Delete("/api/admin/roles/{id}", adminRoleController.Destroy)

		// Users CRUD
		router.Get("/api/admin/users", adminUserController.Index)
		router.Get("/api/admin/users/{id}", adminUserController.Show)
		router.Post("/api/admin/users", adminUserController.Store)
		router.Put("/api/admin/users/{id}", adminUserController.Update)
		router.Delete("/api/admin/users/{id}", adminUserController.Destroy)

		// Customers CRUD
		router.Get("/api/admin/customers", adminCustomerController.Index)
		router.Get("/api/admin/customers/{id}", adminCustomerController.Show)
		router.Post("/api/admin/customers", adminCustomerController.Store)
		router.Put("/api/admin/customers/{id}", adminCustomerController.Update)
		router.Delete("/api/admin/customers/{id}", adminCustomerController.Destroy)

		// Studios CRUD
		router.Get("/api/admin/studios", adminStudioController.Index)
		router.Get("/api/admin/studios/{id}", adminStudioController.Show)
		router.Post("/api/admin/studios", adminStudioController.Store)
		router.Put("/api/admin/studios/{id}", adminStudioController.Update)
		router.Delete("/api/admin/studios/{id}", adminStudioController.Destroy)

		// Alat Musik CRUD
		router.Get("/api/admin/alat_musiks", adminAlatMusikController.Index)
		router.Get("/api/admin/alat_musiks/{id}", adminAlatMusikController.Show)
		router.Post("/api/admin/alat_musiks", adminAlatMusikController.Store)
		router.Put("/api/admin/alat_musiks/{id}", adminAlatMusikController.Update)
		router.Delete("/api/admin/alat_musiks/{id}", adminAlatMusikController.Destroy)

		// Bookings CRUD & Custom endpoints
		router.Get("/api/admin/bookings", adminBookingController.Index)
		router.Get("/api/admin/bookings/{id}", adminBookingController.Show)
		router.Post("/api/admin/bookings", adminBookingController.Store)
		router.Put("/api/admin/bookings/{id}", adminBookingController.Update)
		router.Delete("/api/admin/bookings/{id}", adminBookingController.Destroy)
		router.Post("/api/admin/bookings/{id}/approve", adminBookingController.Approve)
		router.Post("/api/admin/bookings/{id}/reject", adminBookingController.Reject)
		router.Post("/api/admin/bookings/scan/verify", adminBookingController.Verify)
	})
}
