# CRDB Settings

This tool provides a way to examine CockroachDB cluster settings for specific releases or across
multiple releases. It provides the following:

1. Discover new versions of CockroachDB and store default cluster settings
2. Serve a REST API for examining settings
3. Deploy a web application for exploring settings (coming soon)

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
3. `/settings/history/[setting]` (coming soon)

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

## Web Application

The web application is written in ReactJS.

To deploy the web application:

TODO


