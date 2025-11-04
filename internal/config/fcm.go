package config

import (
	"context"
	"fmt"
	"log"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	fcmApp    *firebase.App
	fcmClient *messaging.Client
	fcmOnce   sync.Once
	fcmMutex  sync.RWMutex
)

// InitFCM initializes Firebase Cloud Messaging client
// Reference: https://firebase.google.com/docs/cloud-messaging/admin/send-messages
func InitFCM(cfg *Config) error {
	if !cfg.FCMEnabled {
		log.Println("INFO: FCM is disabled. Push notifications will not be sent.")
		return nil
	}

	var initErr error
	fcmOnce.Do(func() {
		ctx := context.Background()

		// Initialize Firebase App with service account
		opt := option.WithCredentialsFile(cfg.FCMCredentialsFile)
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize Firebase app: %w", err)
			return
		}

		// Get messaging client
		client, err := app.Messaging(ctx)
		if err != nil {
			initErr = fmt.Errorf("failed to get messaging client: %w", err)
			return
		}

		fcmMutex.Lock()
		fcmApp = app
		fcmClient = client
		fcmMutex.Unlock()

		log.Printf("SUCCESS: FCM initialized for project: %s", cfg.FCMProjectID)
	})

	return initErr
}

// GetFCMClient returns the initialized FCM messaging client
func GetFCMClient() (*messaging.Client, error) {
	fcmMutex.RLock()
	defer fcmMutex.RUnlock()

	if fcmClient == nil {
		return nil, fmt.Errorf("FCM client not initialized. Call InitFCM() first or enable FCM_ENABLED in config")
	}

	return fcmClient, nil
}

// GetFCMApp returns the initialized Firebase app
func GetFCMApp() (*firebase.App, error) {
	fcmMutex.RLock()
	defer fcmMutex.RUnlock()

	if fcmApp == nil {
		return nil, fmt.Errorf("Firebase app not initialized. Call InitFCM() first or enable FCM_ENABLED in config")
	}

	return fcmApp, nil
}

// IsFCMEnabled checks if FCM is enabled and initialized
func IsFCMEnabled() bool {
	fcmMutex.RLock()
	defer fcmMutex.RUnlock()

	return fcmClient != nil
}

// CloseFCM closes the FCM client (if needed for graceful shutdown)
// Note: Firebase Admin SDK v4 doesn't require explicit cleanup
// but this function exists for consistency and future compatibility
func CloseFCM() {
	fcmMutex.Lock()
	defer fcmMutex.Unlock()

	if fcmClient != nil {
		log.Println("INFO: FCM client closed")
		fcmClient = nil
		fcmApp = nil
	}
}
