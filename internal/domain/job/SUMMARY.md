# Job Domain

## Overview

Job domain adalah core business domain untuk Keerja job portal yang mengelola job postings, categories, locations, benefits, skills, dan requirements. Domain ini mencakup 7 entities dengan comprehensive business logic untuk job management, search, matching, dan analytics.

---

## Entities (7)

### 1. **Job** (Main Entity)

Core entity untuk job posting dengan 38+ fields:

**Key Fields:**

- ID, UUID, CompanyID, EmployerUserID, CategoryID
- Title, Slug, JobLevel, EmploymentType
- Description, Requirements, Responsibilities
- Location fields (Location, City, Province, RemoteOption)
- Salary range (SalaryMin, SalaryMax, Currency)
- Experience range (ExperienceMin, ExperienceMax)
- EducationLevel, TotalHires
- Status (draft, published, closed, expired, suspended)
- ViewsCount, ApplicationsCount
- PublishedAt, ExpiredAt, CreatedAt, UpdatedAt

**Relationships:**

- BelongsTo: JobCategory
- HasMany: JobLocation, JobBenefit, JobSkill, JobRequirement

**Helper Methods:**

- `IsPublished()` - Check if job is published
- `IsClosed()` - Check if job is closed/expired
- `IsExpired()` - Check if job has expired
- `IsActive()` - Check if job is active and accepting applications
- `CanApply()` - Check if job accepts applications

**Enums:**

- JobLevel: Internship, Entry Level, Mid Level, Senior Level, Manager, Director
- EmploymentType: Full-Time, Part-Time, Contract, Internship, Freelance
- Status: draft, published, closed, expired, suspended

---

### 2. **JobCategory**

Hierarchical category system untuk job classification:

**Key Fields:**

- ID, ParentID (for hierarchy), Code, Name
- Description, IsActive
- CreatedAt, UpdatedAt

**Relationships:**

- Self-referential: Parent → Children (hierarchical)
- HasMany: JobSubcategory, Job

---

### 3. **JobSubcategory**

Subcategory untuk granular job classification:

**Key Fields:**

- ID, CategoryID, Code, Name
- Description, IsActive
- CreatedAt, UpdatedAt

**Relationships:**

- BelongsTo: JobCategory

---

### 4. **JobLocation**

Multiple locations untuk satu job (onsite, hybrid, remote):

**Key Fields:**

- ID, JobID, CompanyID, LocationType
- Address, City, Province, PostalCode, Country
- Latitude, Longitude (for geolocation search)
- GooglePlaceID, MapURL
- IsPrimary
- CreatedAt, UpdatedAt

**Enums:**

- LocationType: onsite, hybrid, remote

**Helper Methods:**

- `IsRemote()` - Check if location is remote
- `IsHybrid()` - Check if location is hybrid

**Features:**

- GIS support dengan GIST index untuk geolocation search
- Primary location designation
- Google Maps integration

---

### 5. **JobBenefit**

Benefits offered dengan job:

**Key Fields:**

- ID, JobID, BenefitID (reference to master), BenefitName
- Description, IsHighlight
- CreatedAt, UpdatedAt

**Features:**

- Can reference benefits_master atau custom benefit
- Highlight important benefits
- Unique constraint pada (JobID, BenefitName)

---

### 6. **JobSkill**

Skills required untuk job (many-to-many dengan skills_master):

**Key Fields:**

- ID, JobID, SkillID
- ImportanceLevel (required, preferred, optional)
- Weight (0.00 - 1.00 for matching algorithm)
- CreatedAt, UpdatedAt

**Helper Methods:**

- `IsRequired()` - Check if skill is required
- `IsPreferred()` - Check if skill is preferred

**Features:**

- Weighted skills untuk matching algorithm
- Importance levels untuk filtering
- Unique constraint pada (JobID, SkillID)

---

### 7. **JobRequirement**

Detailed requirements dengan different types:

**Key Fields:**

- ID, JobID, RequirementType
- RequirementText (detailed text)
- SkillID (optional reference)
- MinExperience, MaxExperience
- EducationLevel, Language
- IsMandatory, Priority
- CreatedAt, UpdatedAt

**Enums:**

- RequirementType: education, experience, skill, language, certification, other

**Helper Methods:**

- `IsEducationRequirement()` - Check if requirement is education type
- `IsExperienceRequirement()` - Check if requirement is experience type
- `IsSkillRequirement()` - Check if requirement is skill type

---

## Repository Interface (90+ methods)

### Job CRUD (7 methods)

- `Create()`, `FindByID()`, `FindByUUID()`, `FindBySlug()`
- `Update()`, `Delete()`, `SoftDelete()`

### Job Listing & Search (7 methods)

- `List()` - List jobs dengan filter
- `ListByCompany()` - Company's jobs
- `ListByEmployer()` - Employer user's jobs
- `SearchJobs()` - Advanced search dengan JobSearchFilter
- `SearchByLocation()` - Geolocation-based search
- `SearchBySkills()` - Search by skill IDs
- `SearchBySalaryRange()` - Search by salary range

### Job Status Operations (7 methods)

- `UpdateStatus()`, `PublishJob()`, `CloseJob()`, `ExpireJob()`, `SuspendJob()`
- `GetExpiredJobs()` - Jobs yang sudah expired
- `GetExpiringJobs()` - Jobs yang akan expire dalam X days

### Job Statistics (4 methods)

- `IncrementViews()` - Track job views
- `IncrementApplications()` - Track applications
- `GetJobStats()` - Individual job stats
- `GetCompanyJobStats()` - Company's overall job stats

### Recommendation & Matching (3 methods)

- `GetRecommendedJobs()` - Personalized recommendations for user
- `GetSimilarJobs()` - Similar jobs berdasarkan job ID
- `GetMatchingJobs()` - Jobs matching user profile

### JobCategory CRUD (7 methods)

- `CreateCategory()`, `FindCategoryByID()`, `FindCategoryByCode()`
- `UpdateCategory()`, `DeleteCategory()`
- `ListCategories()`, `GetCategoryTree()`, `GetActiveCategories()`

### JobSubcategory CRUD (6 methods)

- `CreateSubcategory()`, `FindSubcategoryByID()`, `FindSubcategoryByCode()`
- `UpdateSubcategory()`, `DeleteSubcategory()`
- `ListSubcategories()`, `GetActiveSubcategories()`

### JobLocation Operations (7 methods)

- `CreateLocation()`, `FindLocationByID()`, `UpdateLocation()`, `DeleteLocation()`
- `ListLocationsByJob()`, `GetPrimaryLocation()`, `SetPrimaryLocation()`

### JobBenefit Operations (8 methods)

- `CreateBenefit()`, `FindBenefitByID()`, `UpdateBenefit()`, `DeleteBenefit()`
- `ListBenefitsByJob()`, `GetHighlightedBenefits()`
- `BulkCreateBenefits()`, `BulkDeleteBenefits()`

### JobSkill Operations (8 methods)

- `CreateSkill()`, `FindSkillByID()`, `UpdateSkill()`, `DeleteSkill()`
- `ListSkillsByJob()`, `GetRequiredSkills()`, `GetPreferredSkills()`
- `BulkCreateSkills()`, `BulkDeleteSkills()`

### JobRequirement Operations (7 methods)

- `CreateRequirement()`, `FindRequirementByID()`, `UpdateRequirement()`, `DeleteRequirement()`
- `ListRequirementsByJob()`, `GetMandatoryRequirements()`
- `BulkCreateRequirements()`, `BulkDeleteRequirements()`

### Analytics (3 methods)

- `GetTrendingJobs()` - Most viewed/applied jobs
- `GetPopularCategories()` - Category stats
- `GetJobsByDateRange()` - Jobs in date range

---

## Service Interface (80+ methods)

### Job Management - Employer (8 methods)

- `CreateJob()` - Create new job dengan nested data (locations, benefits, skills, requirements)
- `UpdateJob()` - Update job details
- `DeleteJob()` - Soft delete job
- `GetJob()`, `GetJobBySlug()`, `GetJobByUUID()` - Retrieve job
- `GetMyJobs()` - Employer's jobs
- `GetCompanyJobs()` - Company's all jobs

### Job Status Management (9 methods)

- `PublishJob()` - Publish draft job
- `UnpublishJob()` - Unpublish to draft
- `CloseJob()` - Close job (no more applications)
- `ReopenJob()` - Reopen closed job
- `SuspendJob()` - Suspend job dengan reason
- `SetJobExpiry()` - Set expiry date
- `ExtendJobExpiry()` - Extend by X days
- `AutoExpireJobs()` - Batch expire expired jobs (cron job)

### Job Search & Discovery - Public (9 methods)

- `ListJobs()` - List jobs dengan filter
- `SearchJobs()` - Advanced search dengan facets
- `SearchJobsByLocation()` - Geolocation search dengan radius
- `GetFeaturedJobs()` - Featured/promoted jobs
- `GetLatestJobs()` - Recently posted jobs
- `GetTrendingJobs()` - Popular jobs by views
- `GetRecommendedJobs()` - Personalized untuk user
- `GetSimilarJobs()` - Similar jobs

### Job Matching (2 methods)

- `CalculateMatchScore()` - Calculate match score between job & user
- `GetMatchingJobs()` - Get matching jobs dengan score

### Job Views & Interactions (3 methods)

- `IncrementView()` - Track job view (dengan user tracking)
- `GetJobStats()` - Individual job statistics
- `GetCompanyJobStats()` - Company's job statistics

### Job Details Management (16 methods)

**Location:**

- `AddLocation()`, `UpdateLocation()`, `DeleteLocation()`, `SetPrimaryLocation()`

**Benefits:**

- `AddBenefit()`, `UpdateBenefit()`, `DeleteBenefit()`, `BulkAddBenefits()`

**Skills:**

- `AddSkill()`, `UpdateSkill()`, `DeleteSkill()`, `BulkAddSkills()`

**Requirements:**

- `AddRequirement()`, `UpdateRequirement()`, `DeleteRequirement()`, `BulkAddRequirements()`

### Category Management - Admin (8 methods)

- `CreateCategory()`, `UpdateCategory()`, `DeleteCategory()`
- `GetCategory()`, `GetCategoryByCode()`
- `ListCategories()`, `GetCategoryTree()`, `GetActiveCategories()`

### Subcategory Management - Admin (6 methods)

- `CreateSubcategory()`, `UpdateSubcategory()`, `DeleteSubcategory()`
- `GetSubcategory()`, `ListSubcategories()`, `GetActiveSubcategories()`

### Analytics & Reporting (5 methods)

- `GetJobAnalytics()` - Time-series analytics for job
- `GetCompanyAnalytics()` - Company analytics dengan breakdown
- `GetCategoryAnalytics()` - Category analytics
- `GetPopularCategories()` - Popular categories list
- `GetTopCompanies()` - Top companies by activity

### Bulk Operations (3 methods)

- `BulkPublishJobs()`, `BulkCloseJobs()`, `BulkDeleteJobs()`

### Validation (3 methods)

- `ValidateJob()` - Validate job data
- `CheckJobOwnership()` - Verify employer owns job
- `CheckJobStatus()` - Get current job status

---

## Request DTOs (13)

1. **CreateJobRequest** - Create job dengan 20+ fields + nested data
2. **UpdateJobRequest** - Update job fields
3. **AddLocationRequest** - Add job location dengan geocoding
4. **UpdateLocationRequest** - Update location details
5. **AddBenefitRequest** - Add benefit
6. **UpdateBenefitRequest** - Update benefit
7. **AddSkillRequest** - Add skill dengan importance level
8. **UpdateSkillRequest** - Update skill importance
9. **AddRequirementRequest** - Add requirement
10. **UpdateRequirementRequest** - Update requirement
11. **CreateCategoryRequest** - Create category
12. **UpdateCategoryRequest** - Update category
13. **CreateSubcategoryRequest** - Create subcategory
14. **UpdateSubcategoryRequest** - Update subcategory

---

## Response DTOs (13)

1. **JobSearchResponse** - Search results dengan facets & suggestions
2. **SearchFacets** - Faceted search filters (categories, locations, levels, types, salaries)
3. **FacetItem** - Individual facet dengan count
4. **MatchScore** - Job-user match score dengan breakdown (skill, experience, education, location)
5. **MatchResponse** - Matching jobs dengan scores
6. **JobWithScore** - Job entity dengan match score
7. **JobAnalytics** - Time-series analytics data
8. **CompanyAnalytics** - Company job analytics
9. **CategoryAnalytics** - Category analytics
10. **TimeSeriesData** - Time-series data point
11. **SourceStats** - Traffic source statistics
12. **JobPerformance** - Job performance metrics
13. **CompanyStats** - Company statistics

---

## Filters & Search Types (3)

1. **JobFilter** - Basic filtering (status, company, category, location, level, type, salary, experience, education)
2. **JobSearchFilter** - Advanced search (keyword, location, categories, skills, levels, types, remote, salary, experience, education, companies, posted within)
3. **CategoryFilter** - Category filtering (parent, active, keyword)

---

## Statistics Types (3)

1. **JobStats** - Job statistics (views, applications breakdown by period, conversion rate)
2. **CompanyJobStats** - Company statistics (total jobs, status breakdown, views/applications totals & averages)
3. **CategoryStats** - Category statistics (job count, views, applications)

---

## Business Features

### 1. **Job Posting Workflow**

- Draft → Published → Closed/Expired
- Auto-expiration handling
- Bulk operations
- Status management

### 2. **Advanced Search**

- Full-text search dengan keyword
- Multi-criteria filtering
- Faceted search dengan counts
- Search suggestions
- Geolocation search dengan radius
- Skills-based search
- Salary range search

### 3. **Job Matching & Recommendation**

- Calculate match score (skill, experience, education, location)
- Personalized job recommendations
- Similar jobs discovery
- Weighted skill matching

### 4. **Job Analytics**

- Views tracking (unique & total)
- Applications tracking
- Conversion rate calculation
- Time-series data
- Traffic source analysis
- Performance metrics
- Company-level analytics
- Category-level analytics

### 5. **Multi-Location Support**

- Multiple locations per job
- Location types (onsite, hybrid, remote)
- Primary location designation
- Geocoding support
- Google Maps integration
- Radius-based search

### 6. **Skills Management**

- Weighted skills (0.00 - 1.00)
- Importance levels (required, preferred, optional)
- Bulk operations
- Skills matching algorithm

### 7. **Flexible Requirements**

- Multiple requirement types
- Mandatory vs optional
- Priority ordering
- Structured data (experience range, education level)

### 8. **Benefits Highlighting**

- Custom or master-based benefits
- Highlighted benefits
- Bulk operations

---

## Technical Features

1. **GORM Integration**

   - Proper relationships dengan foreignKey & constraints
   - CASCADE delete untuk child entities
   - SET NULL untuk optional relationships
   - Indexes untuk performance
   - GIST index untuk geolocation

2. **UUID Support**

   - UUID generation untuk external references
   - Slug generation untuk SEO-friendly URLs

3. **Validation**

   - Comprehensive validation tags
   - Enum validation
   - Business rule validation

4. **Soft Delete**

   - Support soft delete via GORM

5. **Timestamps**

   - Auto-managed CreatedAt & UpdatedAt

6. **Pagination**

   - Consistent pagination support
   - Total count tracking

7. **Filtering**

   - Flexible filter structs
   - Multiple filter types

8. **Bulk Operations**
   - Batch create/delete untuk child entities
   - Bulk status updates

---

## Statistics

- **Total Entities:** 7
- **Total Repository Methods:** ~90
- **Total Service Methods:** ~80
- **Total Request DTOs:** 13
- **Total Response DTOs:** 13
- **Total Filter Types:** 3
- **Total Stats Types:** 3
- **Total Lines of Code:** ~800 (entities + repository + service)

---

## Integration Points

### Depends On:

- Company domain (company_id reference)
- User domain (employer_user_id reference)
- Master domain (skills_master, benefits_master)

### Used By:

- Application domain (job applications)
- Search service (indexing)
- Notification service (job alerts)
- Analytics service (reporting)

---

## Next Steps

After Job domain completion:

1. **Application Domain** - Job applications, stages, interviews
2. **Master Domain** - Skills master, Benefits master
3. **Admin Domain** - Admin users & roles
4. **Repository Implementation** - Implement all repository interfaces
5. **Service Implementation** - Implement all business logic

---