package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/nahuelojea/handballscore/config/secretmanager"
	"google.golang.org/api/option"
)

var (
	FirebaseApp      *firebase.App
	FirebaseDBClient *db.Client
	once             sync.Once
	initErr          error
)

// initialize attempts to initialize the Firebase app and client.
func initialize() {
	fmt.Println("Attempting Firebase initialization...")

	firebaseSecret, err := secretmanager.GetFirebaseSecret("firabase")
	if err != nil {
		initErr = fmt.Errorf("error getting Firebase secret: %w", err)
		// fmt.Println(initErr) // Removed
		return
	}

	secretJSON, err := json.Marshal(firebaseSecret)
	if err != nil {
		initErr = fmt.Errorf("error marshalling Firebase secret to JSON: %w", err)
		// fmt.Println(initErr) // Removed
		return
	}

	opt := option.WithCredentialsJSON(secretJSON)
	if firebaseSecret.ProjectID == "" {
		initErr = fmt.Errorf("ProjectID is empty in Firebase secret, cannot form DatabaseURL")
		// fmt.Println(initErr) // Removed
		return
	}
	config := &firebase.Config{
		DatabaseURL: "https://" + firebaseSecret.ProjectID + ".firebaseio.com",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		initErr = fmt.Errorf("error initializing Firebase app: %w", err)
		// fmt.Println(initErr) // Removed
		return
	}
	FirebaseApp = app

	client, err := app.Database(context.Background())
	if err != nil {
		initErr = fmt.Errorf("error initializing Firebase Realtime Database client: %w", err)
		// fmt.Println(initErr) // Removed
		return
	}
	FirebaseDBClient = client

	// fmt.Println("Firebase initialized successfully via sync.Once.") // Removed
}

// GetFirebaseDBClient returns the Firebase Realtime Database client,
// initializing it on the first call.
func GetFirebaseDBClient() (*db.Client, error) {
	once.Do(initialize)
	if initErr != nil {
		return nil, initErr
	}
	if FirebaseDBClient == nil && initErr == nil {
		return nil, fmt.Errorf("Firebase DB client is nil after initialization without an error; this indicates a bug in initialization logic")
	}
	return FirebaseDBClient, nil
}

// GetFirebaseApp returns the Firebase App,
// initializing it on the first call if not already initialized.
// This might be useful if other Firebase services (not just DB) are needed later.
func GetFirebaseApp() (*firebase.App, error) {
	once.Do(initialize)
	if initErr != nil {
		return nil, initErr
	}
	if FirebaseApp == nil && initErr == nil {
		return nil, fmt.Errorf("Firebase App is nil after initialization without an error; this indicates a bug")
	}
	return FirebaseApp, nil
}
