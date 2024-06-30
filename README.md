# Tars.Run

A simple short link application backend written in Go. 

Demo site: <https://tars.run>

> [!NOTE]
> This site may have downtime and loss of data is possible. I recommend self-hosting this application if you require
> durability.

## Installation

The prod version runs in a single docker container with a dependency on [Litestream].

In order to use this project's Dockerfile, you'll need to understand how to use
[Litestream] and fill out the appropriate environment variables.

Environment variables can be found in `.env.example`.

```shell
# create a .env
cp .env.example .env
# or feel free to use any other method to inject the required environment variables into your shell.
```

The application uses [goose] to perform database migrations. To bootstrap the database, you'll
need to install [goose] and create a SQLite database file.

```shell
go get github.com/pressly/goose/v3/cmd/goose
mkdir -p data
touch data/shorty.db.up
task db:migration:up
```

In development, the application can be started using `air`. Simply run the `air` command and a
hot-reloading webserver will be started. The migrations must be applied in order to create the
tables.

```shell
# development
task dev
```

**Note: hot-reloading requires [air] and [templ]**

There is a [Taskfile](/Taskfile.yml) with a large number of commands.

## Web version

### How it works

This version of the application is self-contained using Go templates and embedded files. Adding
new links is done via a `POST` request to the `v1/links` endpoint which saves it to the database
and renders it in the template. All other endpoints are `GET` requests and accessible via the
web pages `/` and `/{hash}/analytics`.

When developing, if adding new Tailwind classes you may be required to re-run the Tailwind
compiler. To do that just run `task tailwind` which will output a minified CSS file named `theme.css`
in `assets/static/css`. This is used in the `assets/static/view/layout/base.templ` file. Both Alpine.js and
Tailwind are self-contained in this binary for portability and to reduce network calls.

All templates use [templ] as it provides type safety over Go's standard library template package.

### Endpoints

| Method | URL         | Handler            | Action                  | Authentication |
|--------|-------------|--------------------|-------------------------|----------------|
| GET    | `/`         | `handleHome`       | Retrieves all links     | False          |
| GET    | `/:hash`    | `handleRedirect`   | Retrieves a single link | False          |
| POST   | `/v1/links` | `handleCreateLink` | Creates a new link      | False          |

## Self Hosting

This application is designed to be deployed as a container.

The following environment variables are required to run this if using Litestream.

```shell
LITESTREAM_SECRET_ACCESS_KEY=
LITESTREAM_ACCESS_KEY_ID=
DB_PATH=./data/data.db
LITESTREAM_BUCKET=
LITESTREAM_ENDPOINT=
LITESTREAM_RETENTION=72h
DB_SYNC_INTERVAL=60s
LITESTREAM_DB_NAME=
```

Refer to `.env.example` for more configuration options such as Social Media accounts and application name.

## Stack

- [go-chi](https://github.com/go-chi/chi)
- Sqlite (with [Litestream])
- [templ]

## Development

You will require the following tools:

- [templ]
- [air]
- [goose]

## License

Licenced under Apache 2.0, see [LICENSE](/LICENSE).

[air]: https://github.com/air-verse/air

[templ]: https://templ.guide

[litestream]: https://litestream.io

[goose]: https://github.com/pressly/goose