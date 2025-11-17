package seeders

// IndustryList provides standardized industry options for Indonesia
var IndustryList = []string{
	// Technology & Digital
	"Technology",
	"Information Technology",
	"Software Development",
	"Financial Technology",
	"Education Technology",
	"Healthcare Technology",
	"E-Commerce",
	"Digital Marketing",

	// Financial Services
	"Banking & Finance",
	"Insurance",
	"Investment",
	"Accounting",

	// Professional Services
	"Consulting",
	"Legal Services",
	"Human Resources",
	"Marketing & Advertising",
	"Public Relations",

	// Manufacturing & Industrial
	"Manufacturing",
	"Automotive",
	"Electronics",
	"Chemical",
	"Food & Beverage",
	"Textile & Apparel",
	"Pharmaceutical",

	// Infrastructure & Construction
	"Construction",
	"Real Estate",
	"Architecture",
	"Engineering",
	"Infrastructure Development",

	// Transportation & Logistics
	"Transportation",
	"Logistics & Supply Chain",
	"Aviation",
	"Maritime",
	"Warehousing",

	// Retail & Consumer
	"Retail",
	"Wholesale",
	"Consumer Goods",
	"Fashion & Lifestyle",

	// Hospitality & Tourism
	"Hotels & Hospitality",
	"Travel & Tourism",
	"Restaurant & Food Service",
	"Event Management",

	// Media & Entertainment
	"Media & Publishing",
	"Entertainment",
	"Broadcasting",
	"Film & Video Production",
	"Gaming",

	// Healthcare
	"Healthcare Services",
	"Medical Devices",
	"Hospital & Clinic",
	"Laboratory Services",

	// Education
	"Education",
	"Training & Development",
	"Research",

	// Energy & Utilities
	"Oil & Gas",
	"Mining",
	"Energy",
	"Utilities",
	"Renewable Energy",

	// Telecommunications
	"Telecommunications",
	"Internet Service Provider",

	// Agriculture
	"Agriculture",
	"Fishery",
	"Forestry",
	"Agribusiness",

	// Government & Non-Profit
	"Government",
	"Non-Profit Organization",
	"NGO",
	"International Organization",

	// Other
	"Other",
}

// CompanySizeList provides standardized company size categories
var CompanySizeList = []struct {
	Code        string
	Description string
	MinEmployee int
	MaxEmployee *int
}{
	{Code: "1-10", Description: "1-10 employees (Startup/Small)", MinEmployee: 1, MaxEmployee: ptr(10)},
	{Code: "11-50", Description: "11-50 employees (Small)", MinEmployee: 11, MaxEmployee: ptr(50)},
	{Code: "51-200", Description: "51-200 employees (Medium)", MinEmployee: 51, MaxEmployee: ptr(200)},
	{Code: "201-1000", Description: "201-1000 employees (Large)", MinEmployee: 201, MaxEmployee: ptr(1000)},
	{Code: "1000+", Description: "1000+ employees (Enterprise)", MinEmployee: 1000, MaxEmployee: nil},
}

// CompanyTypeList provides standardized company types
var CompanyTypeList = []struct {
	Code        string
	Description string
}{
	{Code: "private", Description: "Private Company (PT)"},
	{Code: "public", Description: "Public Company (Tbk)"},
	{Code: "startup", Description: "Startup"},
	{Code: "ngo", Description: "Non-Governmental Organization"},
	{Code: "government", Description: "Government Agency"},
}

// IndonesianCities provides major cities in Indonesia
var IndonesianCities = []struct {
	City     string
	Province string
	Region   string
}{
	// Java
	{City: "Jakarta", Province: "DKI Jakarta", Region: "Java"},
	{City: "Jakarta Selatan", Province: "DKI Jakarta", Region: "Java"},
	{City: "Jakarta Pusat", Province: "DKI Jakarta", Region: "Java"},
	{City: "Jakarta Barat", Province: "DKI Jakarta", Region: "Java"},
	{City: "Jakarta Utara", Province: "DKI Jakarta", Region: "Java"},
	{City: "Jakarta Timur", Province: "DKI Jakarta", Region: "Java"},
	{City: "Tangerang", Province: "Banten", Region: "Java"},
	{City: "Tangerang Selatan", Province: "Banten", Region: "Java"},
	{City: "Bekasi", Province: "Jawa Barat", Region: "Java"},
	{City: "Depok", Province: "Jawa Barat", Region: "Java"},
	{City: "Bogor", Province: "Jawa Barat", Region: "Java"},
	{City: "Bandung", Province: "Jawa Barat", Region: "Java"},
	{City: "Semarang", Province: "Jawa Tengah", Region: "Java"},
	{City: "Yogyakarta", Province: "DI Yogyakarta", Region: "Java"},
	{City: "Surabaya", Province: "Jawa Timur", Region: "Java"},
	{City: "Malang", Province: "Jawa Timur", Region: "Java"},

	// Sumatra
	{City: "Medan", Province: "Sumatera Utara", Region: "Sumatra"},
	{City: "Palembang", Province: "Sumatera Selatan", Region: "Sumatra"},
	{City: "Batam", Province: "Kepulauan Riau", Region: "Sumatra"},
	{City: "Pekanbaru", Province: "Riau", Region: "Sumatra"},
	{City: "Padang", Province: "Sumatera Barat", Region: "Sumatra"},
	{City: "Bandar Lampung", Province: "Lampung", Region: "Sumatra"},

	// Kalimantan
	{City: "Balikpapan", Province: "Kalimantan Timur", Region: "Kalimantan"},
	{City: "Samarinda", Province: "Kalimantan Timur", Region: "Kalimantan"},
	{City: "Pontianak", Province: "Kalimantan Barat", Region: "Kalimantan"},
	{City: "Banjarmasin", Province: "Kalimantan Selatan", Region: "Kalimantan"},

	// Sulawesi
	{City: "Makassar", Province: "Sulawesi Selatan", Region: "Sulawesi"},
	{City: "Manado", Province: "Sulawesi Utara", Region: "Sulawesi"},

	// Bali & Nusa Tenggara
	{City: "Denpasar", Province: "Bali", Region: "Bali"},
	{City: "Mataram", Province: "Nusa Tenggara Barat", Region: "Nusa Tenggara"},

	// Papua & Maluku
	{City: "Jayapura", Province: "Papua", Region: "Papua"},
	{City: "Ambon", Province: "Maluku", Region: "Maluku"},
}

// ExperienceLevelList provides standardized experience levels for jobs
var ExperienceLevelList = []struct {
	Code        string
	Description string
	MinYears    int
	MaxYears    *int
}{
	{Code: "entry", Description: "Entry Level / Fresh Graduate", MinYears: 0, MaxYears: ptr(1)},
	{Code: "junior", Description: "Junior (1-3 years)", MinYears: 1, MaxYears: ptr(3)},
	{Code: "mid", Description: "Mid-Level (3-5 years)", MinYears: 3, MaxYears: ptr(5)},
	{Code: "senior", Description: "Senior (5-8 years)", MinYears: 5, MaxYears: ptr(8)},
	{Code: "lead", Description: "Lead / Manager (8+ years)", MinYears: 8, MaxYears: nil},
	{Code: "executive", Description: "Executive / C-Level (10+ years)", MinYears: 10, MaxYears: nil},
}

// SalaryRanges provides common salary ranges in Indonesia (in IDR millions per month)
var SalaryRanges = []struct {
	Min         int
	Max         int
	Description string
	Level       string
}{
	{Min: 3, Max: 5, Description: "3-5 juta/bulan", Level: "entry"},
	{Min: 5, Max: 8, Description: "5-8 juta/bulan", Level: "entry"},
	{Min: 8, Max: 12, Description: "8-12 juta/bulan", Level: "junior"},
	{Min: 12, Max: 18, Description: "12-18 juta/bulan", Level: "junior"},
	{Min: 18, Max: 25, Description: "18-25 juta/bulan", Level: "mid"},
	{Min: 25, Max: 35, Description: "25-35 juta/bulan", Level: "senior"},
	{Min: 35, Max: 50, Description: "35-50 juta/bulan", Level: "senior"},
	{Min: 50, Max: 75, Description: "50-75 juta/bulan", Level: "lead"},
	{Min: 75, Max: 100, Description: "75-100 juta/bulan", Level: "lead"},
	{Min: 100, Max: 0, Description: "100+ juta/bulan", Level: "executive"},
}

// Helper function to create int pointer
func ptr(i int) *int {
	return &i
}
