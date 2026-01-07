# CRDB Settings (and other things) Tool

This tool was originally designed to provide a way to examine CockroachDB cluster settings for specific releases or across
multiple releases. However, it was recently modified to also support querying metrics. In the future, it may be renamed or split into different repositories.

It provides the following:

1. Discover new versions of CockroachDB and store default cluster settings and metrics
2. Serve a REST API for examining settings and metrics

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

Show settings details:

```
./crdb-settings settings detail --setting [setting] --url $DBURL
```

Find Github issues related to a setting:

```
./crdb-settings settings github --setting [setting] --url $DBURL
```

### Metrics

Update metrics stored in database (by default, start with most recent release and go backwards):

```
./crdb-settings metrics update --url $DBURL --release=recent-50
```

### Github

Update settings from Github mentions:

```
./crdb-settings settings github --url $DBURL
```


## REST API

The REST API is defined via an OpenAPI spec and can be served via a web server.

The following operations are supported:

1. `/settings/release/[release]`
2. `/settings/compare/[release1]..[release2]`
3. `/settings/detail/[setting]`
4. `/metrics/release/[release]`
5. `/metrics/compare/[release1]..[release2]`


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

Web application: https://distributedbites.com

REST API:
* View settings by release endpoint: https://api.distributedbites.com/settings/release/v23.1.22
* Compare release settings endpoint: https://api.distributedbites.com/settings/compare/v23.1.22..v23.2.7

