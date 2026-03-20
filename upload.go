package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
)

func uploadFilesToGCS(ctx context.Context, client *storage.Client, flags Flags) error {
	bucket := client.Bucket(flags.Bucket)

	tagKey, tagValue := parseTag(flags.Tag)

	var wg sync.WaitGroup
	concurrencyChan := make(chan struct{}, flags.ConcurrentUploads)

	err := filepath.Walk(flags.SourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			concurrencyChan <- struct{}{}
			wg.Add(1)
			go func(path string) {
				defer func() {
					<-concurrencyChan
					wg.Done()
				}()

				relativePath, _ := filepath.Rel(flags.SourceDirectory, path)
				gcsPath := filepath.ToSlash(filepath.Join(flags.TargetDirectory, relativePath))

				file, err := os.Open(path)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer func() { _ = file.Close() }()

				contentType, err := getContentType(file)
				if err != nil {
					fmt.Println("Error getting content type:", err, "File:", path)
					return
				}
				if _, err := file.Seek(0, 0); err != nil {
					fmt.Println("Error seeking file:", err, "File:", path)
					return
				}

				cacheControl := flags.DefaultCacheControl
				if filepath.Base(gcsPath) == "index.html" {
					cacheControl = flags.IndexHTMLCacheControl
				}

				fmt.Println("Uploading:", gcsPath, "ContentType:", contentType, "Tag:", flags.Tag)

				obj := bucket.Object(gcsPath)
				writer := obj.NewWriter(ctx)
				writer.ContentType = contentType
				writer.CacheControl = cacheControl
				writer.Metadata = map[string]string{
					tagKey: tagValue,
				}

				if _, err := io.Copy(writer, file); err != nil {
					fmt.Println("Error uploading file:", err, "File:", path)
					_ = writer.Close()
					return
				}
				if err := writer.Close(); err != nil {
					fmt.Println("Error closing writer:", err, "File:", path)
					return
				}
			}(path)
		}
		return nil
	})
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func parseTag(tag string) (string, string) {
	parts := strings.SplitN(tag, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "tag", tag
}
