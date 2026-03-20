# gcs-uploader

CLI tool for uploading files to Google Cloud Storage Bucket with version tagging and automatic cleanup of old versions.
Tags are sorted as semver and oldest versions are deleted based on `--num-latest-tags-to-keep`.

## Flags

### Required

| Flag                        | Description                                  |
| --------------------------- | -------------------------------------------- |
| `--bucket`                  | GCS bucket name                              |
| `--source-directory`        | Local source directory                       |
| `--num-latest-tags-to-keep` | Number of latest tags to keep                |
| `--tag`                     | Tag for uploaded files (format: `key=value`) |

### Optional

| Flag                     | Default                   | Description                         |
| ------------------------ | ------------------------- | ----------------------------------- |
| `--target-directory`     | `""` (bucket root)        | GCS target directory                |
| `--concurrent-uploads`   | `100`                     | Number of concurrent uploads        |
| `--concurrent-deletions` | `100`                     | Number of concurrent deletions      |
| `--cache-control`        | `max-age=31536000,public` | Cache-Control header for files      |
| `--index-cache-control`  | `no-cache`                | Cache-Control header for index.html |

## Authentication

Uses [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials).

Options:

- Set `GOOGLE_APPLICATION_CREDENTIALS` env var to a service account key file path
- Run `gcloud auth application-default login` for local development
- Automatic on GCE, Cloud Run, GKE with attached service account

## Docker Images

| Image                                    | Description                       |
| ---------------------------------------- | --------------------------------- |
| `ghcr.io/entigolabs/gcs-uploader`        | Base image with gcs-uploader only |
| `ghcr.io/entigolabs/gcs-uploader-gcloud` | Includes gcloud CLI               |

## Quick Start (Docker)

```
docker run \
  -v $(pwd)/source:/source \
  -v /path/to/credentials.json:/credentials.json \
  -e GOOGLE_APPLICATION_CREDENTIALS=/credentials.json \
  ghcr.io/entigolabs/gcs-uploader:latest \
  --bucket mybucket \
  --source-directory /source \
  --target-directory target/ \
  --num-latest-tags-to-keep 3 \
  --tag version=1.0.0
```

## CI/CD Usage (with gcloud CLI)

```yaml
# Bitbucket Pipelines example
- step:
    name: Deploy
    image: ghcr.io/entigolabs/gcs-uploader-gcloud:latest
    script:
      - echo $BASE64_GOOGLE_APPLICATION_CREDENTIALS | base64 -d > credentials.json
      - export GOOGLE_APPLICATION_CREDENTIALS=credentials.json
      - gcloud auth activate-service-account --key-file credentials.json
      - gcs-uploader --bucket mybucket --source-directory source --num-latest-tags-to-keep 3 --tag version=1.0.0
      - gcloud compute url-maps invalidate-cdn-cache my-url-map --path '/*'
```
