package main

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func deleteObjectsWithTags(ctx context.Context, client *storage.Client, flags Flags, tagsToDelete []string) error {
	bucket := client.Bucket(flags.Bucket)

	it := bucket.Objects(ctx, nil)

	var wg sync.WaitGroup
	concurrencyChan := make(chan struct{}, flags.ConcurrentDeletions)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		wg.Add(1)
		concurrencyChan <- struct{}{}
		go func(attrs *storage.ObjectAttrs) {
			defer func() {
				<-concurrencyChan
				wg.Done()
			}()

			for _, tagValue := range attrs.Metadata {
				if contains(tagsToDelete, tagValue) {
					err := bucket.Object(attrs.Name).Delete(ctx)
					if err != nil {
						fmt.Printf("Error deleting object %s: %v\n", attrs.Name, err)
						return
					}
					fmt.Printf("Deleted object: %s, Tag: %s\n", attrs.Name, tagValue)
					break
				}
			}
		}(attrs)
	}
	wg.Wait()
	return nil
}
