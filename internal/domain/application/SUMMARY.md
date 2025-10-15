# Application Domain

## Overview

Application domain mengelola proses hiring dari application submission sampai final decision (hired/rejected). Domain ini mencakup 5 entities dengan comprehensive business logic untuk application tracking, stage management, document handling, notes, dan interview scheduling.

---

## Entities (5)

### 1. **JobApplication** (Main Entity)

Core entity untuk job application dengan 14 fields:

**Key Fields:**

- ID, JobID, UserID, CompanyID
- AppliedAt, Status, Source
- MatchScore (calculated from job-user matching)
- NotesText (internal notes)
- ViewedByEmployer, IsBookmarked (employer tracking)
- ResumeURL
- CreatedAt, UpdatedAt

**Relationships:**

- HasMany: JobApplicationStage, ApplicationDocument, ApplicationNote, Interview

**Helper Methods:**

- `IsApplied()` - Check if status is applied
- `IsInProgress()` - Check if in hiring process (screening/shortlisted/interview/offered)
- `IsCompleted()` - Check if has final status (hired/rejected/withdrawn)
- `IsHired()`, `IsRejected()`, `IsWithdrawn()` - Check specific statuses
- `CanWithdraw()` - Check if user can withdraw

**Enums:**

- Status: applied, screening, shortlisted, interview, offered, hired, rejected, withdrawn (8 stages)

**Constraints:**

- Unique constraint on (JobID, UserID) - one application per job per user

---

### 2. **JobApplicationStage**

Stage tracking untuk hiring workflow:

**Key Fields:**

- ID, ApplicationID, StageName
- Description, HandledBy (recruiter/admin)
- StartedAt, CompletedAt
- Duration (generated column: CompletedAt - StartedAt)
- Notes
- CreatedAt, UpdatedAt

**Relationships:**

- BelongsTo: JobApplication
- HasMany: ApplicationNote (stage-specific notes), Interview

**Helper Methods:**

- `IsCompleted()` - Check if stage completed
- `IsInProgress()` - Check if stage ongoing
- `Complete()` - Mark stage as completed

**Features:**

- Auto-calculated duration via PostgreSQL
- Stage history tracking
- Multiple stages per application

---

### 3. **ApplicationDocument**

Document management (CV, cover letter, portfolio, etc.):

**Key Fields:**

- ID, ApplicationID, UserID
- DocumentType (cv/cover_letter/portfolio/certificate/transcript/other)
- FileName, FileURL, FileType, FileSize
- UploadedAt, IsVerified, VerifiedBy, VerifiedAt
- Notes
- CreatedAt, UpdatedAt

**Helper Methods:**

- `IsCV()`, `IsCoverLetter()` - Check document type
- `Verify()` - Mark document as verified with verifier ID

**Features:**

- Multiple documents per application
- Document verification workflow
- Admin verification tracking

---

### 4. **ApplicationNote**

Notes dan evaluations dari recruiters:

**Key Fields:**

- ID, ApplicationID, StageID (optional - stage-specific note)
- AuthorID (recruiter/admin), NoteType, NoteText
- Visibility (internal/public), Sentiment (positive/neutral/negative)
- IsPinned
- CreatedAt, UpdatedAt

**Enums:**

- NoteType: evaluation, feedback, reminder, internal (4 types)
- Visibility: internal, public (2 types)
- Sentiment: positive, neutral, negative (3 types)

**Helper Methods:**

- `IsInternal()`, `IsPublic()` - Check visibility
- `IsEvaluation()`, `IsFeedback()` - Check note type
- `IsPositive()`, `IsNegative()` - Check sentiment
- `Pin()`, `Unpin()` - Manage pinned status

**Features:**

- Can be linked to specific stage
- Sentiment analysis
- Pin important notes
- Public notes visible to candidate

---

### 5. **Interview**

Interview scheduling dan evaluation:

**Key Fields:**

- ID, ApplicationID, StageID, InterviewerID
- ScheduledAt, EndedAt
- InterviewType (online/onsite/hybrid)
- MeetingLink, Location
- Status (scheduled/completed/rescheduled/cancelled/no_show)
- Evaluation scores (4 dimensions):
  - OverallScore (0-100)
  - TechnicalScore (0-100)
  - CommunicationScore (0-100)
  - PersonalityScore (0-100)
- Remarks, FeedbackSummary
- CreatedAt, UpdatedAt

**Helper Methods:**

- `IsScheduled()`, `IsCompleted()`, `IsCancelled()`, `IsNoShow()` - Check status
- `IsOnline()`, `IsOnsite()` - Check interview type
- `Complete()`, `Cancel()`, `MarkNoShow()` - Status management
- `HasScores()` - Check if evaluated
- `CalculateAverageScore()` - Calculate average from 3 dimension scores

**Features:**

- Multi-dimensional evaluation
- Flexible interview types
- Meeting link for online interviews
- Location for onsite interviews
- Reschedule tracking

---

## Repository Interface (80+ methods)

### JobApplication CRUD (5 methods)

- `Create()`, `FindByID()`, `FindByJobAndUser()`
- `Update()`, `Delete()`

### Application Listing (9 methods)

- `List()` - List all applications dengan filter
- `ListByUser()` - User's applications
- `ListByJob()` - Job's applications
- `ListByCompany()` - Company's applications
- `UpdateStatus()`, `BulkUpdateStatus()`
- `GetApplicationsByStatus()`
- `MarkAsViewed()`, `ToggleBookmark()`, `GetBookmarkedApplications()`

### Application Search (3 methods)

- `SearchApplications()` - Advanced search
- `GetApplicationsWithHighScore()` - Filter by match score
- Plus advanced filtering in List methods

### Application Statistics (4 methods)

- `GetApplicationStats()` - Individual application stats
- `GetUserApplicationStats()` - User's overall stats
- `GetJobApplicationStats()` - Job's application stats
- `GetCompanyApplicationStats()` - Company's overall stats

### JobApplicationStage Operations (7 methods)

- `CreateStage()`, `FindStageByID()`, `UpdateStage()`, `CompleteStage()`
- `ListStagesByApplication()`, `GetCurrentStage()`, `GetStageHistory()`

### ApplicationDocument Operations (9 methods)

- `CreateDocument()`, `FindDocumentByID()`, `UpdateDocument()`, `DeleteDocument()`
- `ListDocumentsByApplication()`, `ListDocumentsByUser()`, `GetDocumentsByType()`
- `VerifyDocument()`, `GetUnverifiedDocuments()`

### ApplicationNote Operations (10 methods)

- `CreateNote()`, `FindNoteByID()`, `UpdateNote()`, `DeleteNote()`
- `ListNotesByApplication()`, `ListNotesByStage()`, `ListNotesByAuthor()`
- `GetPinnedNotes()`, `PinNote()`, `UnpinNote()`

### Interview Operations (13 methods)

- `CreateInterview()`, `FindInterviewByID()`, `UpdateInterview()`, `DeleteInterview()`
- `ListInterviewsByApplication()`, `ListInterviewsByInterviewer()`
- `GetUpcomingInterviews()`, `GetInterviewsByDateRange()`
- `UpdateInterviewStatus()`, `CompleteInterview()`, `RescheduleInterview()`, `CancelInterview()`

### Analytics & Reporting (7 methods)

- `GetApplicationTrends()` - Time-series trends
- `GetConversionFunnel()` - Hiring funnel metrics
- `GetAverageTimePerStage()` - Stage duration analysis
- `GetTopApplicants()` - Best applicants by score
- `GetApplicationSourceStats()` - Application source breakdown

### Bulk Operations (2 methods)

- `BulkCreateApplications()`, `BulkDeleteApplications()`

---

## Service Interface (70+ methods)

### Application Submission - Job Seeker (5 methods)

- `ApplyForJob()` - Submit application with documents
- `WithdrawApplication()` - Withdraw application
- `GetMyApplications()` - View own applications
- `GetApplicationDetail()` - View application detail
- `GetMyApplicationStats()` - View own statistics

### Application Review - Employer (6 methods)

- `GetJobApplications()` - View job's applications
- `GetCompanyApplications()` - View company's applications
- `GetApplicationForReview()` - Get application detail
- `MarkAsViewed()`, `ToggleBookmark()`, `GetBookmarkedApplications()`

### Status Workflow - Employer (7 methods)

- `MoveToScreening()`, `MoveToShortlist()`, `MoveToInterview()`
- `MakeOffer()`, `MarkAsHired()`, `RejectApplication()`
- `BulkUpdateStatus()` - Bulk status updates

### Stage Management (4 methods)

- `GetApplicationStages()`, `GetCurrentStage()`, `GetStageHistory()`
- `CompleteStage()` - Mark stage as complete

### Document Management (7 methods)

- `UploadApplicationDocument()`, `UpdateDocument()`, `DeleteDocument()`
- `GetApplicationDocuments()`, `GetDocumentsByType()`
- `VerifyDocument()`, `GetUnverifiedDocuments()`

### Notes Management - Employer (7 methods)

- `AddNote()`, `UpdateNote()`, `DeleteNote()`
- `GetApplicationNotes()`, `GetStageNotes()`
- `PinNote()`, `UnpinNote()`, `GetPinnedNotes()`

### Interview Scheduling (9 methods)

- `ScheduleInterview()`, `RescheduleInterview()`, `CancelInterview()`
- `CompleteInterview()` - Complete with evaluation scores
- `MarkInterviewNoShow()`
- `GetApplicationInterviews()`, `GetInterviewDetail()`
- `GetUpcomingInterviews()`, `GetInterviewsByDateRange()`
- `SendInterviewReminder()` - Send reminder notification

### Search & Filtering (3 methods)

- `SearchApplications()` - Advanced search
- `GetHighScoreApplications()` - Filter by match score
- `GetRecentApplications()` - Recent applications

### Analytics & Reporting (9 methods)

- `GetApplicationAnalytics()` - Individual application analytics
- `GetJobApplicationAnalytics()` - Job analytics
- `GetCompanyApplicationAnalytics()` - Company analytics
- `GetConversionFunnel()` - Funnel metrics
- `GetApplicationTrends()` - Trends over time
- `GetAverageTimePerStage()` - Stage duration
- `GetTopApplicants()` - Best applicants
- `GetApplicationSourceAnalytics()` - Source breakdown

### Notifications (4 methods)

- `NotifyApplicationReceived()` - Notify employer
- `NotifyStatusUpdate()` - Notify candidate
- `NotifyInterviewScheduled()`, `NotifyInterviewReminder()`

### Validation & Permissions (4 methods)

- `ValidateApplication()`, `CheckApplicationOwnership()`
- `CheckEmployerAccess()`, `CanApplyForJob()`

### Bulk Operations (3 methods)

- `BulkRejectApplications()`, `BulkMoveToStage()`, `ExportApplications()`

---

## Request DTOs (9)

1. **ApplyJobRequest** - Submit application dengan documents
2. **UploadDocumentRequest** - Upload document
3. **UpdateDocumentRequest** - Update document info
4. **AddNoteRequest** - Add note dengan type, visibility, sentiment
5. **UpdateNoteRequest** - Update note
6. **ScheduleInterviewRequest** - Schedule interview
7. **RescheduleInterviewRequest** - Reschedule interview
8. **CompleteInterviewRequest** - Complete dengan evaluation scores

---

## Response DTOs (15+)

1. **ApplicationListResponse** - Paginated list dengan stats
2. **ApplicationSummary** - Summary for listing
3. **ApplicationDetailResponse** - Complete detail dengan job, applicant, stages, documents, notes, interviews
4. **JobDetail** - Job info in application context
5. **ApplicantProfile** - Applicant profile dengan skills, education
6. **ListStats** - Statistics untuk list (viewed, bookmarked, match score)
7. **ApplicationAnalytics** - Detailed analytics dengan timeline, stage progress, document stats, interview stats
8. **TimelineEvent** - Timeline event
9. **StageProgress** - Stage progress detail
10. **DocumentStats** - Document statistics
11. **InterviewStats** - Interview statistics
12. **MatchAnalysis** - Match score breakdown
13. **ActivityLogEntry** - Activity log
14. **JobApplicationAnalytics** - Job analytics
15. **CompanyApplicationAnalytics** - Company analytics
16. **TimeSeriesData** - Time-series data point
17. **JobStats** - Job statistics

---

## Filters & Search Types (3)

1. **ApplicationFilter** - Basic filtering (status, job, user, company, score range, viewed, bookmarked, source, date range, sort)
2. **ApplicationSearchFilter** - Advanced search (keyword, job IDs, company IDs, statuses, score, sources, applied within, has documents, has interviews)
3. **InterviewFilter** - Interview filtering (status, type, scheduled date range, completed only)

---

## Statistics Types (12)

1. **ApplicationStats** - Individual application stats
2. **UserApplicationStats** - User's overall stats
3. **JobApplicationStats** - Job's application stats dengan source breakdown
4. **CompanyApplicationStats** - Company stats dengan monthly breakdown
5. **ApplicationTrend** - Trend data point
6. **ConversionFunnel** - Hiring funnel metrics
7. **StageTimeStats** - Average time per stage
8. **SourceStats** - Application source statistics
9. **SourceCount** - Source count
10. **JobPerformance** - Job performance metrics
11. **MonthlyCount** - Monthly count
12. **InterviewScores** - Interview evaluation scores

---

## Business Features

### 1. **Application Submission Workflow**

- Apply with resume + documents
- Match score calculation
- Duplicate prevention (one application per job per user)
- Application source tracking
- Withdraw functionality

### 2. **Hiring Stage Management**

- 8-stage workflow: applied → screening → shortlisted → interview → offered → hired/rejected/withdrawn
- Stage history tracking
- Duration calculation per stage
- Stage-specific notes
- Automatic stage transitions

### 3. **Document Management**

- Multiple document types (CV, cover letter, portfolio, certificate, transcript)
- Document verification workflow
- File metadata tracking (type, size)
- Admin verification with notes
- Document type filtering

### 4. **Notes & Collaboration**

- Internal vs public notes
- Note types: evaluation, feedback, reminder, internal
- Sentiment analysis (positive, neutral, negative)
- Pin important notes
- Stage-specific notes
- Note author tracking

### 5. **Interview Scheduling**

- Flexible interview types (online, onsite, hybrid)
- Meeting link for online interviews
- Location for onsite interviews
- Reschedule functionality
- Cancel functionality
- No-show tracking
- Reminder notifications

### 6. **Interview Evaluation**

- Multi-dimensional scoring:
  - Overall score
  - Technical score
  - Communication score
  - Personality score
- Remarks and feedback summary
- Average score calculation
- Interview completion tracking

### 7. **Application Tracking**

- Viewed by employer tracking
- Bookmark functionality
- Application timeline
- Activity log
- Status change history
- Stage progress tracking

### 8. **Search & Discovery**

- Advanced search with multiple criteria
- Filter by match score
- Filter by status
- Filter by date range
- Filter by source
- Filter by documents/interviews presence
- Sort by various fields

### 9. **Analytics & Reporting**

- Application trends over time
- Conversion funnel analysis
- Average time per stage
- Top applicants ranking
- Source effectiveness analysis
- Job performance metrics
- Company-level analytics
- Monthly breakdown

### 10. **Bulk Operations**

- Bulk status updates
- Bulk rejection with reason
- Bulk stage movement
- Export applications

### 11. **Notifications**

- Application received notification
- Status update notification
- Interview scheduled notification
- Interview reminder notification

---

## Technical Features

1. **GORM Integration**

   - Proper relationships dengan foreignKey & constraints
   - CASCADE delete untuk child entities
   - SET NULL untuk optional relationships
   - Generated column (duration) via PostgreSQL
   - Indexes untuk performance

2. **Unique Constraints**

   - One application per job per user
   - Composite unique index (JobID, UserID)

3. **Validation**

   - Comprehensive validation tags
   - Enum validation
   - Score range validation (0-100)
   - Business rule validation

4. **Timestamps**

   - Auto-managed CreatedAt & UpdatedAt
   - AppliedAt tracking
   - UploadedAt for documents
   - ScheduledAt, EndedAt for interviews
   - VerifiedAt for documents

5. **Helper Methods**

   - Status check methods
   - Action methods (Complete, Cancel, Verify)
   - Calculation methods (CalculateAverageScore)

6. **Filtering & Pagination**

   - Flexible filter structs
   - Date range filtering
   - Score range filtering
   - Boolean filters

7. **Analytics**
   - Time-series data
   - Conversion metrics
   - Source tracking
   - Performance metrics

---

## Statistics

- **Total Entities:** 5
- **Total Repository Methods:** ~80
- **Total Service Methods:** ~70
- **Total Request DTOs:** 9
- **Total Response DTOs:** 15+
- **Total Filter Types:** 3
- **Total Stats Types:** 12
- **Total Lines of Code:** ~850 (entities + repository + service)

---

## Integration Points

### Depends On:

- User domain (user_id reference)
- Job domain (job_id reference)
- Company domain (company_id reference)
- Admin domain (admin_users for verification, handling, interviewing)

### Used By:

- Notification service (status updates, interview reminders)
- Email service (application receipts, interview invitations)
- Analytics service (reporting)
- Export service (data export)

---

## Key Workflows

### 1. Application Flow (Job Seeker)

```
User browses jobs →
Applies with resume/documents →
Receives confirmation →
Tracks application status →
Can withdraw if needed
```

### 2. Review Flow (Employer)

```
Receives application →
Reviews profile & documents →
Adds notes/evaluation →
Moves through stages (screening → shortlist → interview) →
Schedules interview →
Evaluates candidate →
Makes offer/rejects
```

### 3. Interview Flow

```
Schedule interview →
Send invitation →
Send reminder (1 day before) →
Conduct interview →
Complete evaluation with scores →
Record feedback
```

---

## Next Steps

After Application domain completion:

1. **Admin & Master Domains** - AdminUser, AdminRole, SkillsMaster, BenefitsMaster
2. **Repository Implementation** - Implement all repository interfaces
3. **Service Implementation** - Implement all business logic
4. **Email Templates** - Design email templates for notifications

---
