# User Preferences API

This document describes the new user preferences endpoints.

Endpoints

1) GET /api/v1/users/me/preferences
   - Authenticated endpoint (Bearer token)
   - Returns current user's preferences as JSON. `language_preference` is nullable (can be null).

2) PUT /api/v1/users/me/preferences
   - Authenticated endpoint (Bearer token)
   - Partial update / upsert - if user has no preferences yet, this will create a new preferences row.
   - Validation rules:
     - `language_preference`: optional, max length 10 characters
     - `theme_preference`: optional, one of `light`, `dark`
     - `preferred_job_type`: optional, max length 50
     - `profile_visibility`: optional, one of `public`, `private`, `recruiter-only`
     - `preferred_salary_min` / `preferred_salary_max`: optional, must be >= 0

Example request

PUT /api/v1/users/me/preferences
{
  "language_preference": "en",
  "theme_preference": "dark",
  "profile_visibility": "private",
  "preferred_salary_min": 50000
}

Example response

200 OK
{
  "success": true,
  "message": "Preferences updated successfully",
  "data": {
    "id": 12,
    "language_preference": "en",
    "theme_preference": "dark",
    "preferred_job_type": "remote",
    "updated_at": "2025-11-28T12:34:56Z"
  }
}
