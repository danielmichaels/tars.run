# Short Links

A simple short link application. This is exists solely for learning and personal purposes only.

## Features

- Gorilla mux webserver
- Sqlite backend
- Go templates

## Installation

The prod version runs in a single docker container

In development, the application can be started using `air`. Simply run the `air` command and a 
hot-reloading webserver will be started.

## How it works

When provided a URL, the server will return a shortened URL. Anyone who clicks on the shortened 
URL will be redirected to the original URL. 

To make this possible, a mapping must be created between the original URL and its shortened 
sibling. This application holds the mapping within a Sqlite database.

Each time a short link is clicked, that event is logged. A summary page with this data is 
available for each short link that is entered into the system.

## Endpoints

| Method | URL | Handler | Action | Authentication |
|---|---|---|---|---|
| GET | `/api/v1/healthcheck` | `healthcheckHandler` | Get system status | False |
| GET | `/api/v1/links` | `getLinksHandler` | Retrieves all links | True |
| GET | `/api/v1/links/:hash` | `showLinkHandler` | Retrieves a single link | False |
| POST | `/api/v1/links` | `createLinkHandler` | Creates a new link | False |
| PATCH | `/api/v1/links/:hash` | `updateLinkHandler` | Updates a single link | True |
| DELETE | `/api/v1/links/:hash` | `deleteLinkHandler` | Deletes a single link | True |
| - | | | |
| POST | `/api/v1/users` | `createUserHandler` | Creates a new user | False |
| GET | `/api/v1/users/:id` | `showUserHandler` | Retrieves a single user | True |
| PATCH | `/api/v1/users/:id` | `updateUserHandler` | Updates a single user | True |
| DELETE | `/api/v1/users/:id` | `deleteUserHandler` | Deletes a single user | True |

