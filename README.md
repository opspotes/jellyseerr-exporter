# Jellyseerr Exporter

Prometheus exporter for [Jellyseerr](https://github.com/Fallenbagel/jellyseerr).

## Usage

```bash
docker run --rm -p 9850:9850 ghcr.io/opspotes/jellyseerr-exporter:latest \
  "--jellyseerr.address=https://jellyseerr.example.com" \
  "--jellyseerr.apiKey=examplesecretapikey"
```

## Metrics

Two main metric groups are exported: Requests and Users.

### Requests

The requests on the Jellyseerr server are counted. Request counts have the following labels:

| Label            |                      Description                       | Configurable |
| :--------------- | :----------------------------------------------------: | -----------: |
| `request_status` |  The approval status for the requests (e.g. Approved)  |           no |
| `media_status`   | The media status for requested items (e.g. Available)  |           no |
| `media_type`     |       The category of request media (e.g. movie)       |           no |
| `is_4k`          |      Requested on a 4k tagged service (e.g. true)      |           no |
| `genre`          |       The main genre for a requested media item        |          yes |
| `company`        | The production company or network for a requested item |          yes |

> ⚠️ Collecting Genre/Company info can take a lot of time with large request quantities

### Users

User request counts of an Jellyseerr server are collected with the following labels:

| Label   |          Description          | Configurable |
| :------ | :---------------------------: | -----------: |
| `email` | The email address of the user |           no |

## Configuration

| Flag                 |                               Description                               | Default |
| :------------------- | :---------------------------------------------------------------------: | :------ |
| `log`                |                 Sets the logging level for the exporter                 | `fatal` |
| `jellyseerr.address` |                   The URI of the Jellyseerr instance                    |         |
| `jellyseerr.apiKey`  |              The admin API key of the Jellyseerr instance               |         |
| `jellyseerr.locale`  |                  The locale of the Jellyseerr instance                  | `en`    |
| `fullData`           | Do not collect genre and company to reduce API requests and cardinality | `false` |

You **must** provide the Overseerr address and API key!

## Original project

This project is a fork from [`WillFantom/overseerr-exporter`](https://github.com/WillFantom/overseerr-exporter), with some improvements. Thanks a lot to him for this work released under a GPL license.
