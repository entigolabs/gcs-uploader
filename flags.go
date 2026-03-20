package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	NumLatestTagsToKeep   int
	SourceDirectory       string
	TargetDirectory       string
	Bucket                string
	Tag                   string
	ConcurrentUploads     int
	ConcurrentDeletions   int
	DefaultCacheControl   string
	IndexHTMLCacheControl string
}

func (c *Flags) getValues() error {
	flag.IntVar(&c.NumLatestTagsToKeep, "num-latest-tags-to-keep", 0, "Number of latest tags to keep")
	flag.StringVar(&c.SourceDirectory, "source-directory", "", "Source directory")
	flag.StringVar(&c.TargetDirectory, "target-directory", "", "Target directory")
	flag.StringVar(&c.Bucket, "bucket", "", "GCS bucket name")
	flag.StringVar(&c.Tag, "tag", "", "Tag")
	flag.IntVar(&c.ConcurrentUploads, "concurrent-uploads", 100, "Number of concurrent uploads")
	flag.IntVar(&c.ConcurrentDeletions, "concurrent-deletions", 100, "Number of concurrent deletions")
	flag.StringVar(&c.DefaultCacheControl, "cache-control", "max-age=31536000,public", "Cache-Control header for uploaded files")
	flag.StringVar(&c.IndexHTMLCacheControl, "index-cache-control", "no-cache", "Cache-Control header for index.html")
	flag.Parse()

	if c.NumLatestTagsToKeep == 0 || c.SourceDirectory == "" || c.Bucket == "" || c.Tag == "" {
		return fmt.Errorf("all flags must be set")
	}

	fmt.Println("Number of latest tags to keep:", c.NumLatestTagsToKeep)
	fmt.Println("Source directory:", c.SourceDirectory)
	fmt.Println("Target directory:", c.TargetDirectory, "(bucket root if empty)")
	fmt.Println("GCS bucket name:", c.Bucket)
	fmt.Println("Tag:", c.Tag)

	return nil
}
