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
