package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func getContentType(file *os.File) (string, error) {
	ext := filepath.Ext(file.Name())
	switch ext {
	case ".js":
		return "text/javascript", nil
	case ".css":
		return "text/css", nil
	default:
		mime, err := mimetype.DetectReader(file)
		if err != nil {
			return "", err
		}
		return mime.String(), nil
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func compareVersions(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		numA, _ := strconv.Atoi(partsA[i])
		numB, _ := strconv.Atoi(partsB[i])

		if numA < numB {
			return -1
		} else if numA > numB {
			return 1
		}
	}

	if len(partsA) < len(partsB) {
		return -1
	} else if len(partsA) > len(partsB) {
		return 1
	}

	return 0
}
