# Company Domain 

### 1. Entity (entity.go)

**File:** `internal/domain/company/entity.go`

**8 Entities Created:**

1. **Company** - Main company entity

   - Basic info (name, slug, legal name, registration)
   - Location (address, city, province, coordinates)
   - Media (logo, banner)
   - Verification status
   - Relationships to all related entities

2. **CompanyProfile** - Detailed company profile

   - Marketing content (tagline, descriptions, mission, vision)
   - Gallery and media (images, video)
   - SEO optimization fields
   - Social media links (JSONB)
   - Publication status

3. **CompanyIndustry** - Industry classifications

   - Hierarchical structure (parent-child)
   - Code and name
   - Active status

4. **CompanyFollower** - User follows company

   - User-Company relationship
   - Follow/unfollow tracking
   - Active status

5. **CompanyReview** - Employee/ex-employee reviews

   - Multiple rating dimensions (culture, work-life, salary, management)
   - Pros, cons, advice
   - Anonymous option
   - Moderation workflow

6. **CompanyDocument** - Legal documents

   - Document types (SIUP, NPWP, NIB, AKTA, TDP, ISO, etc.)
   - Document verification workflow
   - Expiry tracking

7. **CompanyEmployee** - Employee records

   - Employment details (type, status, dates)
   - Salary range
   - Visibility controls
   - Verification status

8. **CompanyVerification** - Company verification status

   - Verification workflow
   - Score and notes
   - Badge system
   - Expiry tracking

9. **EmployerUser** - Users with employer privileges
   - Role-based access (owner, admin, recruiter, viewer)
   - Company-specific credentials
   - Permission methods

**Features:**

- All GORM tags configured
- Validation tags included
- JSON tags with omitempty
- Proper relationships (ForeignKey, OnDelete)
- Check constraints for enums
- Helper methods (IsVerified, IsOwner, etc.)
- PostgreSQL array and JSONB support

---

### 2. Repository Interface (repository.go)

**File:** `internal/domain/company/repository.go`

**70+ Methods Defined:**

**Company Operations (8 methods):**

- Create, FindByID, FindByUUID, FindBySlug
- Update, Delete, List, SearchCompanies

**Profile Operations (3 methods):**

- CreateProfile, FindProfileByCompanyID, UpdateProfile

**Follower Operations (6 methods):**

- FollowCompany, UnfollowCompany, IsFollowing
- GetFollowers, GetFollowedCompanies, CountFollowers

**Review Operations (9 methods):**

- Create, Update, Delete, FindByID
- GetReviewsByCompanyID, GetReviewsByUserID
- ApproveReview, RejectReview
- CalculateAverageRatings

**Document Operations (7 methods):**

- Create, Update, Delete, FindByID
- GetDocumentsByCompanyID
- ApproveDocument, RejectDocument

**Employee Operations (5 methods):**

- Add, Update, Delete
- GetEmployeesByCompanyID, CountEmployees

**Employer User Operations (7 methods):**

- Create, Update, Delete, FindByID
- FindByUserAndCompany
- GetEmployerUsersByCompanyID, GetCompaniesByUserID

**Verification Operations (7 methods):**

- Create, Update, FindByCompanyID
- RequestVerification, ApproveVerification, RejectVerification
- GetPendingVerifications

**Industry Operations (7 methods):**

- Create, Update, Delete
- FindByID, FindByCode
- GetAllIndustries, GetIndustryTree

**Analytics (4 methods):**

- SearchCompanies, GetVerifiedCompanies
- GetTopRatedCompanies
- GetCompaniesNeedingVerificationRenewal

**Supporting Types:**

- CompanyFilter
- ReviewFilter
- AverageRatings

---

### 3. Service Interface (service.go)

**File:** `internal/domain/company/service.go`

**80+ Methods Defined:**

**Company Management (6 methods):**

- RegisterCompany, GetCompany, GetCompanyBySlug
- UpdateCompany, DeleteCompany
- ListCompanies, SearchCompanies

**Profile Management (5 methods):**

- CreateProfile, UpdateProfile, GetProfile
- PublishProfile, UnpublishProfile

**Media Management (4 methods):**

- UploadLogo, UploadBanner
- DeleteLogo, DeleteBanner

**Follower Management (6 methods):**

- FollowCompany, UnfollowCompany, IsFollowing
- GetFollowers, GetFollowedCompanies, GetFollowerCount

**Review Management (9 methods):**

- AddReview, UpdateReview, DeleteReview, GetReview
- GetCompanyReviews, GetUserReviews, GetAverageRatings
- ApproveReview, RejectReview, HideReview
- GetPendingReviews (admin)

**Document Management (7 methods):**

- UploadDocument, UpdateDocument, DeleteDocument, GetDocuments
- ApproveDocument, RejectDocument (admin)
- CheckExpiredDocuments

**Employee Management (5 methods):**

- AddEmployee, UpdateEmployee, RemoveEmployee
- GetEmployees, GetEmployeeCount

**Employer User Management (7 methods):**

- InviteEmployer, AcceptInvitation
- UpdateEmployerRole, RemoveEmployerUser
- GetEmployerUsers, GetUserCompanies
- CheckEmployerPermission

**Verification Management (7 methods):**

- RequestVerification, GetVerificationStatus
- ApproveVerification, RejectVerification
- GetPendingVerifications, RenewVerification
- CheckVerificationExpiry

**Industry Management (6 methods):**

- CreateIndustry, UpdateIndustry, DeleteIndustry
- GetIndustry, GetAllIndustries, GetIndustryTree

**Analytics (4 methods):**

- GetCompanyStats, GetTopRatedCompanies
- GetVerifiedCompanies, GetCompanyEngagement

**Request DTOs (11 types):**

- RegisterCompanyRequest
- UpdateCompanyRequest
- CreateProfileRequest
- UpdateProfileRequest
- AddReviewRequest
- UpdateReviewRequest
- UploadDocumentRequest
- UpdateDocumentRequest
- AddEmployeeRequest
- UpdateEmployeeRequest
- InviteEmployerRequest
- CreateIndustryRequest
- UpdateIndustryRequest

**Response DTOs (2 types):**

- CompanyStats
- EngagementStats


## Key Features

### Business Logic Covered:

Company registration & management
Profile creation with SEO
Company following system
Employee review system with moderation
Document verification workflow
Employee management
Multi-user employer access with roles
Verification badge system
Industry hierarchy
Analytics and stats

---

## Next Steps

Lanjutkan ke domain berikutnya:

1. **Job Domain** - Job postings, categories, requirements
2. **Application Domain** - Job applications, stages, interviews
3. **Admin Domain** - Admin users and roles
4. **Master Domain** - Skills and benefits master data

---

