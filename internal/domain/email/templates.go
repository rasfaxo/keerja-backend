package email

import (
	"bytes"
	"fmt"
	"html/template"
)

// EmailTemplate represents email template types
type EmailTemplate string

const (
	TemplateVerification       EmailTemplate = "verification"
	TemplateForgotPassword     EmailTemplate = "forgot_password"
	TemplatePasswordReset      EmailTemplate = "password_reset"
	TemplateWelcome            EmailTemplate = "welcome"
	TemplateApplicationUpdate  EmailTemplate = "application_update"
	TemplateInterviewInvite    EmailTemplate = "interview_invite"
	TemplateJobAlert           EmailTemplate = "job_alert"
	TemplateCompanyVerified    EmailTemplate = "company_verified"
	TemplateOTP                EmailTemplate = "otp"
	TemplateOTPRegistration    EmailTemplate = "otp_registration"
	TemplateCompanyInvitation  EmailTemplate = "company_invitation"
	TemplateInvitationAccepted EmailTemplate = "invitation_accepted"
	TemplateInvitationExpired  EmailTemplate = "invitation_expired"
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
	// Invitation specific fields
	InviterName string
	Position    string
	Role        string
	InviteURL   string
	ExpiryDays  string
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

	TemplateCompanyInvitation: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Undangan Bergabung ke Tim {{.CompanyName}}</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f5f7fa;">
    <div style="max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <!-- Header -->
        <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">🎉 Undangan Bergabung</h1>
            <p style="color: #e0e7ff; margin: 10px 0 0 0; font-size: 16px;">Anda diundang untuk bergabung dengan tim</p>
        </div>
        
        <!-- Content -->
        <div style="padding: 40px 30px;">
            <p style="font-size: 16px; margin: 0 0 20px 0;">Halo <strong>{{.Name}}</strong>,</p>
            
            <p style="font-size: 16px; line-height: 1.8; margin: 0 0 25px 0;">
                <strong>{{.InviterName}}</strong> dari <strong style="color: #667eea;">{{.CompanyName}}</strong> mengundang Anda untuk bergabung sebagai anggota tim mereka di platform Keerja.
            </p>

            <!-- Info Box -->
            <div style="background-color: #f8fafc; border-left: 4px solid #667eea; padding: 20px; margin: 25px 0; border-radius: 6px;">
                <table style="width: 100%; border-collapse: collapse;">
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px; width: 120px;">📍 Posisi:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px;">{{.Position}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px;">👤 Role:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px; text-transform: capitalize;">{{.Role}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px;">⏰ Berlaku hingga:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px;">{{.ExpiryDays}} hari</td>
                    </tr>
                </table>
            </div>

            <p style="font-size: 15px; color: #64748b; margin: 25px 0;">
                Klik tombol di bawah untuk menerima undangan dan mulai berkolaborasi dengan tim:
            </p>

            <!-- CTA Button -->
            <div style="text-align: center; margin: 35px 0;">
                <a href="{{.InviteURL}}" style="display: inline-block; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; padding: 16px 40px; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4); transition: all 0.3s;">
                    ✨ Terima Undangan
                </a>
            </div>

            <p style="font-size: 14px; color: #94a3b8; margin: 25px 0 10px 0;">
                Atau salin dan tempel link berikut di browser Anda:
            </p>
            <div style="background-color: #f1f5f9; padding: 12px; border-radius: 6px; word-break: break-all; font-size: 13px; color: #475569; font-family: 'Courier New', monospace;">
                {{.InviteURL}}
            </div>

            <!-- Warning Box -->
            <div style="background-color: #fef3c7; border-left: 4px solid #f59e0b; padding: 15px; margin: 30px 0 0 0; border-radius: 6px;">
                <p style="margin: 0; color: #92400e; font-size: 13px;"><strong>⚠️ Penting:</strong></p>
                <ul style="margin: 8px 0 0 0; padding-left: 20px; color: #92400e; font-size: 13px;">
                    <li>Link undangan ini akan kadaluarsa dalam {{.ExpiryDays}} hari</li>
                    <li>Link hanya dapat digunakan satu kali</li>
                    <li>Pastikan ini adalah email yang terdaftar di akun Keerja Anda</li>
                </ul>
            </div>
        </div>

        <!-- Footer -->
        <div style="background-color: #f8fafc; padding: 25px 30px; border-top: 1px solid #e2e8f0;">
            <p style="font-size: 13px; color: #94a3b8; margin: 0 0 10px 0; text-align: center;">
                Jika Anda tidak mengenal perusahaan ini atau tidak mengharapkan undangan ini, abaikan email ini.
            </p>
            <hr style="border: none; border-top: 1px solid #e2e8f0; margin: 15px 0;">
            <p style="font-size: 12px; color: #cbd5e1; margin: 0; text-align: center;">
                Butuh bantuan? Hubungi kami di <a href="mailto:{{.SupportEmail}}" style="color: #667eea; text-decoration: none;">{{.SupportEmail}}</a><br>
                © {{.Year}} Keerja. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>
`,

	TemplateInvitationAccepted: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Undangan Diterima - {{.Name}} Bergabung ke Tim</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f5f7fa;">
    <div style="max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <!-- Header -->
        <div style="background: linear-gradient(135deg, #10b981 0%, #059669 100%); padding: 40px 30px; text-align: center;">
            <div style="font-size: 48px; margin: 0 0 15px 0;">✅</div>
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Undangan Diterima!</h1>
            <p style="color: #d1fae5; margin: 10px 0 0 0; font-size: 16px;">Anggota baru telah bergabung dengan tim Anda</p>
        </div>
        
        <!-- Content -->
        <div style="padding: 40px 30px;">
            <p style="font-size: 16px; margin: 0 0 20px 0;">Halo <strong>{{.InviterName}}</strong>,</p>
            
            <p style="font-size: 16px; line-height: 1.8; margin: 0 0 25px 0;">
                Kabar baik! <strong>{{.Name}}</strong> telah menerima undangan Anda dan sekarang menjadi bagian dari tim <strong style="color: #10b981;">{{.CompanyName}}</strong>.
            </p>

            <!-- Success Box -->
            <div style="background-color: #f0fdf4; border: 2px solid #86efac; padding: 25px; margin: 25px 0; border-radius: 8px; text-align: center;">
                <div style="font-size: 40px; margin: 0 0 15px 0;">🎊</div>
                <h3 style="color: #166534; margin: 0 0 10px 0; font-size: 20px;">Selamat Datang di Tim!</h3>
                <p style="color: #166534; margin: 0; font-size: 14px;">{{.Name}} siap untuk mulai berkolaborasi</p>
            </div>

            <!-- Member Info -->
            <div style="background-color: #f8fafc; border-left: 4px solid #10b981; padding: 20px; margin: 25px 0; border-radius: 6px;">
                <h4 style="margin: 0 0 15px 0; color: #1e293b; font-size: 16px;">📋 Detail Anggota Baru:</h4>
                <table style="width: 100%; border-collapse: collapse;">
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px; width: 120px;">👤 Nama:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px;">{{.Name}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px;">📧 Email:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px;">{{.Email}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px;">📍 Posisi:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px;">{{.Position}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748b; font-size: 14px;">🔑 Role:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px; text-transform: capitalize;">{{.Role}}</td>
                    </tr>
                </table>
            </div>

            <p style="font-size: 15px; color: #64748b; margin: 25px 0;">
                Anda sekarang dapat berkolaborasi dengan anggota baru di dashboard perusahaan:
            </p>

            <!-- CTA Button -->
            <div style="text-align: center; margin: 35px 0;">
                <a href="{{.DashboardURL}}" style="display: inline-block; background: linear-gradient(135deg, #10b981 0%, #059669 100%); color: #ffffff; padding: 16px 40px; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);">
                    🚀 Buka Dashboard
                </a>
            </div>

            <!-- Tips Box -->
            <div style="background-color: #eff6ff; border-left: 4px solid #3b82f6; padding: 15px; margin: 30px 0 0 0; border-radius: 6px;">
                <p style="margin: 0; color: #1e40af; font-size: 13px;"><strong>💡 Tips:</strong></p>
                <ul style="margin: 8px 0 0 0; padding-left: 20px; color: #1e40af; font-size: 13px;">
                    <li>Perkenalkan anggota baru kepada tim</li>
                    <li>Berikan orientasi tentang tools dan prosedur</li>
                    <li>Pastikan mereka memiliki akses yang sesuai</li>
                </ul>
            </div>
        </div>

        <!-- Footer -->
        <div style="background-color: #f8fafc; padding: 25px 30px; border-top: 1px solid #e2e8f0;">
            <hr style="border: none; border-top: 1px solid #e2e8f0; margin: 0 0 15px 0;">
            <p style="font-size: 12px; color: #cbd5e1; margin: 0; text-align: center;">
                Butuh bantuan? Hubungi kami di <a href="mailto:{{.SupportEmail}}" style="color: #10b981; text-decoration: none;">{{.SupportEmail}}</a><br>
                © {{.Year}} Keerja. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>
`,

	TemplateInvitationExpired: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Undangan Kadaluarsa</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f5f7fa;">
    <div style="max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <!-- Header -->
        <div style="background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); padding: 40px 30px; text-align: center;">
            <div style="font-size: 48px; margin: 0 0 15px 0;">⏰</div>
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Undangan Kadaluarsa</h1>
            <p style="color: #fecaca; margin: 10px 0 0 0; font-size: 16px;">Link undangan tidak lagi berlaku</p>
        </div>
        
        <!-- Content -->
        <div style="padding: 40px 30px;">
            <p style="font-size: 16px; margin: 0 0 20px 0;">Halo <strong>{{.Name}}</strong>,</p>
            
            <p style="font-size: 16px; line-height: 1.8; margin: 0 0 25px 0;">
                Undangan Anda untuk bergabung dengan <strong style="color: #ef4444;">{{.CompanyName}}</strong> sebagai <strong>{{.Position}}</strong> telah kadaluarsa dan tidak dapat lagi digunakan.
            </p>

            <!-- Info Box -->
            <div style="background-color: #fef2f2; border-left: 4px solid #ef4444; padding: 20px; margin: 25px 0; border-radius: 6px;">
                <h4 style="margin: 0 0 15px 0; color: #991b1b; font-size: 16px;">📋 Detail Undangan:</h4>
                <table style="width: 100%; border-collapse: collapse;">
                    <tr>
                        <td style="padding: 8px 0; color: #991b1b; font-size: 14px; width: 140px;">🏢 Perusahaan:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px; color: #991b1b;">{{.CompanyName}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #991b1b; font-size: 14px;">📍 Posisi:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px; color: #991b1b;">{{.Position}}</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #991b1b; font-size: 14px;">❌ Status:</td>
                        <td style="padding: 8px 0; font-weight: 600; font-size: 14px; color: #991b1b;">Kadaluarsa</td>
                    </tr>
                </table>
            </div>

            <p style="font-size: 15px; color: #64748b; margin: 25px 0;">
                Jika Anda masih tertarik untuk bergabung dengan tim ini, silakan hubungi <strong>{{.InviterName}}</strong> untuk meminta undangan baru.
            </p>

            <!-- Action Box -->
            <div style="background-color: #eff6ff; border: 2px dashed #3b82f6; padding: 25px; margin: 30px 0; border-radius: 8px; text-align: center;">
                <h3 style="color: #1e40af; margin: 0 0 10px 0; font-size: 18px;">💬 Butuh Undangan Baru?</h3>
                <p style="color: #1e40af; margin: 0 0 20px 0; font-size: 14px;">
                    Hubungi admin perusahaan untuk mendapatkan link undangan yang baru
                </p>
                <a href="mailto:{{.SupportEmail}}?subject=Request%20New%20Invitation%20-%20{{.CompanyName}}" style="display: inline-block; background-color: #3b82f6; color: #ffffff; padding: 12px 30px; text-decoration: none; border-radius: 6px; font-weight: 600; font-size: 14px;">
                    📧 Hubungi Support
                </a>
            </div>

            <!-- Info Box -->
            <div style="background-color: #fefce8; border-left: 4px solid #eab308; padding: 15px; margin: 25px 0 0 0; border-radius: 6px;">
                <p style="margin: 0; color: #854d0e; font-size: 13px;"><strong>ℹ️ Informasi:</strong></p>
                <ul style="margin: 8px 0 0 0; padding-left: 20px; color: #854d0e; font-size: 13px;">
                    <li>Link undangan biasanya berlaku selama 7 hari</li>
                    <li>Setelah kadaluarsa, link tidak dapat digunakan lagi</li>
                    <li>Admin perusahaan dapat mengirim undangan baru kapan saja</li>
                </ul>
            </div>
        </div>

        <!-- Footer -->
        <div style="background-color: #f8fafc; padding: 25px 30px; border-top: 1px solid #e2e8f0;">
            <hr style="border: none; border-top: 1px solid #e2e8f0; margin: 0 0 15px 0;">
            <p style="font-size: 12px; color: #cbd5e1; margin: 0; text-align: center;">
                Butuh bantuan? Hubungi kami di <a href="mailto:{{.SupportEmail}}" style="color: #ef4444; text-decoration: none;">{{.SupportEmail}}</a><br>
                © {{.Year}} Keerja. All rights reserved.
            </p>
        </div>
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
		TemplateVerification:       "Verifikasi Email Anda - Keerja",
		TemplateForgotPassword:     "Reset Password Akun Keerja Anda",
		TemplatePasswordReset:      "Password Anda Telah Direset",
		TemplateWelcome:            "Selamat Datang di Keerja!",
		TemplateApplicationUpdate:  "Update Status Lamaran Pekerjaan",
		TemplateInterviewInvite:    "Undangan Interview - Keerja",
		TemplateJobAlert:           "Job Alert: Pekerjaan Baru Sesuai Preferensi Anda",
		TemplateCompanyVerified:    "Perusahaan Anda Telah Terverifikasi",
		TemplateOTP:                "Kode OTP Verifikasi - Keerja",
		TemplateOTPRegistration:    "Verifikasi Email Registrasi - Keerja",
		TemplateCompanyInvitation:  "Undangan Bergabung ke Tim - Keerja",
		TemplateInvitationAccepted: "Undangan Diterima - Anggota Baru Bergabung",
		TemplateInvitationExpired:  "Undangan Kadaluarsa - Keerja",
	}

	if subject, ok := subjects[templateType]; ok {
		return subject
	}
	return "Notifikasi dari Keerja"
}
