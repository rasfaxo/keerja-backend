package seeders

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Province represents a province entity
type Province struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
	Code string `gorm:"uniqueIndex;not null"`
}

func (Province) TableName() string {
	return "provinces"
}

// City represents a city entity
type City struct {
	ID         int64  `gorm:"primaryKey"`
	ProvinceID int64  `gorm:"not null"`
	Name       string `gorm:"not null"`
	Type       string `gorm:"not null"` // Kota or Kabupaten
	Code       string `gorm:"uniqueIndex;not null"`
}

func (City) TableName() string {
	return "cities"
}

// District represents a district entity
type District struct {
	ID         int64   `gorm:"primaryKey"`
	CityID     int64   `gorm:"not null"`
	Name       string  `gorm:"not null"`
	Code       string  `gorm:"uniqueIndex;not null"`
	PostalCode *string `gorm:"size:10"`
}

func (District) TableName() string {
	return "districts"
}

// LocationSeeder seeds provinces, cities, and districts
func LocationSeeder(db *gorm.DB) error {
	if err := SeedProvinces(db); err != nil {
		return err
	}
	if err := SeedCities(db); err != nil {
		return err
	}
	if err := SeedDistricts(db); err != nil {
		return err
	}
	return nil
}

// SeedProvinces seeds all 34 provinces of Indonesia
func SeedProvinces(db *gorm.DB) error {
	log.Println("Seeding provinces table...")

	provinces := []Province{
		{Name: "Aceh", Code: "11"},
		{Name: "Sumatera Utara", Code: "12"},
		{Name: "Sumatera Barat", Code: "13"},
		{Name: "Riau", Code: "14"},
		{Name: "Jambi", Code: "15"},
		{Name: "Sumatera Selatan", Code: "16"},
		{Name: "Bengkulu", Code: "17"},
		{Name: "Lampung", Code: "18"},
		{Name: "Kepulauan Bangka Belitung", Code: "19"},
		{Name: "Kepulauan Riau", Code: "21"},
		{Name: "DKI Jakarta", Code: "31"},
		{Name: "Jawa Barat", Code: "32"},
		{Name: "Jawa Tengah", Code: "33"},
		{Name: "DI Yogyakarta", Code: "34"},
		{Name: "Jawa Timur", Code: "35"},
		{Name: "Banten", Code: "36"},
		{Name: "Bali", Code: "51"},
		{Name: "Nusa Tenggara Barat", Code: "52"},
		{Name: "Nusa Tenggara Timur", Code: "53"},
		{Name: "Kalimantan Barat", Code: "61"},
		{Name: "Kalimantan Tengah", Code: "62"},
		{Name: "Kalimantan Selatan", Code: "63"},
		{Name: "Kalimantan Timur", Code: "64"},
		{Name: "Kalimantan Utara", Code: "65"},
		{Name: "Sulawesi Utara", Code: "71"},
		{Name: "Sulawesi Tengah", Code: "72"},
		{Name: "Sulawesi Selatan", Code: "73"},
		{Name: "Sulawesi Tenggara", Code: "74"},
		{Name: "Gorontalo", Code: "75"},
		{Name: "Sulawesi Barat", Code: "76"},
		{Name: "Maluku", Code: "81"},
		{Name: "Maluku Utara", Code: "82"},
		{Name: "Papua Barat", Code: "91"},
		{Name: "Papua", Code: "94"},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&provinces)

	if result.Error != nil {
		log.Printf("Failed to seed provinces: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d provinces", len(provinces))
	return nil
}

// SeedCities seeds major cities in Indonesia
func SeedCities(db *gorm.DB) error {
	log.Println("Seeding cities table...")

	// Helper to get province ID by code
	getProvinceID := func(code string) int64 {
		var province Province
		if err := db.Where("code = ?", code).First(&province).Error; err != nil {
			return 0
		}
		return province.ID
	}

	// DKI Jakarta (31)
	jakartaID := getProvinceID("31")
	// Jawa Barat (32)
	jabarID := getProvinceID("32")
	// Banten (36)
	bantenID := getProvinceID("36")
	// Jawa Tengah (33)
	jatengID := getProvinceID("33")
	// DI Yogyakarta (34)
	yogyaID := getProvinceID("34")
	// Jawa Timur (35)
	jatimID := getProvinceID("35")
	// Sumatera Utara (12)
	sumutID := getProvinceID("12")
	// Sumatera Selatan (16)
	sumsalID := getProvinceID("16")
	// Kepulauan Riau (21)
	kepriID := getProvinceID("21")
	// Bali (51)
	baliID := getProvinceID("51")
	// Sulawesi Selatan (73)
	sulselID := getProvinceID("73")
	// Kalimantan Timur (64)
	kaltimID := getProvinceID("64")

	cities := []City{
		// DKI Jakarta
		{ProvinceID: jakartaID, Name: "Jakarta Pusat", Type: "Kota", Code: "3171"},
		{ProvinceID: jakartaID, Name: "Jakarta Utara", Type: "Kota", Code: "3172"},
		{ProvinceID: jakartaID, Name: "Jakarta Barat", Type: "Kota", Code: "3173"},
		{ProvinceID: jakartaID, Name: "Jakarta Selatan", Type: "Kota", Code: "3174"},
		{ProvinceID: jakartaID, Name: "Jakarta Timur", Type: "Kota", Code: "3175"},
		{ProvinceID: jakartaID, Name: "Kepulauan Seribu", Type: "Kabupaten", Code: "3101"},

		// Jawa Barat
		{ProvinceID: jabarID, Name: "Bogor", Type: "Kabupaten", Code: "3201"},
		{ProvinceID: jabarID, Name: "Sukabumi", Type: "Kabupaten", Code: "3202"},
		{ProvinceID: jabarID, Name: "Cianjur", Type: "Kabupaten", Code: "3203"},
		{ProvinceID: jabarID, Name: "Bandung", Type: "Kabupaten", Code: "3204"},
		{ProvinceID: jabarID, Name: "Garut", Type: "Kabupaten", Code: "3205"},
		{ProvinceID: jabarID, Name: "Tasikmalaya", Type: "Kabupaten", Code: "3206"},
		{ProvinceID: jabarID, Name: "Ciamis", Type: "Kabupaten", Code: "3207"},
		{ProvinceID: jabarID, Name: "Kuningan", Type: "Kabupaten", Code: "3208"},
		{ProvinceID: jabarID, Name: "Cirebon", Type: "Kabupaten", Code: "3209"},
		{ProvinceID: jabarID, Name: "Majalengka", Type: "Kabupaten", Code: "3210"},
		{ProvinceID: jabarID, Name: "Sumedang", Type: "Kabupaten", Code: "3211"},
		{ProvinceID: jabarID, Name: "Indramayu", Type: "Kabupaten", Code: "3212"},
		{ProvinceID: jabarID, Name: "Subang", Type: "Kabupaten", Code: "3213"},
		{ProvinceID: jabarID, Name: "Purwakarta", Type: "Kabupaten", Code: "3214"},
		{ProvinceID: jabarID, Name: "Karawang", Type: "Kabupaten", Code: "3215"},
		{ProvinceID: jabarID, Name: "Bekasi", Type: "Kabupaten", Code: "3216"},
		{ProvinceID: jabarID, Name: "Bandung Barat", Type: "Kabupaten", Code: "3217"},
		{ProvinceID: jabarID, Name: "Pangandaran", Type: "Kabupaten", Code: "3218"},
		{ProvinceID: jabarID, Name: "Kota Bogor", Type: "Kota", Code: "3271"},
		{ProvinceID: jabarID, Name: "Kota Sukabumi", Type: "Kota", Code: "3272"},
		{ProvinceID: jabarID, Name: "Kota Bandung", Type: "Kota", Code: "3273"},
		{ProvinceID: jabarID, Name: "Kota Cirebon", Type: "Kota", Code: "3274"},
		{ProvinceID: jabarID, Name: "Kota Bekasi", Type: "Kota", Code: "3275"},
		{ProvinceID: jabarID, Name: "Kota Depok", Type: "Kota", Code: "3276"},
		{ProvinceID: jabarID, Name: "Kota Cimahi", Type: "Kota", Code: "3277"},
		{ProvinceID: jabarID, Name: "Kota Tasikmalaya", Type: "Kota", Code: "3278"},
		{ProvinceID: jabarID, Name: "Kota Banjar", Type: "Kota", Code: "3279"},

		// Banten
		{ProvinceID: bantenID, Name: "Pandeglang", Type: "Kabupaten", Code: "3601"},
		{ProvinceID: bantenID, Name: "Lebak", Type: "Kabupaten", Code: "3602"},
		{ProvinceID: bantenID, Name: "Tangerang", Type: "Kabupaten", Code: "3603"},
		{ProvinceID: bantenID, Name: "Serang", Type: "Kabupaten", Code: "3604"},
		{ProvinceID: bantenID, Name: "Kota Tangerang", Type: "Kota", Code: "3671"},
		{ProvinceID: bantenID, Name: "Kota Cilegon", Type: "Kota", Code: "3672"},
		{ProvinceID: bantenID, Name: "Kota Serang", Type: "Kota", Code: "3673"},
		{ProvinceID: bantenID, Name: "Kota Tangerang Selatan", Type: "Kota", Code: "3674"},

		// Jawa Tengah
		{ProvinceID: jatengID, Name: "Kota Semarang", Type: "Kota", Code: "3374"},
		{ProvinceID: jatengID, Name: "Kota Solo", Type: "Kota", Code: "3372"},

		// DI Yogyakarta
		{ProvinceID: yogyaID, Name: "Kota Yogyakarta", Type: "Kota", Code: "3471"},
		{ProvinceID: yogyaID, Name: "Sleman", Type: "Kabupaten", Code: "3404"},
		{ProvinceID: yogyaID, Name: "Bantul", Type: "Kabupaten", Code: "3402"},

		// Jawa Timur
		{ProvinceID: jatimID, Name: "Kota Surabaya", Type: "Kota", Code: "3578"},
		{ProvinceID: jatimID, Name: "Kota Malang", Type: "Kota", Code: "3573"},
		{ProvinceID: jatimID, Name: "Sidoarjo", Type: "Kabupaten", Code: "3515"},

		// Sumatera Utara
		{ProvinceID: sumutID, Name: "Kota Medan", Type: "Kota", Code: "1271"},

		// Sumatera Selatan
		{ProvinceID: sumsalID, Name: "Kota Palembang", Type: "Kota", Code: "1671"},

		// Kepulauan Riau
		{ProvinceID: kepriID, Name: "Kota Batam", Type: "Kota", Code: "2171"},

		// Bali
		{ProvinceID: baliID, Name: "Kota Denpasar", Type: "Kota", Code: "5171"},
		{ProvinceID: baliID, Name: "Badung", Type: "Kabupaten", Code: "5103"},

		// Sulawesi Selatan
		{ProvinceID: sulselID, Name: "Kota Makassar", Type: "Kota", Code: "7371"},

		// Kalimantan Timur
		{ProvinceID: kaltimID, Name: "Kota Balikpapan", Type: "Kota", Code: "6471"},
		{ProvinceID: kaltimID, Name: "Kota Samarinda", Type: "Kota", Code: "6472"},
	}

	// Insert cities one by one, skip if already exists
	insertedCount := 0
	for _, city := range cities {
		if city.ProvinceID == 0 {
			continue
		}

		// Check if city with this code already exists
		var existing City
		if err := db.Where("code = ?", city.Code).First(&existing).Error; err == nil {
			// Already exists, skip
			continue
		}

		if err := db.Create(&city).Error; err != nil {
			log.Printf("Failed to seed city %s: %v", city.Name, err)
			continue
		}
		insertedCount++
	}

	log.Printf("Successfully seeded %d cities", insertedCount)
	return nil
}

// SeedDistricts seeds sample districts in major cities
func SeedDistricts(db *gorm.DB) error {
	log.Println("Seeding districts table...")

	// Helper to get city ID by code
	getCityID := func(code string) int64 {
		var city City
		if err := db.Where("code = ?", code).First(&city).Error; err != nil {
			return 0
		}
		return city.ID
	}

	// Helper for postal code pointer
	postal := func(s string) *string { return &s }

	// Jakarta Pusat (3171)
	jakpusID := getCityID("3171")
	// Jakarta Selatan (3174)
	jakselID := getCityID("3174")
	// Kota Bandung (3273)
	bandungID := getCityID("3273")
	// Bandung Barat (3217)
	bandungBaratID := getCityID("3217")
	// Tangerang Selatan (3674)
	tangselID := getCityID("3674")

	districts := []District{
		// Jakarta Pusat
		{CityID: jakpusID, Name: "Gambir", Code: "3171010", PostalCode: postal("10110")},
		{CityID: jakpusID, Name: "Tanah Abang", Code: "3171020", PostalCode: postal("10120")},
		{CityID: jakpusID, Name: "Menteng", Code: "3171030", PostalCode: postal("10130")},
		{CityID: jakpusID, Name: "Senen", Code: "3171040", PostalCode: postal("10140")},
		{CityID: jakpusID, Name: "Cempaka Putih", Code: "3171050", PostalCode: postal("10150")},
		{CityID: jakpusID, Name: "Johar Baru", Code: "3171060", PostalCode: postal("10160")},
		{CityID: jakpusID, Name: "Kemayoran", Code: "3171070", PostalCode: postal("10170")},
		{CityID: jakpusID, Name: "Sawah Besar", Code: "3171080", PostalCode: postal("10180")},

		// Jakarta Selatan
		{CityID: jakselID, Name: "Tebet", Code: "3174010", PostalCode: postal("12810")},
		{CityID: jakselID, Name: "Setiabudi", Code: "3174020", PostalCode: postal("12910")},
		{CityID: jakselID, Name: "Mampang Prapatan", Code: "3174030", PostalCode: postal("12790")},
		{CityID: jakselID, Name: "Pasar Minggu", Code: "3174040", PostalCode: postal("12510")},
		{CityID: jakselID, Name: "Kebayoran Lama", Code: "3174050", PostalCode: postal("12210")},
		{CityID: jakselID, Name: "Cilandak", Code: "3174060", PostalCode: postal("12430")},
		{CityID: jakselID, Name: "Kebayoran Baru", Code: "3174070", PostalCode: postal("12110")},
		{CityID: jakselID, Name: "Pancoran", Code: "3174080", PostalCode: postal("12780")},
		{CityID: jakselID, Name: "Jagakarsa", Code: "3174090", PostalCode: postal("12620")},
		{CityID: jakselID, Name: "Pesanggrahan", Code: "3174100", PostalCode: postal("12270")},

		// Kota Bandung
		{CityID: bandungID, Name: "Bandung Wetan", Code: "3273010", PostalCode: postal("40114")},
		{CityID: bandungID, Name: "Sumur Bandung", Code: "3273020", PostalCode: postal("40111")},
		{CityID: bandungID, Name: "Cibeunying Kaler", Code: "3273030", PostalCode: postal("40171")},
		{CityID: bandungID, Name: "Cibeunying Kidul", Code: "3273040", PostalCode: postal("40121")},
		{CityID: bandungID, Name: "Coblong", Code: "3273050", PostalCode: postal("40132")},
		{CityID: bandungID, Name: "Sukasari", Code: "3273060", PostalCode: postal("40152")},
		{CityID: bandungID, Name: "Cidadap", Code: "3273070", PostalCode: postal("40141")},

		// Bandung Barat
		{CityID: bandungBaratID, Name: "Lembang", Code: "3217010", PostalCode: postal("40391")},
		{CityID: bandungBaratID, Name: "Parongpong", Code: "3217020", PostalCode: postal("40559")},
		{CityID: bandungBaratID, Name: "Cisarua", Code: "3217030", PostalCode: postal("40551")},
		{CityID: bandungBaratID, Name: "Cikalong Wetan", Code: "3217040", PostalCode: postal("40560")},
		{CityID: bandungBaratID, Name: "Cipeundeuy", Code: "3217050", PostalCode: postal("40558")},
		{CityID: bandungBaratID, Name: "Ngamprah", Code: "3217060", PostalCode: postal("40552")},
		{CityID: bandungBaratID, Name: "Cipatat", Code: "3217070", PostalCode: postal("40553")},
		{CityID: bandungBaratID, Name: "Padalarang", Code: "3217080", PostalCode: postal("40553")},
		{CityID: bandungBaratID, Name: "Batujajar", Code: "3217090", PostalCode: postal("40561")},
		{CityID: bandungBaratID, Name: "Cihampelas", Code: "3217100", PostalCode: postal("40562")},

		// Tangerang Selatan
		{CityID: tangselID, Name: "Serpong", Code: "3674010", PostalCode: postal("15310")},
		{CityID: tangselID, Name: "Serpong Utara", Code: "3674020", PostalCode: postal("15310")},
		{CityID: tangselID, Name: "Pondok Aren", Code: "3674030", PostalCode: postal("15224")},
		{CityID: tangselID, Name: "Ciputat", Code: "3674040", PostalCode: postal("15411")},
		{CityID: tangselID, Name: "Ciputat Timur", Code: "3674050", PostalCode: postal("15412")},
		{CityID: tangselID, Name: "Pamulang", Code: "3674060", PostalCode: postal("15417")},
		{CityID: tangselID, Name: "Setu", Code: "3674070", PostalCode: postal("15314")},
	}

	// Insert districts one by one, skip if already exists
	insertedCount := 0
	for _, district := range districts {
		if district.CityID == 0 {
			continue
		}

		// Check if district with this code already exists
		var existing District
		if err := db.Where("code = ?", district.Code).First(&existing).Error; err == nil {
			// Already exists, skip
			continue
		}

		if err := db.Create(&district).Error; err != nil {
			log.Printf("Failed to seed district %s: %v", district.Name, err)
			continue
		}
		insertedCount++
	}

	log.Printf("Successfully seeded %d districts", insertedCount)
	return nil
}
