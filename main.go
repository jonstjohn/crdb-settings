/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"github.com/jonstjohn/crdb-settings/cmd"
	"github.com/jonstjohn/crdb-settings/pkg/api"
	"log"
	"net/http"
	"os"
)

func main() {
	cmd.Execute()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

}

func getDbUrl() (string, error) {
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer c.Close()

	name := fmt.Sprintf("projects/%s/CRDB_SETTINGS_DBURL/latest", os.Getenv("GOOGLE_CLOUD_PROJECT"))

	req := &secretmanagerpb.GetSecretRequest{Name: name}
	secret, err := c.GetSecret(ctx, req)
	if err != nil {
		return "", err
	}
	return secret.String(), nil
}

func entry() {
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
