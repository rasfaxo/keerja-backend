package email

import (
	"bytes"
	"fmt"
	"html/template"
)

// EmailTemplate represents email template types
type EmailTemplate string

const (
	TemplateVerification      EmailTemplate = "verification"
	TemplateForgotPassword    EmailTemplate = "forgot_password"
	TemplatePasswordReset     EmailTemplate = "password_reset"
	TemplateWelcome           EmailTemplate = "welcome"
	TemplateApplicationUpdate EmailTemplate = "application_update"
	TemplateInterviewInvite   EmailTemplate = "interview_invite"
	TemplateJobAlert          EmailTemplate = "job_alert"
	TemplateCompanyVerified   EmailTemplate = "company_verified"
	TemplateOTP               EmailTemplate = "otp"
	TemplateOTPRegistration   EmailTemplate = "otp_registration"
)

// TemplateData holds data for email templates
type TemplateData struct {
	Name          string
	Email         string
	Token         string
	VerifyURL     string
	ResetURL      string
	LoginURL      string
	DashboardURL  string
	SupportEmail  string
	CompanyName   string
	JobTitle      string
	InterviewDate string
	InterviewTime string
	InterviewURL  string
	ApplicationID string
	Status        string
	Message       string
	OTPCode       string
	Purpose       string
	ExpiryMinutes string
	Year          int
}

// Templates stores HTML templates
var templates = map[EmailTemplate]string{
	TemplateVerification: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verifikasi Email Anda</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">Selamat Datang di Keerja!</h2>
        <p>Halo {{.Name}},</p>
        <p>Terima kasih telah mendaftar di Keerja. Silakan verifikasi email Anda dengan mengklik tombol di bawah ini:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.VerifyURL}}" style="background-color: #4CAF50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
                Verifikasi Email
            </a>
        </div>
        <p>Atau salin dan tempel link berikut di browser Anda:</p>
        <p style="word-break: break-all; color: #666;">{{.VerifyURL}}</p>
        <p>Link verifikasi ini akan kadaluarsa dalam 24 jam.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            Jika Anda tidak mendaftar di Keerja, abaikan email ini.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}
        </p>
    </div>
</body>
</html>
`,

	TemplateForgotPassword: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Password Anda</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #FF9800;">Reset Password</h2>
        <p>Halo {{.Name}},</p>
        <p>Kami menerima permintaan untuk reset password akun Anda. Klik tombol di bawah untuk membuat password baru:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.ResetURL}}" style="background-color: #FF9800; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
                Reset Password
            </a>
        </div>
        <p>Atau salin dan tempel link berikut di browser Anda:</p>
        <p style="word-break: break-all; color: #666;">{{.ResetURL}}</p>
        <p>Link reset password ini akan kadaluarsa dalam 1 jam.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            Jika Anda tidak meminta reset password, abaikan email ini. Password Anda tetap aman.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}
        </p>
    </div>
</body>
</html>
`,

	TemplateWelcome: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Selamat Datang di Keerja!</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">Selamat Datang di Keerja!</h2>
        <p>Halo {{.Name}},</p>
        <p>Email Anda telah berhasil diverifikasi! Sekarang Anda dapat mulai menggunakan Keerja untuk mencari pekerjaan impian Anda.</p>
        <h3>Langkah Selanjutnya:</h3>
        <ul>
            <li>Lengkapi profil Anda untuk meningkatkan visibilitas</li>
            <li>Upload CV dan dokumen pendukung</li>
            <li>Mulai cari dan lamar pekerjaan</li>
            <li>Set job alert untuk notifikasi pekerjaan baru</li>
        </ul>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.DashboardURL}}" style="background-color: #4CAF50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
                Ke Dashboard
            </a>
        </div>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            © {{.Year}} Keerja. All rights reserved.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}
        </p>
    </div>
</body>
</html>
`,

	TemplateApplicationUpdate: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Update Status Lamaran</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2196F3;">Update Status Lamaran</h2>
        <p>Halo {{.Name}},</p>
        <p>Status lamaran Anda untuk posisi <strong>{{.JobTitle}}</strong> di <strong>{{.CompanyName}}</strong> telah diupdate.</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <p style="margin: 0;"><strong>ID Lamaran:</strong> {{.ApplicationID}}</p>
            <p style="margin: 10px 0;"><strong>Status:</strong> <span style="color: #2196F3;">{{.Status}}</span></p>
            <p style="margin: 10px 0 0 0;"><strong>Pesan:</strong> {{.Message}}</p>
        </div>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.DashboardURL}}" style="background-color: #2196F3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
                Lihat Detail
            </a>
        </div>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            © {{.Year}} Keerja. All rights reserved.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}
        </p>
    </div>
</body>
</html>
`,

	TemplateInterviewInvite: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Undangan Interview</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #9C27B0;">Undangan Interview</h2>
        <p>Halo {{.Name}},</p>
        <p>Selamat! Anda telah dipilih untuk interview untuk posisi <strong>{{.JobTitle}}</strong> di <strong>{{.CompanyName}}</strong>.</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <p style="margin: 0;"><strong>Tanggal:</strong> {{.InterviewDate}}</p>
            <p style="margin: 10px 0;"><strong>Waktu:</strong> {{.InterviewTime}}</p>
            <p style="margin: 10px 0 0 0;"><strong>Link Interview:</strong> <a href="{{.InterviewURL}}" style="color: #2196F3;">{{.InterviewURL}}</a></p>
        </div>
        <p>{{.Message}}</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.InterviewURL}}" style="background-color: #9C27B0; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
                Join Interview
            </a>
        </div>
        <p>Pastikan Anda sudah siap dan hadir tepat waktu. Good luck!</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            © {{.Year}} Keerja. All rights reserved.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}
        </p>
    </div>
</body>
</html>
`,

	TemplateOTP: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Kode OTP Anda</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">Kode Verifikasi OTP</h2>
        <p>Halo {{.Name}},</p>
        <p>Berikut adalah kode OTP untuk {{.Purpose}}:</p>
        <div style="background-color: #f5f5f5; padding: 20px; border-radius: 4px; margin: 30px 0; text-align: center;">
            <h1 style="margin: 0; font-size: 48px; letter-spacing: 10px; color: #4CAF50;">{{.OTPCode}}</h1>
        </div>
        <p><strong>Kode ini akan kadaluarsa dalam {{.ExpiryMinutes}} menit.</strong></p>
        <p>Jangan bagikan kode ini kepada siapapun. Tim Keerja tidak akan pernah meminta kode OTP Anda.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            Jika Anda tidak meminta kode OTP ini, abaikan email ini.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}
        </p>
    </div>
</body>
</html>
`,

	TemplateOTPRegistration: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verifikasi Email Registrasi</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">Selamat Datang di Keerja!</h2>
        <p>Halo {{.Name}},</p>
        <p>Terima kasih telah mendaftar di Keerja. Untuk menyelesaikan registrasi, silakan verifikasi email Anda dengan memasukkan kode OTP berikut:</p>
        <div style="background-color: #f5f5f5; padding: 30px; border-radius: 8px; margin: 30px 0; text-align: center; border: 2px solid #4CAF50;">
            <h1 style="margin: 0; font-size: 56px; letter-spacing: 15px; color: #4CAF50; font-weight: bold;">{{.OTPCode}}</h1>
        </div>
        <p style="text-align: center; font-size: 14px; color: #666;"><strong>Kode OTP ini akan kadaluarsa dalam {{.ExpiryMinutes}} menit.</strong></p>
        
        <div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 4px;">
            <p style="margin: 0; color: #856404;"><strong>⚠️ Perhatian Keamanan:</strong></p>
            <ul style="margin: 10px 0 0 0; color: #856404;">
                <li>Jangan bagikan kode OTP ini kepada siapapun</li>
                <li>Tim Keerja tidak akan pernah meminta kode OTP Anda</li>
                <li>Kode ini hanya valid untuk satu kali penggunaan</li>
            </ul>
        </div>

        <p>Setelah verifikasi berhasil, Anda dapat mulai:</p>
        <ul>
            <li>📝 Lengkapi profil Anda</li>
            <li>📄 Upload CV dan dokumen pendukung</li>
            <li>🔍 Cari dan lamar pekerjaan impian</li>
            <li>🔔 Aktifkan job alert untuk notifikasi lowongan terbaru</li>
        </ul>

        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            Jika Anda tidak mendaftar di Keerja, abaikan email ini.<br>
            Butuh bantuan? Hubungi kami di {{.SupportEmail}}<br>
            © {{.Year}} Keerja. All rights reserved.
        </p>
    </div>
</body>
</html>
`,
}

// RenderTemplate renders email template with data
func RenderTemplate(templateType EmailTemplate, data TemplateData) (string, error) {
	tmplStr, ok := templates[templateType]
	if !ok {
		return "", fmt.Errorf("template %s not found", templateType)
	}

	tmpl, err := template.New(string(templateType)).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// GetSubject returns email subject based on template type
func GetSubject(templateType EmailTemplate) string {
	subjects := map[EmailTemplate]string{
		TemplateVerification:      "Verifikasi Email Anda - Keerja",
		TemplateForgotPassword:    "Reset Password Akun Keerja Anda",
		TemplatePasswordReset:     "Password Anda Telah Direset",
		TemplateWelcome:           "Selamat Datang di Keerja!",
		TemplateApplicationUpdate: "Update Status Lamaran Pekerjaan",
		TemplateInterviewInvite:   "Undangan Interview - Keerja",
		TemplateJobAlert:          "Job Alert: Pekerjaan Baru Sesuai Preferensi Anda",
		TemplateCompanyVerified:   "Perusahaan Anda Telah Terverifikasi",
		TemplateOTP:               "Kode OTP Verifikasi - Keerja",
		TemplateOTPRegistration:   "Verifikasi Email Registrasi - Keerja",
	}

	if subject, ok := subjects[templateType]; ok {
		return subject
	}
	return "Notifikasi dari Keerja"
}
