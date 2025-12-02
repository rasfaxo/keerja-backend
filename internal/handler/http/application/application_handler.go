package applicationhandler

import (
	"keerja-backend/internal/domain/application"
)

// ApplicationHandler handles application-related HTTP requests
type ApplicationHandler struct {
	appService application.ApplicationService
}

// NewApplicationHandler creates a new instance of ApplicationHandler
func NewApplicationHandler(appService application.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}
