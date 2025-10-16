# Dokumentasi API Keerja Backend

## Pendahuluan

Dokumen ini berisi spesifikasi lengkap API untuk Keerja Backend yang dapat digunakan oleh tim frontend. Semua endpoint berada di bawah base path: `/api/v1`.

## Format Response

### Success Response

```json
{
  "success": true,
  "message": "Success message",
  "data": { ... } // Data bervariasi sesuai endpoint
}
```

### Error Response

```json
{
  "success": false,
  "message": "Error message",
  "errors": [ ... ] // Array detail error (opsional)
}
```

## Autentikasi

Untuk endpoint yang memerlukan autentikasi, tambahkan header:

```
Authorization: Bearer {jwt_token}
```

---

## 1. Auth Endpoints

### Register User

- **URL**: POST `/api/v1/auth/register`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "Password123!",
    "name": "John Doe",
    "phone": "+628123456789",
    "role": "job_seeker" // atau "employer"
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Registration successful. Please check your email for verification",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "job_seeker",
      "is_verified": false,
      "created_at": "2025-10-16T08:00:00Z"
    }
  }
  ```

### Login

- **URL**: POST `/api/v1/auth/login`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "Password123!"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Login successful",
    "data": {
      "user": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "email": "user@example.com",
        "name": "John Doe",
        "role": "job_seeker"
      },
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_at": "2025-10-16T09:00:00Z"
    }
  }
  ```

### Email Verification

- **URL**: POST `/api/v1/auth/verify-email`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "token": "verification_token_from_email"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Email verified successfully",
    "data": {
      "is_verified": true
    }
  }
  ```

### Forgot Password

- **URL**: POST `/api/v1/auth/forgot-password`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Password reset instructions sent to your email"
  }
  ```

### Reset Password

- **URL**: POST `/api/v1/auth/reset-password`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "token": "reset_token_from_email",
    "password": "NewPassword123!",
    "confirm_password": "NewPassword123!"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Password reset successful. Please login with your new password"
  }
  ```

### Resend Verification Email

- **URL**: POST `/api/v1/auth/resend-verification`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Verification email sent"
  }
  ```

### Refresh Token

- **URL**: POST `/api/v1/auth/refresh-token`
- **Auth**: Required
- **Body**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Token refreshed",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_at": "2025-10-16T10:00:00Z"
    }
  }
  ```

### Logout

- **URL**: POST `/api/v1/auth/logout`
- **Auth**: Required
- **Body**: Empty
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Logged out successfully"
  }
  ```

---

## 2. User Endpoints

### Get Current User Profile

- **URL**: GET `/api/v1/users/me`
- **Auth**: Required
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Profile retrieved",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe",
      "phone": "+628123456789",
      "role": "job_seeker",
      "profile": {
        "headline": "Software Engineer",
        "about": "Experienced software engineer...",
        "location": "Jakarta, Indonesia",
        "avatar_url": "https://storage.example.com/avatars/user.jpg"
      },
      "experience": [...],
      "education": [...],
      "skills": [...]
    }
  }
  ```

### Update User Profile

- **URL**: PUT `/api/v1/users/me`
- **Auth**: Required
- **Body**:
  ```json
  {
    "name": "John Smith",
    "headline": "Senior Software Engineer",
    "about": "Experienced in building scalable web applications",
    "location": "Jakarta, Indonesia",
    "phone": "+628987654321"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Profile updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Smith",
      "headline": "Senior Software Engineer",
      "about": "Experienced in building scalable web applications",
      "location": "Jakarta, Indonesia",
      "phone": "+628987654321",
      "updated_at": "2025-10-16T08:30:00Z"
    }
  }
  ```

### Add Education

- **URL**: POST `/api/v1/users/me/education`
- **Auth**: Required
- **Body**:
  ```json
  {
    "institution": "University of Indonesia",
    "degree": "Bachelor of Computer Science",
    "field_of_study": "Computer Science",
    "start_date": "2017-08-01",
    "end_date": "2021-06-30",
    "grade": "3.8",
    "description": "Focus on software engineering"
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Education added successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "institution": "University of Indonesia",
      "degree": "Bachelor of Computer Science",
      "field_of_study": "Computer Science",
      "start_date": "2017-08-01",
      "end_date": "2021-06-30",
      "grade": "3.8",
      "description": "Focus on software engineering",
      "created_at": "2025-10-16T08:35:00Z"
    }
  }
  ```

### Update Education

- **URL**: PUT `/api/v1/users/me/education/:id`
- **Auth**: Required
- **Params**: id (Education ID)
- **Body**:
  ```json
  {
    "institution": "University of Indonesia",
    "degree": "Bachelor of Computer Science",
    "field_of_study": "Computer Science",
    "start_date": "2017-08-01",
    "end_date": "2021-07-15",
    "grade": "3.9",
    "description": "Focus on software engineering and AI"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Education updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "institution": "University of Indonesia",
      "degree": "Bachelor of Computer Science",
      "field_of_study": "Computer Science",
      "start_date": "2017-08-01",
      "end_date": "2021-07-15",
      "grade": "3.9",
      "description": "Focus on software engineering and AI",
      "updated_at": "2025-10-16T08:40:00Z"
    }
  }
  ```

### Delete Education

- **URL**: DELETE `/api/v1/users/me/education/:id`
- **Auth**: Required
- **Params**: id (Education ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Education deleted successfully"
  }
  ```

### Add Experience

- **URL**: POST `/api/v1/users/me/experience`
- **Auth**: Required
- **Body**:
  ```json
  {
    "company": "Tech Company",
    "title": "Software Engineer",
    "location": "Jakarta, Indonesia",
    "is_current": true,
    "start_date": "2021-07-01",
    "end_date": null,
    "description": "Developing web applications using React and Node.js"
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Experience added successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "company": "Tech Company",
      "title": "Software Engineer",
      "location": "Jakarta, Indonesia",
      "is_current": true,
      "start_date": "2021-07-01",
      "end_date": null,
      "description": "Developing web applications using React and Node.js",
      "created_at": "2025-10-16T08:45:00Z"
    }
  }
  ```

### Update Experience

- **URL**: PUT `/api/v1/users/me/experience/:id`
- **Auth**: Required
- **Params**: id (Experience ID)
- **Body**:
  ```json
  {
    "company": "Tech Company",
    "title": "Senior Software Engineer",
    "location": "Jakarta, Indonesia",
    "is_current": true,
    "start_date": "2021-07-01",
    "end_date": null,
    "description": "Leading development of web applications using React and Node.js"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Experience updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "company": "Tech Company",
      "title": "Senior Software Engineer",
      "location": "Jakarta, Indonesia",
      "is_current": true,
      "start_date": "2021-07-01",
      "end_date": null,
      "description": "Leading development of web applications using React and Node.js",
      "updated_at": "2025-10-16T08:50:00Z"
    }
  }
  ```

### Delete Experience

- **URL**: DELETE `/api/v1/users/me/experience/:id`
- **Auth**: Required
- **Params**: id (Experience ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Experience deleted successfully"
  }
  ```

### Add Skill

- **URL**: POST `/api/v1/users/me/skills`
- **Auth**: Required
- **Body**:
  ```json
  {
    "skill_id": 1, // ID dari skills_master
    "proficiency_level": "advanced" // beginner, intermediate, advanced, expert
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Skill added successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "skill": {
        "id": 1,
        "name": "JavaScript",
        "category": "Programming Language"
      },
      "proficiency_level": "advanced",
      "created_at": "2025-10-16T08:55:00Z"
    }
  }
  ```

### Delete Skill

- **URL**: DELETE `/api/v1/users/me/skills/:id`
- **Auth**: Required
- **Params**: id (Skill ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Skill deleted successfully"
  }
  ```

### Upload Document

- **URL**: POST `/api/v1/users/me/documents`
- **Auth**: Required
- **Body**: Multipart form-data
  - file: (Binary file data)
  - document_type: "resume" | "cover_letter" | "certificate" | "other"
  - title: "My Resume"
  - description: "Updated resume 2025"
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Document uploaded successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "title": "My Resume",
      "document_type": "resume",
      "description": "Updated resume 2025",
      "file_url": "https://storage.example.com/documents/resume.pdf",
      "file_size": 1024000,
      "file_type": "application/pdf",
      "created_at": "2025-10-16T09:00:00Z"
    }
  }
  ```

---

## 3. Job Endpoints

### List Jobs

- **URL**: GET `/api/v1/jobs`
- **Auth**: Tidak perlu
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - location: string (optional)
  - job_type: string (optional) - "full_time", "part_time", "contract", "internship"
  - experience_level: string (optional) - "entry", "mid", "senior", "executive"
  - salary_min: number (optional)
  - salary_max: number (optional)
  - sort: string (optional) - "newest", "salary_high", "salary_low"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Jobs retrieved successfully",
    "data": {
      "jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440010",
          "title": "Software Engineer",
          "company": {
            "id": "550e8400-e29b-41d4-a716-446655440020",
            "name": "Tech Company",
            "logo_url": "https://storage.example.com/logos/company.jpg",
            "is_verified": true
          },
          "location": "Jakarta, Indonesia",
          "job_type": "full_time",
          "salary_min": 10000000,
          "salary_max": 15000000,
          "is_salary_hidden": false,
          "posted_at": "2025-10-10T08:00:00Z",
          "application_count": 12
        }
        // more jobs...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 120,
        "total_pages": 12
      }
    }
  }
  ```

### Get Job Details

- **URL**: GET `/api/v1/jobs/:id`
- **Auth**: Tidak perlu
- **Params**: id (Job ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job details retrieved",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "title": "Software Engineer",
      "description": "We are looking for a software engineer...",
      "requirements": "- 3+ years of experience\n- Proficiency in JavaScript\n- Experience with React",
      "responsibilities": "- Develop web applications\n- Collaborate with the team\n- Write clean code",
      "benefits": [
        {
          "id": 1,
          "name": "Health Insurance",
          "icon": "health"
        },
        {
          "id": 2,
          "name": "Remote Work",
          "icon": "remote"
        }
      ],
      "skills": [
        {
          "id": 1,
          "name": "JavaScript"
        },
        {
          "id": 2,
          "name": "React"
        }
      ],
      "company": {
        "id": "550e8400-e29b-41d4-a716-446655440020",
        "name": "Tech Company",
        "description": "Leading tech company...",
        "logo_url": "https://storage.example.com/logos/company.jpg",
        "website": "https://techcompany.com",
        "industry": "Software Development",
        "company_size": "51-200",
        "founded_year": 2015,
        "is_verified": true
      },
      "location": "Jakarta, Indonesia",
      "is_remote": false,
      "job_type": "full_time",
      "experience_level": "mid",
      "salary_min": 10000000,
      "salary_max": 15000000,
      "is_salary_hidden": false,
      "application_deadline": "2025-11-10T23:59:59Z",
      "posted_at": "2025-10-10T08:00:00Z",
      "application_count": 12,
      "has_applied": false // if user is authenticated
    }
  }
  ```

### Search Jobs

- **URL**: POST `/api/v1/jobs/search`
- **Auth**: Tidak perlu
- **Body**:
  ```json
  {
    "query": "software engineer",
    "location": "Jakarta",
    "job_type": ["full_time", "contract"],
    "experience_level": ["mid", "senior"],
    "salary_min": 8000000,
    "salary_max": 20000000,
    "skills": [1, 2, 3],
    "page": 1,
    "limit": 10,
    "sort": "newest"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Search results",
    "data": {
      "jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440010",
          "title": "Senior Software Engineer",
          "company": {
            "id": "550e8400-e29b-41d4-a716-446655440020",
            "name": "Tech Company",
            "logo_url": "https://storage.example.com/logos/company.jpg",
            "is_verified": true
          },
          "location": "Jakarta, Indonesia",
          "job_type": "full_time",
          "salary_min": 12000000,
          "salary_max": 18000000,
          "is_salary_hidden": false,
          "posted_at": "2025-10-12T08:00:00Z",
          "application_count": 8,
          "relevance_score": 0.95
        }
        // more jobs...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 45,
        "total_pages": 5
      }
    }
  }
  ```

### Create Job (Employer)

- **URL**: POST `/api/v1/jobs`
- **Auth**: Required (Employer only)
- **Body**:
  ```json
  {
    "title": "Frontend Developer",
    "description": "We are looking for a frontend developer...",
    "requirements": "- 2+ years of experience\n- Proficiency in React\n- Understanding of UI/UX principles",
    "responsibilities": "- Develop user interfaces\n- Collaborate with designers\n- Write clean code",
    "company_id": "550e8400-e29b-41d4-a716-446655440020",
    "location": "Jakarta, Indonesia",
    "is_remote": false,
    "job_type": "full_time",
    "experience_level": "mid",
    "salary_min": 8000000,
    "salary_max": 12000000,
    "is_salary_hidden": false,
    "application_deadline": "2025-11-15T23:59:59Z",
    "benefit_ids": [1, 2, 3],
    "skill_ids": [1, 2, 3]
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Job created successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440011",
      "title": "Frontend Developer",
      "description": "We are looking for a frontend developer...",
      "company_id": "550e8400-e29b-41d4-a716-446655440020",
      "status": "active",
      "posted_at": "2025-10-16T09:30:00Z",
      "created_at": "2025-10-16T09:30:00Z"
    }
  }
  ```

### Update Job (Employer)

- **URL**: PUT `/api/v1/jobs/:id`
- **Auth**: Required (Employer only)
- **Params**: id (Job ID)
- **Body**:
  ```json
  {
    "title": "Senior Frontend Developer",
    "description": "We are looking for a senior frontend developer...",
    "requirements": "- 4+ years of experience\n- Advanced proficiency in React\n- Understanding of UI/UX principles",
    "responsibilities": "- Lead frontend development\n- Collaborate with designers\n- Write clean code",
    "location": "Jakarta, Indonesia",
    "is_remote": true,
    "job_type": "full_time",
    "experience_level": "senior",
    "salary_min": 15000000,
    "salary_max": 20000000,
    "is_salary_hidden": false,
    "application_deadline": "2025-11-15T23:59:59Z",
    "benefit_ids": [1, 2, 3, 4],
    "skill_ids": [1, 2, 3, 4]
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440011",
      "title": "Senior Frontend Developer",
      "description": "We are looking for a senior frontend developer...",
      "updated_at": "2025-10-16T09:35:00Z"
    }
  }
  ```

### Delete Job (Employer)

- **URL**: DELETE `/api/v1/jobs/:id`
- **Auth**: Required (Employer only)
- **Params**: id (Job ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job deleted successfully"
  }
  ```

### Get My Jobs (Employer)

- **URL**: GET `/api/v1/jobs/my-jobs`
- **Auth**: Required (Employer only)
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - status: string (optional) - "active", "closed", "draft", "expired"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "My jobs retrieved successfully",
    "data": {
      "jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440011",
          "title": "Senior Frontend Developer",
          "status": "active",
          "application_count": 5,
          "posted_at": "2025-10-16T09:30:00Z",
          "application_deadline": "2025-11-15T23:59:59Z",
          "is_expired": false
        }
        // more jobs...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 5,
        "total_pages": 1
      },
      "summary": {
        "total_jobs": 5,
        "active_jobs": 3,
        "closed_jobs": 1,
        "draft_jobs": 0,
        "expired_jobs": 1
      }
    }
  }
  ```

### Get Job Applications (Employer)

- **URL**: GET `/api/v1/jobs/:id/applications`
- **Auth**: Required (Employer only)
- **Params**: id (Job ID)
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - stage: string (optional) - "applied", "screening", "interview", "offer", "rejected"
  - sort: string (optional) - "newest", "oldest", "rating"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job applications retrieved successfully",
    "data": {
      "applications": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440030",
          "user": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "name": "John Doe",
            "headline": "Software Engineer",
            "avatar_url": "https://storage.example.com/avatars/user.jpg"
          },
          "stage": "screening",
          "applied_at": "2025-10-12T10:00:00Z",
          "last_updated_at": "2025-10-14T14:30:00Z",
          "rating": 4,
          "has_resume": true
        }
        // more applications...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 5,
        "total_pages": 1
      },
      "summary": {
        "total_applications": 5,
        "stages": {
          "applied": 2,
          "screening": 1,
          "interview": 1,
          "offer": 1,
          "rejected": 0
        }
      }
    }
  }
  ```

---

## 4. Application Endpoints

### Apply for Job

- **URL**: POST `/api/v1/applications/jobs/:id/apply`
- **Auth**: Required (Job Seeker only)
- **Params**: id (Job ID)
- **Body**:
  ```json
  {
    "resume_id": "550e8400-e29b-41d4-a716-446655440004", // Document ID
    "cover_letter_id": "550e8400-e29b-41d4-a716-446655440005", // Optional
    "phone": "+628123456789",
    "email": "john@example.com",
    "note": "I am very interested in this position because..."
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Application submitted successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "job": {
        "id": "550e8400-e29b-41d4-a716-446655440010",
        "title": "Software Engineer",
        "company": {
          "name": "Tech Company",
          "logo_url": "https://storage.example.com/logos/company.jpg"
        }
      },
      "stage": "applied",
      "applied_at": "2025-10-16T10:00:00Z"
    }
  }
  ```

### Get My Applications

- **URL**: GET `/api/v1/applications/my-applications`
- **Auth**: Required
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - status: string (optional) - "active", "completed", "rejected"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Applications retrieved successfully",
    "data": {
      "applications": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440030",
          "job": {
            "id": "550e8400-e29b-41d4-a716-446655440010",
            "title": "Software Engineer",
            "company": {
              "name": "Tech Company",
              "logo_url": "https://storage.example.com/logos/company.jpg"
            }
          },
          "stage": "screening",
          "status": "active",
          "applied_at": "2025-10-12T10:00:00Z",
          "last_updated_at": "2025-10-14T14:30:00Z"
        }
        // more applications...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 3,
        "total_pages": 1
      }
    }
  }
  ```

### Get Application Details

- **URL**: GET `/api/v1/applications/:id`
- **Auth**: Required
- **Params**: id (Application ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Application details retrieved",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "job": {
        "id": "550e8400-e29b-41d4-a716-446655440010",
        "title": "Software Engineer",
        "company": {
          "id": "550e8400-e29b-41d4-a716-446655440020",
          "name": "Tech Company",
          "logo_url": "https://storage.example.com/logos/company.jpg"
        },
        "location": "Jakarta, Indonesia",
        "job_type": "full_time"
      },
      "stage": "screening",
      "status": "active",
      "applied_at": "2025-10-12T10:00:00Z",
      "last_updated_at": "2025-10-14T14:30:00Z",
      "documents": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440004",
          "title": "My Resume",
          "document_type": "resume",
          "file_url": "https://storage.example.com/documents/resume.pdf"
        },
        {
          "id": "550e8400-e29b-41d4-a716-446655440005",
          "title": "Cover Letter",
          "document_type": "cover_letter",
          "file_url": "https://storage.example.com/documents/cover_letter.pdf"
        }
      ],
      "timeline": [
        {
          "stage": "applied",
          "timestamp": "2025-10-12T10:00:00Z",
          "note": "Application submitted"
        },
        {
          "stage": "screening",
          "timestamp": "2025-10-14T14:30:00Z",
          "note": "Application is being reviewed"
        }
      ],
      "interviews": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440040",
          "schedule": "2025-10-20T13:00:00Z",
          "duration": 60,
          "type": "video",
          "location": "https://meet.google.com/abc-defg-hij",
          "note": "Technical interview with the team lead"
        }
      ]
    }
  }
  ```

### Withdraw Application

- **URL**: DELETE `/api/v1/applications/:id/withdraw`
- **Auth**: Required (Job Seeker only)
- **Params**: id (Application ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Application withdrawn successfully"
  }
  ```

### Update Application Stage (Employer)

- **URL**: PUT `/api/v1/applications/:id/stage`
- **Auth**: Required (Employer only)
- **Params**: id (Application ID)
- **Body**:
  ```json
  {
    "stage": "interview",
    "note": "Moving to interview stage after successful screening"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Application stage updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "stage": "interview",
      "updated_at": "2025-10-16T10:30:00Z"
    }
  }
  ```

### Add Application Note (Employer)

- **URL**: POST `/api/v1/applications/:id/notes`
- **Auth**: Required (Employer only)
- **Params**: id (Application ID)
- **Body**:
  ```json
  {
    "note": "Candidate shows strong technical skills but needs improvement in communication",
    "visibility": "internal" // internal or shared
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Note added successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440050",
      "note": "Candidate shows strong technical skills but needs improvement in communication",
      "visibility": "internal",
      "created_at": "2025-10-16T10:35:00Z",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440060",
        "name": "Recruiter Name"
      }
    }
  }
  ```

### Schedule Interview (Employer)

- **URL**: POST `/api/v1/applications/:id/schedule-interview`
- **Auth**: Required (Employer only)
- **Params**: id (Application ID)
- **Body**:
  ```json
  {
    "schedule": "2025-10-25T14:00:00Z",
    "duration": 60, // minutes
    "type": "video", // in_person, video, phone
    "location": "https://meet.google.com/abc-defg-hij", // or physical location for in_person
    "note": "Technical interview with the engineering team",
    "interviewer_ids": [
      "550e8400-e29b-41d4-a716-446655440065",
      "550e8400-e29b-41d4-a716-446655440066"
    ]
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Interview scheduled successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440040",
      "schedule": "2025-10-25T14:00:00Z",
      "duration": 60,
      "type": "video",
      "location": "https://meet.google.com/abc-defg-hij",
      "note": "Technical interview with the engineering team",
      "created_at": "2025-10-16T10:40:00Z"
    }
  }
  ```

---

## 5. Company Endpoints

### List Companies

- **URL**: GET `/api/v1/companies`
- **Auth**: Tidak perlu
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - industry: string (optional)
  - location: string (optional)
  - size: string (optional) - "1-10", "11-50", "51-200", "201-500", "501-1000", "1000+"
  - sort: string (optional) - "newest", "oldest", "rating", "job_count"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Companies retrieved successfully",
    "data": {
      "companies": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440020",
          "name": "Tech Company",
          "logo_url": "https://storage.example.com/logos/company.jpg",
          "industry": "Software Development",
          "location": "Jakarta, Indonesia",
          "company_size": "51-200",
          "is_verified": true,
          "active_job_count": 8,
          "rating": 4.5
        }
        // more companies...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 50,
        "total_pages": 5
      }
    }
  }
  ```

### Get Company Details

- **URL**: GET `/api/v1/companies/:id`
- **Auth**: Tidak perlu
- **Params**: id (Company ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company details retrieved",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "name": "Tech Company",
      "description": "Leading tech company specializing in software development...",
      "logo_url": "https://storage.example.com/logos/company.jpg",
      "cover_url": "https://storage.example.com/covers/company_cover.jpg",
      "website": "https://techcompany.com",
      "linkedin": "https://linkedin.com/company/techcompany",
      "industry": "Software Development",
      "location": "Jakarta, Indonesia",
      "company_size": "51-200",
      "founded_year": 2015,
      "is_verified": true,
      "active_job_count": 8,
      "rating": 4.5,
      "review_count": 45,
      "followers_count": 2300,
      "is_following": false, // if user is authenticated
      "benefits": [
        {
          "id": 1,
          "name": "Health Insurance",
          "icon": "health"
        },
        {
          "id": 2,
          "name": "Remote Work",
          "icon": "remote"
        }
      ],
      "photos": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440070",
          "url": "https://storage.example.com/company_photos/office1.jpg",
          "caption": "Our modern office"
        }
        // more photos...
      ],
      "active_jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440010",
          "title": "Software Engineer",
          "location": "Jakarta, Indonesia",
          "job_type": "full_time",
          "posted_at": "2025-10-10T08:00:00Z"
        }
        // more jobs (limited to 5)...
      ]
    }
  }
  ```

### Register Company

- **URL**: POST `/api/v1/companies`
- **Auth**: Required
- **Body**:
  ```json
  {
    "name": "New Tech Company",
    "description": "We are a new tech company...",
    "website": "https://newtechcompany.com",
    "linkedin": "https://linkedin.com/company/newtechcompany",
    "industry": "Software Development",
    "location": "Bandung, Indonesia",
    "company_size": "11-50",
    "founded_year": 2022
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Company registered successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440021",
      "name": "New Tech Company",
      "is_verified": false,
      "created_at": "2025-10-16T11:00:00Z"
    }
  }
  ```

### Update Company

- **URL**: PUT `/api/v1/companies/:id`
- **Auth**: Required
- **Params**: id (Company ID)
- **Body**:
  ```json
  {
    "name": "Updated Tech Company",
    "description": "We are an updated tech company...",
    "website": "https://updatedtechcompany.com",
    "linkedin": "https://linkedin.com/company/updatedtechcompany",
    "industry": "Software Development",
    "location": "Bandung, Indonesia",
    "company_size": "11-50",
    "founded_year": 2022
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440021",
      "name": "Updated Tech Company",
      "updated_at": "2025-10-16T11:05:00Z"
    }
  }
  ```

### Follow Company

- **URL**: POST `/api/v1/companies/:id/follow`
- **Auth**: Required
- **Params**: id (Company ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company followed successfully",
    "data": {
      "is_following": true,
      "followers_count": 2301
    }
  }
  ```

### Unfollow Company

- **URL**: DELETE `/api/v1/companies/:id/follow`
- **Auth**: Required
- **Params**: id (Company ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company unfollowed successfully",
    "data": {
      "is_following": false,
      "followers_count": 2300
    }
  }
  ```

### Add Company Review

- **URL**: POST `/api/v1/companies/:id/review`
- **Auth**: Required
- **Params**: id (Company ID)
- **Body**:
  ```json
  {
    "rating": 4,
    "title": "Great place to work",
    "pros": "Good work-life balance, great colleagues",
    "cons": "Limited career growth opportunities",
    "is_anonymous": true,
    "employment_status": "current", // current, former
    "job_title": "Software Engineer",
    "year_worked": 2024
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Review submitted successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440080",
      "rating": 4,
      "created_at": "2025-10-16T11:10:00Z"
    }
  }
  ```

### Request Company Verification

- **URL**: POST `/api/v1/companies/:id/verify`
- **Auth**: Required
- **Params**: id (Company ID)
- **Body**:
  ```json
  {
    "contact_name": "John Smith",
    "contact_email": "john@newtechcompany.com",
    "contact_phone": "+628123456789",
    "document_ids": ["550e8400-e29b-41d4-a716-446655440090"] // Company registration document
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Verification request submitted successfully",
    "data": {
      "request_id": "550e8400-e29b-41d4-a716-446655440095",
      "status": "pending",
      "submitted_at": "2025-10-16T11:15:00Z"
    }
  }
  ```

---

## 6. Admin Endpoints

### Get Admin Dashboard

- **URL**: GET `/api/v1/admin/dashboard`
- **Auth**: Required (Admin only)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Dashboard data retrieved",
    "data": {
      "users": {
        "total": 5000,
        "job_seekers": 4500,
        "employers": 500,
        "new_today": 25,
        "new_this_week": 150,
        "active_today": 1200
      },
      "companies": {
        "total": 300,
        "verified": 220,
        "pending_verification": 30,
        "new_today": 3,
        "new_this_week": 15
      },
      "jobs": {
        "total": 1500,
        "active": 800,
        "closed": 400,
        "expired": 300,
        "new_today": 20,
        "new_this_week": 120
      },
      "applications": {
        "total": 12000,
        "today": 350,
        "this_week": 2200
      },
      "charts": {
        "user_registration": [...], // Time series data
        "job_posting": [...],
        "applications": [...]
      }
    }
  }
  ```

### Get Users (Admin)

- **URL**: GET `/api/v1/admin/users`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - role: string (optional) - "job_seeker", "employer", "admin"
  - status: string (optional) - "active", "inactive", "suspended"
  - search: string (optional) - search by name, email
  - sort: string (optional) - "newest", "oldest"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Users retrieved successfully",
    "data": {
      "users": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "name": "John Doe",
          "email": "john@example.com",
          "role": "job_seeker",
          "status": "active",
          "is_verified": true,
          "created_at": "2025-10-01T08:00:00Z",
          "last_login": "2025-10-15T10:30:00Z"
        }
        // more users...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 5000,
        "total_pages": 500
      }
    }
  }
  ```

### Get User Detail (Admin)

- **URL**: GET `/api/v1/admin/users/:id`
- **Auth**: Required (Admin only)
- **Params**: id (User ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "User details retrieved",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+628123456789",
      "role": "job_seeker",
      "status": "active",
      "is_verified": true,
      "created_at": "2025-10-01T08:00:00Z",
      "last_login": "2025-10-15T10:30:00Z",
      "profile_completion": 85,
      "activity": {
        "applications_count": 10,
        "companies_followed": 5,
        "documents_uploaded": 3
      },
      "login_history": [
        {
          "timestamp": "2025-10-15T10:30:00Z",
          "ip": "192.168.1.1",
          "user_agent": "Mozilla/5.0..."
        }
        // more login history...
      ]
    }
  }
  ```

### Update User Status (Admin)

- **URL**: PUT `/api/v1/admin/users/:id/status`
- **Auth**: Required (Admin only)
- **Params**: id (User ID)
- **Body**:
  ```json
  {
    "status": "suspended",
    "reason": "Violated terms of service",
    "duration": 7 // days, null for permanent
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "User status updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "suspended",
      "updated_at": "2025-10-16T12:00:00Z"
    }
  }
  ```

### Get Companies (Admin)

- **URL**: GET `/api/v1/admin/companies`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - status: string (optional) - "active", "inactive", "suspended"
  - verification: string (optional) - "verified", "pending", "rejected", "not_requested"
  - search: string (optional) - search by name
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Companies retrieved successfully",
    "data": {
      "companies": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440020",
          "name": "Tech Company",
          "status": "active",
          "is_verified": true,
          "owner_id": "550e8400-e29b-41d4-a716-446655440060",
          "owner_name": "Jane Smith",
          "active_job_count": 8,
          "created_at": "2025-08-01T08:00:00Z"
        }
        // more companies...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 300,
        "total_pages": 30
      }
    }
  }
  ```

### Get Company Detail (Admin)

- **URL**: GET `/api/v1/admin/companies/:id`
- **Auth**: Required (Admin only)
- **Params**: id (Company ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company details retrieved",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "name": "Tech Company",
      "description": "Leading tech company...",
      "website": "https://techcompany.com",
      "industry": "Software Development",
      "status": "active",
      "is_verified": true,
      "verification_history": [
        {
          "status": "verified",
          "timestamp": "2025-09-01T10:00:00Z",
          "admin_id": "550e8400-e29b-41d4-a716-446655440100",
          "admin_name": "Admin User",
          "note": "All documents verified"
        }
      ],
      "owner": {
        "id": "550e8400-e29b-41d4-a716-446655440060",
        "name": "Jane Smith",
        "email": "jane@techcompany.com"
      },
      "jobs": {
        "total": 25,
        "active": 8,
        "closed": 12,
        "expired": 5
      },
      "activity": {
        "job_applications_received": 350,
        "followers_count": 2300,
        "reviews_count": 45,
        "average_rating": 4.5
      },
      "verification_documents": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440090",
          "title": "Company Registration",
          "file_url": "https://storage.example.com/documents/company_registration.pdf",
          "uploaded_at": "2025-08-15T09:00:00Z"
        }
      ]
    }
  }
  ```

### Verify Company (Admin)

- **URL**: PUT `/api/v1/admin/companies/:id/verify`
- **Auth**: Required (Admin only)
- **Params**: id (Company ID)
- **Body**:
  ```json
  {
    "is_verified": true,
    "note": "All documents verified",
    "request_id": "550e8400-e29b-41d4-a716-446655440095" // Verification request ID
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company verification status updated",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "is_verified": true,
      "updated_at": "2025-10-16T12:10:00Z"
    }
  }
  ```

### Update Company Status (Admin)

- **URL**: PUT `/api/v1/admin/companies/:id/status`
- **Auth**: Required (Admin only)
- **Params**: id (Company ID)
- **Body**:
  ```json
  {
    "status": "suspended",
    "reason": "Violated terms of service",
    "duration": 30 // days, null for permanent
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Company status updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "status": "suspended",
      "updated_at": "2025-10-16T12:15:00Z"
    }
  }
  ```

### Get Jobs (Admin)

- **URL**: GET `/api/v1/admin/jobs`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - status: string (optional) - "active", "closed", "expired", "flagged"
  - company_id: string (optional)
  - search: string (optional) - search by title
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Jobs retrieved successfully",
    "data": {
      "jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440010",
          "title": "Software Engineer",
          "company": {
            "id": "550e8400-e29b-41d4-a716-446655440020",
            "name": "Tech Company"
          },
          "status": "active",
          "application_count": 12,
          "posted_at": "2025-10-10T08:00:00Z",
          "flags_count": 0
        }
        // more jobs...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 1500,
        "total_pages": 150
      }
    }
  }
  ```

### Update Job Status (Admin)

- **URL**: PUT `/api/v1/admin/jobs/:id/status`
- **Auth**: Required (Admin only)
- **Params**: id (Job ID)
- **Body**:
  ```json
  {
    "status": "flagged",
    "reason": "Job description contains inappropriate content",
    "notify_employer": true
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job status updated successfully",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "status": "flagged",
      "updated_at": "2025-10-16T12:20:00Z"
    }
  }
  ```

### Get Applications (Admin)

- **URL**: GET `/api/v1/admin/applications`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
  - status: string (optional) - "active", "completed", "withdrawn"
  - job_id: string (optional)
  - company_id: string (optional)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Applications retrieved successfully",
    "data": {
      "applications": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440030",
          "user": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "name": "John Doe"
          },
          "job": {
            "id": "550e8400-e29b-41d4-a716-446655440010",
            "title": "Software Engineer"
          },
          "company": {
            "id": "550e8400-e29b-41d4-a716-446655440020",
            "name": "Tech Company"
          },
          "stage": "screening",
          "status": "active",
          "applied_at": "2025-10-12T10:00:00Z"
        }
        // more applications...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 12000,
        "total_pages": 1200
      }
    }
  }
  ```

### Get User Report (Admin)

- **URL**: GET `/api/v1/admin/reports/users`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - period: string (optional) - "daily", "weekly", "monthly", "yearly" (default: "monthly")
  - start_date: string (optional) - YYYY-MM-DD
  - end_date: string (optional) - YYYY-MM-DD
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "User report retrieved",
    "data": {
      "summary": {
        "total_users": 5000,
        "active_users": 3800,
        "inactive_users": 1000,
        "suspended_users": 200,
        "job_seekers": 4500,
        "employers": 500,
        "verification_rate": 0.85
      },
      "registration_trend": [
        {
          "period": "2025-09",
          "count": 250,
          "job_seekers": 220,
          "employers": 30
        },
        {
          "period": "2025-10",
          "count": 280,
          "job_seekers": 245,
          "employers": 35
        }
        // more periods...
      ],
      "activity_metrics": {
        "daily_active_users": 1200,
        "weekly_active_users": 2500,
        "monthly_active_users": 3800,
        "average_session_duration": 15.3 // minutes
      },
      "demographics": {
        "locations": [
          { "name": "Jakarta", "count": 2200 },
          { "name": "Bandung", "count": 800 }
          // more locations...
        ],
        "industries": [
          { "name": "Technology", "count": 1500 },
          { "name": "Finance", "count": 800 }
          // more industries...
        ]
      }
    }
  }
  ```

### Get Job Report (Admin)

- **URL**: GET `/api/v1/admin/reports/jobs`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - period: string (optional) - "daily", "weekly", "monthly", "yearly" (default: "monthly")
  - start_date: string (optional) - YYYY-MM-DD
  - end_date: string (optional) - YYYY-MM-DD
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job report retrieved",
    "data": {
      "summary": {
        "total_jobs": 1500,
        "active_jobs": 800,
        "closed_jobs": 400,
        "expired_jobs": 300,
        "flagged_jobs": 15
      },
      "posting_trend": [
        {
          "period": "2025-09",
          "count": 180,
          "active": 100,
          "closed": 50,
          "expired": 30
        },
        {
          "period": "2025-10",
          "count": 210,
          "active": 120,
          "closed": 60,
          "expired": 30
        }
        // more periods...
      ],
      "job_metrics": {
        "average_active_duration": 25.3, // days
        "average_applications_per_job": 8.2,
        "high_performing_jobs": 120 // jobs with above-average applications
      },
      "categorization": {
        "job_types": [
          { "name": "full_time", "count": 1100 },
          { "name": "part_time", "count": 200 },
          { "name": "contract", "count": 150 },
          { "name": "internship", "count": 50 }
        ],
        "experience_levels": [
          { "name": "entry", "count": 400 },
          { "name": "mid", "count": 800 },
          { "name": "senior", "count": 250 },
          { "name": "executive", "count": 50 }
        ],
        "locations": [
          { "name": "Jakarta", "count": 700 },
          { "name": "Bandung", "count": 250 }
          // more locations...
        ],
        "top_skills": [
          { "name": "JavaScript", "count": 300 },
          { "name": "React", "count": 250 }
          // more skills...
        ]
      }
    }
  }
  ```

### Get Application Report (Admin)

- **URL**: GET `/api/v1/admin/reports/applications`
- **Auth**: Required (Admin only)
- **Query Parameters**:
  - period: string (optional) - "daily", "weekly", "monthly", "yearly" (default: "monthly")
  - start_date: string (optional) - YYYY-MM-DD
  - end_date: string (optional) - YYYY-MM-DD
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Application report retrieved",
    "data": {
      "summary": {
        "total_applications": 12000,
        "active_applications": 5000,
        "completed_applications": 6000,
        "withdrawn_applications": 1000
      },
      "application_trend": [
        {
          "period": "2025-09",
          "count": 2800,
          "active": 1200,
          "completed": 1400,
          "withdrawn": 200
        },
        {
          "period": "2025-10",
          "count": 3200,
          "active": 1500,
          "completed": 1500,
          "withdrawn": 200
        }
        // more periods...
      ],
      "stage_metrics": {
        "applied": 2000,
        "screening": 1500,
        "interview": 1000,
        "offer": 500,
        "rejected": 1000,
        "average_time_to_complete": 14.5 // days
      },
      "conversion_rates": {
        "applied_to_screening": 0.75,
        "screening_to_interview": 0.67,
        "interview_to_offer": 0.5,
        "offer_to_hire": 0.8
      }
    }
  }
  ```

### Update Master Data (Admin)

- **URL**: PUT `/api/v1/admin/master/:type/:id`
- **Auth**: Required (Admin only)
- **Params**:
  - type: "skills", "benefits", "industries", "job_types"
  - id: Item ID
- **Body**:
  ```json
  {
    "name": "Updated Skill Name",
    "description": "Updated description",
    "category": "Updated Category",
    "icon": "updated-icon",
    "status": "active"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Master data updated successfully",
    "data": {
      "id": 1,
      "name": "Updated Skill Name",
      "description": "Updated description",
      "updated_at": "2025-10-16T14:20:00Z"
    }
  }
  ```

### Create Master Data (Admin)

- **URL**: POST `/api/v1/admin/master/:type`
- **Auth**: Required (Admin only)
- **Params**:
  - type: "skills", "benefits", "industries", "job_types"
- **Body**:
  ```json
  {
    "name": "New Skill",
    "description": "New skill description",
    "category": "Programming",
    "icon": "code",
    "status": "active"
  }
  ```
- **Response Success** (201):
  ```json
  {
    "success": true,
    "message": "Master data created successfully",
    "data": {
      "id": 10,
      "name": "New Skill",
      "description": "New skill description",
      "created_at": "2025-10-16T14:25:00Z"
    }
  }
  ```

---

## 7. Master Data Endpoints

### Get Skills

- **URL**: GET `/api/v1/master/skills`
- **Auth**: Tidak perlu
- **Query Parameters**:
  - category: string (optional)
  - search: string (optional)
  - page: number (default: 1)
  - limit: number (default: 50)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Skills retrieved successfully",
    "data": {
      "skills": [
        {
          "id": 1,
          "name": "JavaScript",
          "category": "Programming",
          "popularity": 95
        },
        {
          "id": 2,
          "name": "React",
          "category": "Framework",
          "popularity": 90
        }
        // more skills...
      ],
      "pagination": {
        "page": 1,
        "limit": 50,
        "total_items": 200,
        "total_pages": 4
      }
    }
  }
  ```

### Get Benefits

- **URL**: GET `/api/v1/master/benefits`
- **Auth**: Tidak perlu
- **Query Parameters**:
  - category: string (optional)
  - search: string (optional)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Benefits retrieved successfully",
    "data": {
      "benefits": [
        {
          "id": 1,
          "name": "Health Insurance",
          "category": "Healthcare",
          "icon": "health"
        },
        {
          "id": 2,
          "name": "Remote Work",
          "category": "Work Arrangement",
          "icon": "remote"
        }
        // more benefits...
      ]
    }
  }
  ```

### Get Industries

- **URL**: GET `/api/v1/master/industries`
- **Auth**: Tidak perlu
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Industries retrieved successfully",
    "data": {
      "industries": [
        {
          "id": 1,
          "name": "Software Development"
        },
        {
          "id": 2,
          "name": "Financial Services"
        }
        // more industries...
      ]
    }
  }
  ```

### Get Job Types

- **URL**: GET `/api/v1/master/job-types`
- **Auth**: Tidak perlu
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job types retrieved successfully",
    "data": {
      "job_types": [
        {
          "id": 1,
          "name": "full_time",
          "display_name": "Full Time"
        },
        {
          "id": 2,
          "name": "part_time",
          "display_name": "Part Time"
        },
        {
          "id": 3,
          "name": "contract",
          "display_name": "Contract"
        },
        {
          "id": 4,
          "name": "internship",
          "display_name": "Internship"
        }
      ]
    }
  }
  ```

### Get Experience Levels

- **URL**: GET `/api/v1/master/experience-levels`
- **Auth**: Tidak perlu
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Experience levels retrieved successfully",
    "data": {
      "experience_levels": [
        {
          "id": 1,
          "name": "entry",
          "display_name": "Entry Level"
        },
        {
          "id": 2,
          "name": "mid",
          "display_name": "Mid Level"
        },
        {
          "id": 3,
          "name": "senior",
          "display_name": "Senior Level"
        },
        {
          "id": 4,
          "name": "executive",
          "display_name": "Executive Level"
        }
      ]
    }
  }
  ```

---

## 8. Notification Endpoints

### Get Notifications

- **URL**: GET `/api/v1/notifications`
- **Auth**: Required
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 20)
  - is_read: boolean (optional)
  - type: string (optional) - "application", "job", "company", "admin"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Notifications retrieved successfully",
    "data": {
      "notifications": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440200",
          "type": "application",
          "title": "Interview Scheduled",
          "message": "Your interview for Software Engineer position has been scheduled",
          "data": {
            "application_id": "550e8400-e29b-41d4-a716-446655440030",
            "job_id": "550e8400-e29b-41d4-a716-446655440010",
            "company_id": "550e8400-e29b-41d4-a716-446655440020",
            "interview_id": "550e8400-e29b-41d4-a716-446655440040"
          },
          "is_read": false,
          "created_at": "2025-10-15T10:00:00Z"
        }
        // more notifications...
      ],
      "pagination": {
        "page": 1,
        "limit": 20,
        "total_items": 45,
        "total_pages": 3
      },
      "unread_count": 8
    }
  }
  ```

### Mark Notification as Read

- **URL**: PUT `/api/v1/notifications/:id/read`
- **Auth**: Required
- **Params**: id (Notification ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Notification marked as read",
    "data": {
      "id": "550e8400-e29b-41d4-a716-446655440200",
      "is_read": true,
      "updated_at": "2025-10-16T14:30:00Z"
    }
  }
  ```

### Mark All Notifications as Read

- **URL**: PUT `/api/v1/notifications/read-all`
- **Auth**: Required
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "All notifications marked as read",
    "data": {
      "updated_count": 8,
      "unread_count": 0
    }
  }
  ```

### Update Notification Settings

- **URL**: PUT `/api/v1/notifications/settings`
- **Auth**: Required
- **Body**:
  ```json
  {
    "email_notifications": {
      "application_updates": true,
      "job_recommendations": true,
      "company_updates": false,
      "marketing": false
    },
    "push_notifications": {
      "application_updates": true,
      "job_recommendations": false,
      "company_updates": true,
      "marketing": false
    }
  }
  ```
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Notification settings updated successfully",
    "data": {
      "updated_at": "2025-10-16T14:35:00Z"
    }
  }
  ```

---

## 9. Search & Recommendations

### Get Job Recommendations

- **URL**: GET `/api/v1/recommendations/jobs`
- **Auth**: Required
- **Query Parameters**:
  - page: number (default: 1)
  - limit: number (default: 10)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job recommendations retrieved",
    "data": {
      "jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440012",
          "title": "Frontend Developer",
          "company": {
            "id": "550e8400-e29b-41d4-a716-446655440020",
            "name": "Tech Company",
            "logo_url": "https://storage.example.com/logos/company.jpg"
          },
          "location": "Jakarta, Indonesia",
          "job_type": "full_time",
          "posted_at": "2025-10-14T08:00:00Z",
          "match_score": 95,
          "match_reasons": ["skills", "location", "experience"]
        }
        // more jobs...
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total_items": 30,
        "total_pages": 3
      }
    }
  }
  ```

### Get Similar Jobs

- **URL**: GET `/api/v1/jobs/:id/similar`
- **Auth**: Tidak perlu
- **Params**: id (Job ID)
- **Query Parameters**:
  - limit: number (default: 5)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Similar jobs retrieved",
    "data": {
      "jobs": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440013",
          "title": "React Developer",
          "company": {
            "id": "550e8400-e29b-41d4-a716-446655440021",
            "name": "Another Tech Company",
            "logo_url": "https://storage.example.com/logos/company2.jpg"
          },
          "location": "Jakarta, Indonesia",
          "job_type": "full_time",
          "posted_at": "2025-10-15T08:00:00Z",
          "similarity_score": 90
        }
        // more jobs...
      ]
    }
  }
  ```

### Get Job Statistics

- **URL**: GET `/api/v1/jobs/:id/statistics`
- **Auth**: Required (Job Owner/Employer only)
- **Params**: id (Job ID)
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Job statistics retrieved",
    "data": {
      "views": 1250,
      "unique_views": 950,
      "applications": 25,
      "application_rate": 2.63, // percentage
      "view_trend": [
        { "date": "2025-10-10", "views": 150, "unique_views": 120 },
        { "date": "2025-10-11", "views": 200, "unique_views": 150 }
        // more dates...
      ],
      "application_trend": [
        { "date": "2025-10-10", "applications": 3 },
        { "date": "2025-10-11", "applications": 5 }
        // more dates...
      ],
      "candidate_demographics": {
        "experience_levels": [
          { "name": "entry", "count": 5 },
          { "name": "mid", "count": 15 },
          { "name": "senior", "count": 5 }
        ],
        "top_skills": [
          { "name": "JavaScript", "count": 20 },
          { "name": "React", "count": 18 }
          // more skills...
        ]
      }
    }
  }
  ```

---

## 10. Search Suggestion Endpoints

### Get Search Suggestions

- **URL**: GET `/api/v1/search/suggestions`
- **Auth**: Tidak perlu
- **Query Parameters**:
  - q: string - search query
  - type: string (optional) - "job", "company", "skill", "location"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Search suggestions retrieved",
    "data": {
      "suggestions": [
        {
          "type": "job_title",
          "text": "Software Engineer",
          "job_count": 120
        },
        {
          "type": "company",
          "text": "Tech Company",
          "company_id": "550e8400-e29b-41d4-a716-446655440020"
        },
        {
          "type": "skill",
          "text": "JavaScript",
          "skill_id": 1
        },
        {
          "type": "location",
          "text": "Jakarta, Indonesia"
        }
      ]
    }
  }
  ```

### Get Popular Searches

- **URL**: GET `/api/v1/search/popular`
- **Auth**: Tidak perlu
- **Query Parameters**:
  - type: string (optional) - "job", "company", "skill", "location"
- **Response Success** (200):
  ```json
  {
    "success": true,
    "message": "Popular searches retrieved",
    "data": {
      "popular": [
        {
          "type": "job_title",
          "text": "Software Engineer",
          "search_count": 5000
        },
        {
          "type": "job_title",
          "text": "Data Scientist",
          "search_count": 4500
        },
        {
          "type": "company",
          "text": "Tech Company",
          "search_count": 3000
        }
        // more popular searches...
      ]
    }
  }
  ```

---

## Kode Status HTTP dan Pesan Error

### Kode Status HTTP yang Digunakan:

- **200** OK: Request berhasil
- **201** Created: Resource berhasil dibuat
- **400** Bad Request: Request tidak valid
- **401** Unauthorized: Token tidak ada atau tidak valid
- **403** Forbidden: Tidak memiliki izin untuk mengakses resource
- **404** Not Found: Resource tidak ditemukan
- **409** Conflict: Resource sudah ada atau konflik
- **422** Unprocessable Entity: Validasi gagal
- **429** Too Many Requests: Rate limit terlampaui
- **500** Internal Server Error: Terjadi kesalahan di server

### Contoh Pesan Error:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Must be a valid email address"
    },
    {
      "field": "password",
      "message": "Password must be at least 8 characters and contain at least one number, one uppercase letter, and one special character"
    }
  ]
}
```

---

## Rate Limiting

API ini menerapkan rate limiting untuk mencegah penyalahgunaan. Header respons akan menyertakan:

- `X-RateLimit-Limit`: Jumlah request maksimum per periode
- `X-RateLimit-Remaining`: Sisa request yang tersedia
- `X-RateLimit-Reset`: Waktu (dalam detik) hingga batas rate direset

Jika rate limit terlampaui, server akan mengembalikan status 429 Too Many Requests.

---

## Versioning

Semua endpoint diawali dengan `/api/v1`. Jika ada perubahan yang tidak backward compatible, versi akan dinaikkan menjadi `/api/v2`, dan seterusnya.
