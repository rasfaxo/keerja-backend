package utils

import (
	"errors"
	"time"
)

// ParseOptionalDateTime parses optional datetime string in RFC3339 format
func ParseOptionalDateTime(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil, errors.New("datetime must be in RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)")
	}

	return &t, nil
}

// ParseDateTime parses required datetime string in RFC3339 format
func ParseDateTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("datetime is required")
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, errors.New("datetime must be in RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)")
	}

	return t, nil
}

// MustBeFutureTime validates that time is in the future
func MustBeFutureTime(t time.Time) error {
	if t.Before(time.Now()) {
		return errors.New("datetime must be in the future")
	}
	return nil
}

// MustBePastTime validates that time is in the past
func MustBePastTime(t time.Time) error {
	if t.After(time.Now()) {
		return errors.New("datetime must be in the past")
	}
	return nil
}

// MustBeWithinRange validates that time is within specified range
func MustBeWithinRange(t time.Time, min, max time.Time) error {
	if t.Before(min) || t.After(max) {
		return errors.New("datetime must be within the specified range")
	}
	return nil
}

// ValidateDateRange validates that start date is before end date
func ValidateDateRange(start, end time.Time) error {
	if start.After(end) || start.Equal(end) {
		return errors.New("start date must be before end date")
	}
	return nil
}

// ParseOptionalDate parses optional date string (YYYY-MM-DD format)
func ParseOptionalDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, errors.New("date must be in YYYY-MM-DD format")
	}

	return &t, nil
}

// ParseDate parses required date string (YYYY-MM-DD format)
func ParseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("date is required")
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errors.New("date must be in YYYY-MM-DD format")
	}

	return t, nil
}

// FormatDateTime formats time to RFC3339 string
func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatDate formats time to YYYY-MM-DD string
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// TimePtr converts time to pointer
func TimePtr(t time.Time) *time.Time {
	return &t
}

// IsToday checks if time is today
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

// IsPast checks if time is in the past
func IsPast(t time.Time) bool {
	return t.Before(time.Now())
}

// IsFuture checks if time is in the future
func IsFuture(t time.Time) bool {
	return t.After(time.Now())
}

// AddBusinessDays adds business days (excluding weekends) to a date
func AddBusinessDays(t time.Time, days int) time.Time {
	current := t
	daysAdded := 0

	for daysAdded < days {
		current = current.AddDate(0, 0, 1)

		// Skip weekends
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			daysAdded++
		}
	}

	return current
}

// GetStartOfDay returns the start of the day (00:00:00)
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay returns the end of the day (23:59:59)
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// DaysBetween calculates the number of days between two dates
func DaysBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}
