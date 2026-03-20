package main

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func getUniqueGCSObjectTags(ctx context.Context, client *storage.Client, flags Flags) ([]string, error) {
	bucket := client.Bucket(flags.Bucket)

	it := bucket.Objects(ctx, nil)

	unique := make(map[string]bool)
	uniqueTags := []string{}
	mutex := sync.Mutex{}

	var wg sync.WaitGroup
	concurrencyChan := make(chan struct{}, flags.ConcurrentUploads)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		concurrencyChan <- struct{}{}
		go func(attrs *storage.ObjectAttrs) {
			defer func() {
				<-concurrencyChan
				wg.Done()
			}()

			for _, tagValue := range attrs.Metadata {
				mutex.Lock()
				if !unique[tagValue] {
					unique[tagValue] = true
					uniqueTags = append(uniqueTags, tagValue)
				}
				mutex.Unlock()
			}
		}(attrs)
	}
	wg.Wait()

	return uniqueTags, nil
}

func sortTags(versions []string) []string {
	sortedVersions := make([]string, len(versions))
	copy(sortedVersions, versions)

	sort.SliceStable(sortedVersions, func(i, j int) bool {
		return compareVersions(sortedVersions[i], sortedVersions[j]) < 0
	})

	return sortedVersions
}

func getTagsToDelete(flags Flags, tags []string) []string {
	if len(tags) <= flags.NumLatestTagsToKeep {
		fmt.Println("Objects with these tags will remain:", tags)
		fmt.Println("No objects will be deleted.")
		return []string{}
	}
	tagsToKeep := tags[len(tags)-flags.NumLatestTagsToKeep:]
	tagsToDelete := tags[:len(tags)-flags.NumLatestTagsToKeep]

	fmt.Println("Objects with these tags will remain:", tagsToKeep)
	fmt.Println("Objects with these tags will be deleted:", tagsToDelete)

	return tagsToDelete
}
