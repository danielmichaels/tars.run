# Short Links

A simple short link application backend written in Go. This is exists solely for
learning and personal purposes only.

## Installation

The prod version runs in a single docker container and is deployed to my personal CapRover instance.

In order to use this projects Dockerfile, you'll need to understand how to use
[Litestream] (https://litestream.io) and fill out the appropriate environment variables.
Environment variables can be found in `.env.example`.

```shell
# create a .env
cp .env.example .env
# or feel free to use any other method to inject the required environment variables into your shell.
```

In development, the application can be started using `air`. Simply run the `air` command and a
hot-reloading webserver will be started. The migrations must be applied in order to create the
tables.

```shell
# development
make db/migration/up
make run/api
# for hot-reloading (requires air)
make run/api/development # or just `air`
```
**Note: hot-reloading requires [air](https://github.com/cosmtrek/air)**

There is a [Makefile](/Makefile) with a large number of commands.

## How it works
The following explanation assumes a frontend exists to send the data and act on the responses.
For brevity, this only shows the server request, response cycle and makes assumptions about how
certain parts are handled on the frontend.

A POST request that comes to the server's `/v1/links` with the correct payload will create a new
entry in the database with a user supplied link. The server will return a response which
contains the `short_url`.

For example.

```shell
❯ curlie http://localhost:1987/v1/links link=danielms.site
HTTP/2 201
server: nginx
date: Sat, 23 Oct 2021 10:07:16 GMT
content-type: application/json
content-length: 150
location: /v1/links/ApxmhgiNHcc
vary: Origin
vary: Access-Control-Request-Method

{
    "link": {
        "id": 37,
        "created_at": "2021-10-23T10:07:16Z",
        "original_url": "danielms.site",
        "hash": "ApxmhgiNHcc"
    },
    "short_url": "https://localhost:1988/ApxmhgiNHcc"
}
```

When a user visits a valid `short_url` such as above, they will be redirected to its mapped
address, in this case, see `original_url` in the example above.

If you look closely, the `original_url` actually returns a different port number. In production, this may be an entirely different URL.
This is because it is actually the Next.js frontend. When a user visits
`https://localhost:1988/ApxmhgiNHc` it will make a call to the backend and return the following.

```shell
❯ curlie http://localhost:1987/v1/links/ApxmhgiNHcc
HTTP/2 200
server: nginx
date: Sat, 23 Oct 2021 10:13:34 GMT
content-type: application/json
content-length: 107
vary: Origin
vary: Access-Control-Request-Method

{
    "link": {
        "id": 37,
        "created_at": "2021-10-23T10:07:16Z",
        "original_url": "danielms.site",
        "hash": "ApxmhgiNHcc"
    }
}
```

The frontend will then redirect the user to the `original_url`. During the initial lookup,
another database event happens - the user's information such as IP address, user-agent and the
datetime are recorded for analytics purposes.

From the frontend, users can view their shortened URL analytics. This triggers a request to the
server which if a valid short link is found, will return a response similar to below.

```shell
HTTP/2 200
server: nginx
date: Sat, 23 Oct 2021 10:18:20 GMT
content-type: application/json
content-length: 213
vary: Origin
vary: Access-Control-Request-Method

{
    "analytics": [
        {
            "id": 21,
            "ip_address": "102.130.32.28",
            "user_agent": "curl/7.79.1",
            "date_accessed": "2021-10-23T10:13:34Z"
        }
    ],
    "metadata": {
        "current_page": 1,
        "page_size": 20,
        "first_page": 1,
        "last_page": 1,
        "total_records": 1
    }
}

```

## Endpoints

| Method | URL | Handler | Action | Authentication |
|---|---|---|---|---|
| GET | `/v1/healthcheck` | `healthcheckHandler` | Get system status | False |
| GET | `/v1/links` | `getLinksHandler` | Retrieves all links | False |
| GET | `/v1/links/:hash` | `showLinkHandler` | Retrieves a single link | False |
| POST | `/v1/links` | `createLinkHandler` | Creates a new link | False |
| PATCH | `/v1/links/:hash` | `updateLinkHandler` | Updates a single link | False |
| DELETE | `/v1/links/:hash` | `deleteLinkHandler` | Deletes a single link | False |

[comment]: <> (| - | | | |)

[comment]: <> (| POST | `/v1/users` | `createUserHandler` | Creates a new user | False |)

[comment]: <> (| GET | `/v1/users/:id` | `showUserHandler` | Retrieves a single user | True |)

[comment]: <> (| PATCH | `/v1/users/:id` | `updateUserHandler` | Updates a single user | True |)

[comment]: <> (| DELETE | `/v1/users/:id` | `deleteUserHandler` | Deletes a single user | True |)


## Stack

- Gorilla mux webserver
- Sqlite (with Litestream)

## License

Licenced under Apache 2.0, see [LICENSE](/LICENSE).
