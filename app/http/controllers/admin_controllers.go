package controllers

import (
	"strconv"
	"time"

	"kuliah-web-bsm-go/app/facades"
	"kuliah-web-bsm-go/app/models"

	"github.com/goravel/framework/contracts/http"
)

// ==========================================
// AdminRoleController
// ==========================================

type AdminRoleController struct{}

func NewAdminRoleController() *AdminRoleController {
	return &AdminRoleController{}
}

func (c *AdminRoleController) Index(ctx http.Context) http.Response {
	var roles []models.Role
	if err := facades.Orm().Query().Get(&roles); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(roles)
}

func (c *AdminRoleController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var role models.Role
	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Role not found"})
	}
	return ctx.Response().Success().Json(role)
}

func (c *AdminRoleController) Store(ctx http.Context) http.Response {
	var req struct {
		Role string `json:"role"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Role == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	role := models.Role{Role: req.Role}
	if err := facades.Orm().Query().Create(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(role)
}

func (c *AdminRoleController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var role models.Role
	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Role not found"})
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Role == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	role.Role = req.Role
	if err := facades.Orm().Query().Save(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(role)
}

func (c *AdminRoleController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var role models.Role
	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Role not found"})
	}

	if _, err := facades.Orm().Query().Delete(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "Role deleted successfully"})
}

// ==========================================
// AdminUserController
// ==========================================

type AdminUserController struct{}

func NewAdminUserController() *AdminUserController {
	return &AdminUserController{}
}

func (c *AdminUserController) Index(ctx http.Context) http.Response {
	var users []models.User
	if err := facades.Orm().Query().With("Role").With("Customer").Get(&users); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(users)
}

func (c *AdminUserController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var user models.User
	if err := facades.Orm().Query().With("Role").With("Customer").Where("id", id).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "User not found"})
	}
	return ctx.Response().Success().Json(user)
}

func (c *AdminUserController) Store(ctx http.Context) http.Response {
	var req struct {
		RoleID   uint   `json:"role_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Username == "" || req.Email == "" || req.Password == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	hashedPassword, err := facades.Hash().Make(req.Password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": "Hashing failed"})
	}

	roleID := req.RoleID
	if roleID == 0 {
		roleID = 2 // default to customer
	}

	user := models.User{
		RoleID:   roleID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(user)
}

func (c *AdminUserController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var user models.User
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "User not found"})
	}

	var req struct {
		RoleID   uint   `json:"role_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Username == "" || req.Email == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	user.RoleID = req.RoleID
	user.Username = req.Username
	user.Email = req.Email

	if req.Password != "" {
		hashed, err := facades.Hash().Make(req.Password)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": "Hashing failed"})
		}
		user.Password = hashed
	}

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(user)
}

func (c *AdminUserController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var user models.User
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "User not found"})
	}

	if _, err := facades.Orm().Query().Delete(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "User deleted successfully"})
}

// ==========================================
// AdminCustomerController
// ==========================================

type AdminCustomerController struct{}

func NewAdminCustomerController() *AdminCustomerController {
	return &AdminCustomerController{}
}

func (c *AdminCustomerController) Index(ctx http.Context) http.Response {
	var customers []models.Customer
	if err := facades.Orm().Query().With("User").Get(&customers); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(customers)
}

func (c *AdminCustomerController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var customer models.Customer
	if err := facades.Orm().Query().With("User").Where("id", id).First(&customer); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer not found"})
	}
	return ctx.Response().Success().Json(customer)
}

func (c *AdminCustomerController) Store(ctx http.Context) http.Response {
	var req struct {
		UserID  uint   `json:"user_id"`
		Nama    string `json:"nama"`
		Telepon string `json:"telepon"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.UserID == 0 || req.Nama == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	customer := models.Customer{
		UserID:  req.UserID,
		Nama:    req.Nama,
		Telepon: &req.Telepon,
	}

	if err := facades.Orm().Query().Create(&customer); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(customer)
}

func (c *AdminCustomerController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var customer models.Customer
	if err := facades.Orm().Query().Where("id", id).First(&customer); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer not found"})
	}

	var req struct {
		UserID  uint   `json:"user_id"`
		Nama    string `json:"nama"`
		Telepon string `json:"telepon"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.UserID == 0 || req.Nama == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	customer.UserID = req.UserID
	customer.Nama = req.Nama
	customer.Telepon = &req.Telepon

	if err := facades.Orm().Query().Save(&customer); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(customer)
}

func (c *AdminCustomerController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var customer models.Customer
	if err := facades.Orm().Query().Where("id", id).First(&customer); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer not found"})
	}

	if _, err := facades.Orm().Query().Delete(&customer); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "Customer deleted successfully"})
}

// ==========================================
// AdminStudioController
// ==========================================

type AdminStudioController struct{}

func NewAdminStudioController() *AdminStudioController {
	return &AdminStudioController{}
}

func (c *AdminStudioController) Index(ctx http.Context) http.Response {
	var studios []models.Studio
	if err := facades.Orm().Query().Get(&studios); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(studios)
}

func (c *AdminStudioController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var studio models.Studio
	if err := facades.Orm().Query().Where("id", id).First(&studio); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Studio not found"})
	}
	return ctx.Response().Success().Json(studio)
}

func (c *AdminStudioController) Store(ctx http.Context) http.Response {
	var req struct {
		NamaStudio     string `json:"nama_studio"`
		Lokasi         string `json:"lokasi"`
		Luas           int    `json:"luas"`
		JamOperasional string `json:"jam_operasional"` // expected JSON string like {"buka":"08:00","tutup":"22:00"}
		HargaPerJam    int    `json:"harga_per_jam"`
		IsActive       *bool  `json:"is_active"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.NamaStudio == "" || req.Lokasi == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	studio := models.Studio{
		NamaStudio: req.NamaStudio,
		Lokasi:     req.Lokasi,
		Luas:       req.Luas,
		HargaPerJam: req.HargaPerJam,
		IsActive:   isActive,
	}

	if req.JamOperasional != "" {
		studio.JamOperasional = &req.JamOperasional
	}

	if err := facades.Orm().Query().Create(&studio); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(studio)
}

func (c *AdminStudioController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var studio models.Studio
	if err := facades.Orm().Query().Where("id", id).First(&studio); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Studio not found"})
	}

	var req struct {
		NamaStudio     string `json:"nama_studio"`
		Lokasi         string `json:"lokasi"`
		Luas           int    `json:"luas"`
		JamOperasional string `json:"jam_operasional"`
		HargaPerJam    int    `json:"harga_per_jam"`
		IsActive       *bool  `json:"is_active"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.NamaStudio == "" || req.Lokasi == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	studio.NamaStudio = req.NamaStudio
	studio.Lokasi = req.Lokasi
	studio.Luas = req.Luas
	studio.HargaPerJam = req.HargaPerJam
	if req.IsActive != nil {
		studio.IsActive = *req.IsActive
	}

	if req.JamOperasional != "" {
		studio.JamOperasional = &req.JamOperasional
	} else {
		studio.JamOperasional = nil
	}

	if err := facades.Orm().Query().Save(&studio); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(studio)
}

func (c *AdminStudioController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var studio models.Studio
	if err := facades.Orm().Query().Where("id", id).First(&studio); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Studio not found"})
	}

	if _, err := facades.Orm().Query().Delete(&studio); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "Studio deleted successfully"})
}

// ==========================================
// AdminAlatMusikController
// ==========================================

type AdminAlatMusikController struct{}

func NewAdminAlatMusikController() *AdminAlatMusikController {
	return &AdminAlatMusikController{}
}

func (c *AdminAlatMusikController) Index(ctx http.Context) http.Response {
	var instruments []models.AlatMusik
	if err := facades.Orm().Query().With("Studio").Get(&instruments); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(instruments)
}

func (c *AdminAlatMusikController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var instrument models.AlatMusik
	if err := facades.Orm().Query().With("Studio").Where("id", id).First(&instrument); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Instrument not found"})
	}
	return ctx.Response().Success().Json(instrument)
}

func (c *AdminAlatMusikController) Store(ctx http.Context) http.Response {
	var req struct {
		StudioID   uint   `json:"studio_id"`
		NamaAlat   string `json:"nama_alat"`
		Kondisi    string `json:"kondisi"`
		Keterangan string `json:"keterangan"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.StudioID == 0 || req.NamaAlat == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	kondisi := req.Kondisi
	if kondisi == "" {
		kondisi = "Baik"
	}

	instrument := models.AlatMusik{
		StudioID: req.StudioID,
		NamaAlat: req.NamaAlat,
		Kondisi:  kondisi,
	}

	if req.Keterangan != "" {
		instrument.Keterangan = &req.Keterangan
	}

	if err := facades.Orm().Query().Create(&instrument); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(instrument)
}

func (c *AdminAlatMusikController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var instrument models.AlatMusik
	if err := facades.Orm().Query().Where("id", id).First(&instrument); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Instrument not found"})
	}

	var req struct {
		StudioID   uint   `json:"studio_id"`
		NamaAlat   string `json:"nama_alat"`
		Kondisi    string `json:"kondisi"`
		Keterangan string `json:"keterangan"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.StudioID == 0 || req.NamaAlat == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	instrument.StudioID = req.StudioID
	instrument.NamaAlat = req.NamaAlat
	instrument.Kondisi = req.Kondisi

	if req.Keterangan != "" {
		instrument.Keterangan = &req.Keterangan
	} else {
		instrument.Keterangan = nil
	}

	if err := facades.Orm().Query().Save(&instrument); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(instrument)
}

func (c *AdminAlatMusikController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var instrument models.AlatMusik
	if err := facades.Orm().Query().Where("id", id).First(&instrument); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Instrument not found"})
	}

	if _, err := facades.Orm().Query().Delete(&instrument); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "Instrument deleted successfully"})
}

// ==========================================
// AdminBookingController
// ==========================================

type AdminBookingController struct{}

func NewAdminBookingController() *AdminBookingController {
	return &AdminBookingController{}
}

func (c *AdminBookingController) Index(ctx http.Context) http.Response {
	var bookings []models.Booking
	if err := facades.Orm().Query().With("Customer").With("Studio").Get(&bookings); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(bookings)
}

func (c *AdminBookingController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var booking models.Booking
	if err := facades.Orm().Query().With("Customer").With("Studio").Where("id", id).First(&booking); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Booking not found"})
	}
	return ctx.Response().Success().Json(booking)
}

func (c *AdminBookingController) Store(ctx http.Context) http.Response {
	var req struct {
		BookingID  string `json:"booking_id"`
		CustomerID uint   `json:"customer_id"`
		StudioID   uint   `json:"studio_id"`
		Tanggal    string `json:"tanggal"` // YYYY-MM-DD
		Jam        string `json:"jam"`     // HH:MM
		Durasi     int    `json:"durasi"`
		TotalBiaya int    `json:"total_biaya"`
		Status     string `json:"status"`
		Catatan    string `json:"catatan"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.BookingID == "" || req.CustomerID == 0 || req.StudioID == 0 {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	parsedDate, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid date format"})
	}

	durasi := req.Durasi
	if durasi == 0 {
		durasi = 1
	}

	status := req.Status
	if status == "" {
		status = "pending"
	}

	booking := models.Booking{
		BookingID:  req.BookingID,
		CustomerID: req.CustomerID,
		StudioID:   req.StudioID,
		Tanggal:    parsedDate,
		Jam:        req.Jam,
		Durasi:     durasi,
		TotalBiaya: req.TotalBiaya,
		Status:     status,
	}

	if req.Catatan != "" {
		booking.Catatan = &req.Catatan
	}

	if err := facades.Orm().Query().Create(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(booking)
}

func (c *AdminBookingController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var booking models.Booking
	if err := facades.Orm().Query().Where("id", id).First(&booking); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Booking not found"})
	}

	var req struct {
		Status      string `json:"status"`
		Catatan     string `json:"catatan"`
		CheckedIn   bool   `json:"checked_in"`
		ApprovedBy  uint   `json:"approved_by"`
		TotalBiaya  int    `json:"total_biaya"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid request"})
	}

	if req.Status != "" {
		booking.Status = req.Status
	}

	if req.Catatan != "" {
		booking.Catatan = &req.Catatan
	}

	if req.CheckedIn {
		now := time.Now()
		booking.CheckedInAt = &now
	}

	if req.ApprovedBy != 0 {
		booking.ApprovedBy = &req.ApprovedBy
	}

	if req.TotalBiaya != 0 {
		booking.TotalBiaya = req.TotalBiaya
	}

	if err := facades.Orm().Query().Save(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(booking)
}

func (c *AdminBookingController) Approve(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var booking models.Booking
	if err := facades.Orm().Query().Where("id", id).First(&booking); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Booking not found"})
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthorized"})
	}

	booking.Status = "approved"
	booking.ApprovedBy = &user.ID

	if err := facades.Orm().Query().Save(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(booking)
}

func (c *AdminBookingController) Reject(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var booking models.Booking
	if err := facades.Orm().Query().Where("id", id).First(&booking); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Booking not found"})
	}

	booking.Status = "cancelled" // Reject sets status to cancelled

	if err := facades.Orm().Query().Save(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(booking)
}

func (c *AdminBookingController) Verify(ctx http.Context) http.Response {
	// QR scanner verify (matches scan.verify route in Laravel admin scan)
	code := ctx.Request().Input("code")
	if code == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Code is required"})
	}

	var booking models.Booking
	err := facades.Orm().Query().With("Customer").With("Studio").Where("booking_id", code).First(&booking)
	if err != nil || booking.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Booking code not found"})
	}

	if booking.Status != "approved" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Booking status is: " + booking.Status + ". Must be approved to check-in.",
		})
	}

	if booking.CheckedInAt != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Booking already checked-in at " + booking.CheckedInAt.Format("15:04"),
		})
	}

	now := time.Now()
	booking.CheckedInAt = &now
	booking.Status = "completed"

	if err := facades.Orm().Query().Save(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "Check-in successful!",
		"booking": booking,
	})
}

func (c *AdminBookingController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var booking models.Booking
	if err := facades.Orm().Query().Where("id", id).First(&booking); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Booking not found"})
	}

	if _, err := facades.Orm().Query().Delete(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "Booking deleted successfully"})
}

// Convert string to uint helper
func getUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 32)
	return uint(val)
}
