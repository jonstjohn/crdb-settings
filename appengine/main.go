package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/api"
	"log"
	"net/http"
	"os"
)

func getDbUrl() (string, error) {
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer c.Close()

	name := fmt.Sprintf("projects/%s/secrets/CRDB_SETTINGS_DBURL/versions/latest", os.Getenv("GOOGLE_CLOUD_PROJECT"))

	//req := &secretmanagerpb.GetSecretRequest{Name: name}
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: name}
	//secret, err := c.GetSecret(ctx, req)
	secret, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}
	return string(secret.Payload.Data), nil
}

func main() {
	url, err := getDbUrl()
	if err != nil {
		log.Fatal(err)
	}
	sh := api.SettingsHandler{Url: url}
	http.HandleFunc("/", sh.ServeHTTP)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
