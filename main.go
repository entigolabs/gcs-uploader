package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

func main() {
	ctx := context.Background()
	flags := Flags{}

	err := flags.getValues()
	if err != nil {
		fmt.Println("Error getting flag values:", err)
		return
	}

	creds, err := google.FindDefaultCredentials(ctx, storage.ScopeReadWrite)
	if err != nil {
		fmt.Println("Error finding credentials:", err)
		return
	}

	if creds.JSON != nil {
		fmt.Println("Credentials file:", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
		var credFile struct {
			ClientEmail string `json:"client_email"`
		}
		if err := json.Unmarshal(creds.JSON, &credFile); err == nil && credFile.ClientEmail != "" {
			fmt.Println("Service account:", credFile.ClientEmail)
		}
	}
	if creds.ProjectID != "" {
		fmt.Println("Project ID:", creds.ProjectID)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("Error creating GCS client:", err)
		return
	}
	defer func() { _ = client.Close() }()

	err = uploadFilesToGCS(ctx, client, flags)
	if err != nil {
		fmt.Println("Error uploading files to GCS:", err)
		return
	}

	uniqueTags, err := getUniqueGCSObjectTags(ctx, client, flags)
	if err != nil {
		fmt.Println("Error getting GCS object tags:", err)
		return
	}

	sortedUniqueTags := sortTags(uniqueTags)

	tagsToDelete := getTagsToDelete(flags, sortedUniqueTags)

	err = deleteObjectsWithTags(ctx, client, flags, tagsToDelete)
	if err != nil {
		fmt.Println("Error deleting objects with tags:", err)
		return
	}
}
