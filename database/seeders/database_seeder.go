package seeders

import (
	"encoding/json"
	"fmt"
	"time"

	"kuliah-web-bsm-go/app/facades"
	"kuliah-web-bsm-go/app/models"
)

type DatabaseSeeder struct{}

func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

func (s *DatabaseSeeder) Run() error {
	// 1. Create Roles
	adminRole := models.Role{Role: "admin"}
	customerRole := models.Role{Role: "customer"}

	if err := facades.Orm().Query().Create(&adminRole); err != nil {
		return err
	}
	if err := facades.Orm().Query().Create(&customerRole); err != nil {
		return err
	}

	// 2. Users & Customers
	usersData := []struct {
		Username string
		Email    string
		Telepon  string
	}{
		{"Budi Santoso", "budi@gmail.com", "081234567890"},
		{"Siti Aminah", "siti@gmail.com", "081234567891"},
		{"Ahmad Reza", "reza@gmail.com", "081234567892"},
		{"Dewi Lestari", "dewi@gmail.com", "081234567893"},
		{"Rizky Pratama", "rizky@gmail.com", "081234567894"},
		{"Nina Kartika", "nina@gmail.com", "081234567895"},
		{"Faisal Akbar", "faisal@gmail.com", "081234567896"},
		{"Maya Sari", "maya@gmail.com", "081234567897"},
		{"Andi Wijaya", "andi@gmail.com", "081234567898"},
		{"Rina Melati", "rina@gmail.com", "081234567899"},
		{"Yukie Muhammad Billal", "yukiembillal01@gmail.com", "081234567899"},
	}

	passwordHash, err := facades.Hash().Make("password123")
	if err != nil {
		return err
	}

	now := time.Now()
	var customers []models.Customer

	for _, data := range usersData {
		user := models.User{
			RoleID:          customerRole.ID,
			Username:        data.Username,
			Email:           data.Email,
			Password:        passwordHash,
			EmailVerifiedAt: &now,
		}
		if err := facades.Orm().Query().Create(&user); err != nil {
			return err
		}

		phone := data.Telepon
		customer := models.Customer{
			UserID:  user.ID,
			Nama:    data.Username,
			Telepon: &phone,
		}
		if err := facades.Orm().Query().Create(&customer); err != nil {
			return err
		}
		customers = append(customers, customer)
	}

	// 3. Admin User
	adminPasswordHash, err := facades.Hash().Make("admin123")
	if err != nil {
		return err
	}
	adminUser := models.User{
		RoleID:   adminRole.ID,
		Username: "Admin Master",
		Email:    "admin@studio.com",
		Password: adminPasswordHash,
	}
	if err := facades.Orm().Query().Create(&adminUser); err != nil {
		return err
	}

	// 4. Studios
	studio1Ops, _ := json.Marshal(map[string]string{"buka": "09:00", "tutup": "22:00"})
	studio2Ops, _ := json.Marshal(map[string]string{"buka": "09:00", "tutup": "23:00"})
	studio3Ops, _ := json.Marshal(map[string]string{"buka": "10:00", "tutup": "00:00"})

	s1OpsStr := string(studio1Ops)
	s2OpsStr := string(studio2Ops)
	s3OpsStr := string(studio3Ops)

	studios := []models.Studio{
		{
			NamaStudio:     "Studio Akustik A",
			Lokasi:         "Lantai 1",
			Luas:           20,
			HargaPerJam:    75000,
			IsActive:       true,
			JamOperasional: &s1OpsStr,
		},
		{
			NamaStudio:     "Studio Band B",
			Lokasi:         "Lantai 1",
			Luas:           30,
			HargaPerJam:    100000,
			IsActive:       true,
			JamOperasional: &s2OpsStr,
		},
		{
			NamaStudio:     "Studio Rekaman VIP",
			Lokasi:         "Lantai 2",
			Luas:           25,
			HargaPerJam:    250000,
			IsActive:       true,
			JamOperasional: &s3OpsStr,
		},
	}

	for i := range studios {
		if err := facades.Orm().Query().Create(&studios[i]); err != nil {
			return err
		}
	}

	// 5. Instruments / Alat Musik
	alatData := []struct {
		NamaAlat   string
		Kondisi    string
		Keterangan string
	}{
		{"Gitar Akustik Yamaha F310", "Sangat Baik", "Kondisi senar baru diganti"},
		{"Gitar Elektrik Fender Stratocaster", "Sangat Baik", "Warna sunburst"},
		{"Bass Ibanez SR300E", "Baik", "Active pickup"},
		{"Drum Set Pearl Export", "Baik", "Termasuk cymbal Zildjian"},
		{"Keyboard Yamaha PSR-E463", "Cukup", "Beberapa tuts agak keras"},
		{"Microphone Shure SM58", "Sangat Baik", "Vocal mic standard"},
		{"Microphone Shure SM57", "Baik", "Instrument mic"},
		{"Amplifier Marshall MG100HCFX", "Baik", "Head + Cabinet"},
		{"Amplifier Bass Ampeg BA-210", "Sangat Baik", "Punchy sound"},
		{"Mixer Yamaha MG16XU", "Baik", "Analog mixer with effects"},
	}

	for i, data := range alatData {
		desc := data.Keterangan
		instrument := models.AlatMusik{
			StudioID:   studios[i%3].ID,
			NamaAlat:   data.NamaAlat,
			Kondisi:    data.Kondisi,
			Keterangan: &desc,
		}
		if err := facades.Orm().Query().Create(&instrument); err != nil {
			return err
		}
	}

	// 6. Bookings
	statuses := []string{"pending", "approved", "completed", "cancelled"}
	for i := 0; i < 10; i++ {
		customer := customers[i]
		studio := studios[i%3]

		bookingID := fmt.Sprintf("BK-%s-%04d", time.Now().Format("200601"), i+1)
		desc := "Booking untuk latihan reguler"

		booking := models.Booking{
			BookingID:  bookingID,
			CustomerID: customer.ID,
			StudioID:   studio.ID,
			Tanggal:    time.Now().AddDate(0, 0, i),
			Jam:        fmt.Sprintf("%02d:00", 10+i),
			Durasi:     1 + (i % 3),
			TotalBiaya: studio.HargaPerJam * (1 + (i % 3)),
			Status:     statuses[i%4],
			Catatan:    &desc,
		}
		if err := facades.Orm().Query().Create(&booking); err != nil {
			return err
		}
	}

	return nil
}
