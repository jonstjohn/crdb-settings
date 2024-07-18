# CRDB Settings

This tool provides a way to examine CockroachDB cluster settings for specific releases or across
multiple releases. It provides the following:

1. Discover new versions of CockroachDB and store default cluster settings
2. Serve a REST API for examining settings
3. Deploy a web application for exploring settings

## Command-line

These commands require the `$CRDB_SETTINGS_DB_URL` be set, of the `--url` flag to be provided.

### Releases

Update releases stored in the database from an external source (e.g., authoritative yaml file):

```
./crdb-settings releases update
```

List releases from the database:

```
./crdb-settings releases list
```

### Settings

Update settings stored in database (by default, start with most recent release and go backwards):

```
./crdb-settings settings update
```

List settings for a specific version:

```
./crdb-settings settings list [version]
```

Compare settings across 2 versions:

```
./crdb-settings settings compare [version1] [version2]
```

## REST API

The REST API is defined via an OpenAPI spec and can be served via a web server.

### OpenAPI Spec

The following operations are supported:

1. `/settings/list/version`
2. `/settings/compare/version1..version2`

### REST web server

To run the web server:

```
./crdb-settings api serve
```

## Web Application

The web application is written in ReactJS.

To deploy the web application:

TODO

## AWS Lambda Deployment

To ensure that settings are kept up-to-date with the most recent releases, releases and settings
can be deployed to AWS lambda where they can be executed periodically.

TODO



