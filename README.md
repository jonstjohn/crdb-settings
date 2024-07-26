# CRDB Settings

This tool provides a way to examine CockroachDB cluster settings for specific releases or across
multiple releases. It provides the following:

1. Discover new versions of CockroachDB and store default cluster settings
2. Serve a REST API for examining settings

## Command-line

These commands require the database URL to be provided via the `--url` flag.

### Releases

Update releases stored in the database from an external source (e.g., authoritative yaml file):

```
./crdb-settings releases update --url $DBURL
```

List releases from the database:

```
./crdb-settings releases list --url $DBURL
```

### Settings

Update settings stored in database (by default, start with most recent release and go backwards):

```
./crdb-settings settings update --url $DBURL
```

List settings for a specific version:

```
./crdb-settings settings list [version] --url $DBURL
```

Compare settings across 2 versions:

```
./crdb-settings settings compare [version1] [version2] --url $DBURL
```

## REST API

The REST API is defined via an OpenAPI spec and can be served via a web server.

### OpenAPI Spec

The following operations are supported:

1. `/settings/release/[release]`
2. `/settings/compare/[release1]..[release2]`

### REST web server

To run the web server:

```
./crdb-settings api serve --url $DBURL
```

### Deploy to Google App Engine

To deploy to Google App engine, run:

```
gcloud app deploy --project $PROJECT
```

App engine must have access to the secret `CRDB_SETTINGS_DBURL` in the project.

The current deployment uses Google Cloud Build to automatically deploy on push (dev) or tag (prod).

## Public Access

To access the tool without needing to do any installation, the REST API and a web application is exposed via public URLs.

Web application: https://crdb-settings.distributedbites.com

REST API:
* View settings by release endpoint: https://crdb-settings-api.distributedbites.com/settings/release/v23.1.22
* Compare release settings endpoint: https://crdb-settings-api.distributedbites.com/settings/release/v23.1.22..v23.2.7
