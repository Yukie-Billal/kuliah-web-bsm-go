package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"kuliah-web-bsm-go/app/facades"
	"kuliah-web-bsm-go/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/mail"
)

type CustomerBookingController struct{}

type Slot struct {
	Hour      int    `json:"hour"`
	Label     string `json:"label"`
	IsPast    bool   `json:"is_past"`
	IsBooked  bool   `json:"is_booked"`
	Available bool   `json:"available"`
}

func NewCustomerBookingController() *CustomerBookingController {
	return &CustomerBookingController{}
}

// Studios GET /api/studios
func (c *CustomerBookingController) Studios(ctx http.Context) http.Response {
	var studios []models.Studio
	err := facades.Orm().Query().Where("is_active", true).Get(&studios)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to retrieve studios",
		})
	}
	return ctx.Response().Success().Json(studios)
}

// Slots GET /api/studios/{id}/slots
func (c *CustomerBookingController) Slots(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var studio models.Studio
	err := facades.Orm().Query().Where("id", id).First(&studio)
	if err != nil || studio.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Studio not found",
		})
	}

	dateStr := ctx.Request().Query("date", time.Now().Format("2006-01-02"))
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"error": "Tanggal tidak valid",
		})
	}

	// Prevent past dates (compare start of day)
	todayStr := time.Now().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", todayStr)
	if parsedDate.Before(today) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"error": "Tanggal tidak boleh di masa lalu",
		})
	}

	slots := c.buildSlots(studio, dateStr)
	return ctx.Response().Success().Json(http.Json{
		"studio":        studio,
		"selected_date": dateStr,
		"slots":         slots,
	})
}

// Confirm POST /api/bookings/confirm
func (c *CustomerBookingController) Confirm(ctx http.Context) http.Response {
	var req struct {
		StudioID  uint   `json:"studio_id"`
		Date      string `json:"date"`
		StartHour int    `json:"start_hour"`
		Duration  int    `json:"duration"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid request parameters",
		})
	}

	if req.StudioID == 0 || req.Date == "" || req.StartHour < 8 || req.StartHour > 21 || req.Duration < 1 || req.Duration > 4 {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid booking parameters. Jam mulai harus di antara 8-21 dan durasi 1-4 jam.",
		})
	}

	var studio models.Studio
	err := facades.Orm().Query().Where("id", req.StudioID).First(&studio)
	if err != nil || studio.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Studio not found",
		})
	}

	// Validate not in past
	jamMulai := fmt.Sprintf("%02d:00", req.StartHour)
	startDTStr := fmt.Sprintf("%s %s", req.Date, jamMulai)
	startDT, err := time.Parse("2006-01-02 15:04", startDTStr)
	if err != nil || startDT.Before(time.Now()) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "Jam yang dipilih sudah lewat atau tanggal tidak valid.",
		})
	}

	// Check conflict
	if c.hasConflict(studio.ID, req.Date, req.StartHour, req.Duration) {
		return ctx.Response().Json(http.StatusConflict, http.Json{
			"message": "Slot tersebut sudah dipesan orang lain. Silakan pilih jam lain.",
		})
	}

	totalPrice := studio.HargaPerJam * req.Duration
	endDT := startDT.Add(time.Duration(req.Duration) * time.Hour)

	return ctx.Response().Success().Json(http.Json{
		"studio":      studio,
		"date":        req.Date,
		"start_time":  startDT.Format("15:04"),
		"end_time":    endDT.Format("15:04"),
		"duration":    req.Duration,
		"total_biaya": totalPrice,
	})
}

// Store POST /api/bookings/store
func (c *CustomerBookingController) Store(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthorized"})
	}

	var customer models.Customer
	err = facades.Orm().Query().Where("user_id", user.ID).First(&customer)
	if err != nil || customer.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer profile not found"})
	}

	var req struct {
		StudioID  uint   `json:"studio_id"`
		Date      string `json:"date"`
		StartHour int    `json:"start_hour"`
		Duration  int    `json:"duration"`
		Catatan   string `json:"catatan"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid request parameters",
		})
	}

	if req.StudioID == 0 || req.Date == "" || req.StartHour < 8 || req.StartHour > 21 || req.Duration < 1 || req.Duration > 4 {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid booking parameters.",
		})
	}

	var studio models.Studio
	err = facades.Orm().Query().Where("id", req.StudioID).First(&studio)
	if err != nil || studio.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Studio not found",
		})
	}

	jamMulai := fmt.Sprintf("%02d:00", req.StartHour)
	startDTStr := fmt.Sprintf("%s %s", req.Date, jamMulai)
	startDT, err := time.Parse("2006-01-02 15:04", startDTStr)
	if err != nil || startDT.Before(time.Now()) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "Jam yang dipilih sudah lewat.",
		})
	}

	if c.hasConflict(studio.ID, req.Date, req.StartHour, req.Duration) {
		return ctx.Response().Json(http.StatusConflict, http.Json{
			"message": "Slot sudah dipesan. Silakan pilih jam lain.",
		})
	}

	// Generate booking_id SS-YYYY-XXXX
	var lastBooking models.Booking
	currentYear := time.Now().Year()
	err = facades.Orm().Query().Where("booking_id LIKE ?", fmt.Sprintf("SS-%d-%%", currentYear)).
		Order("id DESC").First(&lastBooking)

	seq := 1
	if err == nil && lastBooking.ID > 0 {
		// Parse sequence from last booking_id: SS-2024-0001
		var y int
		var s int
		n, _ := fmt.Sscanf(lastBooking.BookingID, "SS-%d-%d", &y, &s)
		if n == 2 {
			seq = s + 1
		}
	}
	bookingID := fmt.Sprintf("SS-%d-%04d", currentYear, seq)

	parsedDate, _ := time.Parse("2006-01-02", req.Date)

	booking := models.Booking{
		BookingID:  bookingID,
		CustomerID: customer.ID,
		StudioID:   studio.ID,
		Tanggal:    parsedDate,
		Jam:        jamMulai,
		Durasi:     req.Duration,
		TotalBiaya: studio.HargaPerJam * req.Duration,
		Status:     "pending",
	}

	if req.Catatan != "" {
		booking.Catatan = &req.Catatan
	}

	if err := facades.Orm().Query().Create(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to create booking record",
		})
	}

	// Send confirmation email asynchronously
	go func(toEmail string, customerName string, bookingID string, studioName string, tanggal time.Time, jamMulai string, durasi int, totalBiaya int) {
		// Format tanggal
		days := map[string]string{
			"Sunday": "Minggu", "Monday": "Senin", "Tuesday": "Selasa", "Wednesday": "Rabu",
			"Thursday": "Kamis", "Friday": "Jumat", "Saturday": "Sabtu",
		}
		months := map[string]string{
			"January": "Januari", "February": "Februari", "March": "Maret", "April": "April",
			"May": "Mei", "June": "Juni", "July": "Juli", "August": "Agustus",
			"September": "September", "October": "Oktober", "November": "November", "December": "Desember",
		}

		dayEng := tanggal.Format("Monday")
		dayInd := days[dayEng]
		if dayInd == "" {
			dayInd = dayEng
		}
		monthEng := tanggal.Format("January")
		monthInd := months[monthEng]
		if monthInd == "" {
			monthInd = monthEng
		}

		formattedTanggal := fmt.Sprintf("%s, %02d %s %d", dayInd, tanggal.Day(), monthInd, tanggal.Year())

		var startHour int
		fmt.Sscanf(jamMulai[:2], "%d", &startHour)
		jamSelesai := fmt.Sprintf("%02d:00", startHour+durasi)

		// Format currency: 75000 -> 75.000
		totalBiayaStr := ""
		val := totalBiaya
		if val == 0 {
			totalBiayaStr = "0"
		} else {
			var parts []string
			for val > 0 {
				rem := val % 1000
				val = val / 1000
				if val > 0 {
					parts = append([]string{fmt.Sprintf("%03d", rem)}, parts...)
				} else {
					parts = append([]string{fmt.Sprintf("%d", rem)}, parts...)
				}
			}
			for i, p := range parts {
				if i == 0 {
					totalBiayaStr = p
				} else {
					totalBiayaStr = totalBiayaStr + "." + p
				}
			}
		}

		// Render HTML using Go's html/template
		var htmlContent string
		tmpl, err := template.ParseFiles("resources/views/mail/booking_confirmation.html")
		if err == nil {
			var buf bytes.Buffer
			err = tmpl.Execute(&buf, map[string]any{
				"CustomerName": customerName,
				"BookingID":    bookingID,
				"StudioName":   studioName,
				"Tanggal":      formattedTanggal,
				"JamMulai":     jamMulai[:5],
				"JamSelesai":    jamSelesai,
				"Durasi":       durasi,
				"TotalBiaya":   totalBiayaStr,
				"CurrentYear":  time.Now().Year(),
				"AppURL":       facades.Config().GetString("http.url", "http://localhost"),
			})
			if err == nil {
				htmlContent = buf.String()
			}
		}

		if htmlContent == "" {
			htmlContent = fmt.Sprintf("<h1>Booking Confirmed</h1><p>Booking ID: %s</p>", bookingID)
		}

		facades.Mail().To([]string{toEmail}).
			Content(mail.Content{
				Html: htmlContent,
			}).
			Send()
	}(user.Email, customer.Nama, bookingID, studio.NamaStudio, parsedDate, jamMulai, req.Duration, studio.HargaPerJam*req.Duration)

	return ctx.Response().Success().Json(http.Json{
		"message": "Booking berhasil!",
		"booking": booking,
	})
}

// MyBookings GET /api/bookings/saya
func (c *CustomerBookingController) MyBookings(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthorized"})
	}

	var customer models.Customer
	err = facades.Orm().Query().Where("user_id", user.ID).First(&customer)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer profile not found"})
	}

	var bookings []models.Booking
	err = facades.Orm().Query().With("Studio").Where("customer_id", customer.ID).Order("id DESC").Get(&bookings)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to retrieve bookings",
		})
	}

	return ctx.Response().Success().Json(bookings)
}

// Ticket GET /api/bookings/tiket/{bookingId}
func (c *CustomerBookingController) Ticket(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthorized"})
	}

	var customer models.Customer
	err = facades.Orm().Query().Where("user_id", user.ID).First(&customer)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer profile not found"})
	}

	bookingID := ctx.Request().Route("bookingId")
	var booking models.Booking
	err = facades.Orm().Query().With("Studio").With("Customer").Where("booking_id", bookingID).First(&booking)
	if err != nil || booking.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Booking not found",
		})
	}

	// Verify ownership
	if booking.CustomerID != customer.ID {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "Forbidden",
		})
	}

	return ctx.Response().Success().Json(booking)
}

// Cancel DELETE /api/bookings/{id}
func (c *CustomerBookingController) Cancel(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthorized"})
	}

	var customer models.Customer
	err = facades.Orm().Query().Where("user_id", user.ID).First(&customer)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Customer profile not found"})
	}

	id := ctx.Request().Route("id")
	var booking models.Booking
	err = facades.Orm().Query().Where("id", id).First(&booking)
	if err != nil || booking.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Booking not found",
		})
	}

	// Verify ownership
	if booking.CustomerID != customer.ID {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "Forbidden",
		})
	}

	if booking.Status == "cancelled" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Booking is already cancelled",
		})
	}

	booking.Status = "cancelled"
	if err := facades.Orm().Query().Save(&booking); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to cancel booking",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "Booking cancelled successfully",
		"booking": booking,
	})
}

// Helpers
func (c *CustomerBookingController) hasConflict(studioID uint, dateStr string, startHour int, durasi int) bool {
	var jamSlots []any
	for h := startHour; h < startHour+durasi; h++ {
		jamSlots = append(jamSlots, fmt.Sprintf("%02d:00", h))
	}

	count, err := facades.Orm().Query().Model(&models.Booking{}).
		Where("studio_id = ?", studioID).
		Where("tanggal = ?", dateStr).
		WhereNotIn("status", []any{"cancelled"}).
		WhereIn("jam", jamSlots).
		Count()

	return err == nil && count > 0
}

func (c *CustomerBookingController) buildSlots(studio models.Studio, dateStr string) []Slot {
	openHour := 8
	closeHour := 22

	var ops struct {
		Buka  string `json:"buka"`
		Tutup string `json:"tutup"`
	}

	if studio.JamOperasional != nil && *studio.JamOperasional != "" {
		if err := json.Unmarshal([]byte(*studio.JamOperasional), &ops); err == nil {
			if len(ops.Buka) >= 2 {
				fmt.Sscanf(ops.Buka[:2], "%d", &openHour)
			}
			if len(ops.Tutup) >= 2 {
				fmt.Sscanf(ops.Tutup[:2], "%d", &closeHour)
			}
		}
	}

	var bookings []models.Booking
	facades.Orm().Query().Where("studio_id = ?", studio.ID).
		Where("tanggal = ?", dateStr).
		WhereNotIn("status", []any{"cancelled"}).
		Get(&bookings)

	occupiedHours := make(map[int]bool)
	for _, b := range bookings {
		var startH int
		if len(b.Jam) >= 2 {
			fmt.Sscanf(b.Jam[:2], "%d", &startH)
		}
		for h := startH; h < startH+b.Durasi; h++ {
			occupiedHours[h] = true
		}
	}

	now := time.Now()
	isToday := dateStr == now.Format("2006-01-02")
	currentHour := now.Hour()

	var slots []Slot
	for hour := openHour; hour < closeHour; hour++ {
		isPast := isToday && hour <= currentHour
		isBooked := occupiedHours[hour]

		slots = append(slots, Slot{
			Hour:      hour,
			Label:     fmt.Sprintf("%02d:00", hour),
			IsPast:    isPast,
			IsBooked:  isBooked,
			Available: !isPast && !isBooked,
		})
	}

	return slots
}
