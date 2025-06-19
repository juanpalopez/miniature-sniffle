# miniature-sniffle

This project provides a small REST API server written in Go 1.22+. It stores data in memory so no external database is required. Routes use the Go 1.22 `"<METHOD> <URI>"` registration syntax.

## Features

* Create anime entries with a unique title.
* List all anime or fetch a single anime by ID.
* Add ratings and text reviews to an anime.
* List reviews for a specific anime.

## Building and Running

``` 
# build the server
make build

# run the server on port 8080
make run
```

The API will be available at `http://localhost:8080`.

## Running Tests

```
make test
```

GitHub Actions are configured to run vet, build and test on each push.

## Example Requests

Create an anime:

```
curl -X POST http://localhost:8080/animes \
  -H "Content-Type: application/json" \
  -d '{"title":"Naruto","creator":"Masashi Kishimoto","genre":"Shonen","year":2002}'
```

Add a review:

```
curl -X POST http://localhost:8080/animes/1/reviews \
  -H "Content-Type: application/json" \
  -d '{"rating":5,"text":"Great show"}'
```
