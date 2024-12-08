# caching-proxy

Caching proxy server

Runs on specified port, forwards requests to specified origin and caches them

Cache is backed up to a file periodically

# Usage

```sh
caching-proxy --port <number> --origin <url>
```

Example:
```sh
caching-proxy --port 3000 --origin http://dummyjson.com
```

If the origin scheme is omitted, http:// is presumed:
```sh
caching-proxy --port 3000 --origin dummyjson.com
# origin is effectively http://dummyjson.com
```

## Additional flags

```sh
--debug - show debug logs
--backup - specify backup file for the cache
--clear-cache - clear cache (clears backup file before running the app)
```

# Build And Install

To build and install, use `go build` and `go install` respectively, from the project's root directory 
```sh
go build
```
```sh
go install
```

# Testing

```sh
go test ./...
```

# Roadmap reference

https://roadmap.sh/projects/caching-server