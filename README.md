# simple-kv

`simple-kv` is a simple database that stores its values as key-value pairs on disk.

This repository contains both the backend implementation, as well as a HTTP 
frontend to interact with the database through any HTTP client.

## Pre-requirements

1. A version of the [Go toolchain](https://go.dev/dl/) installed, at least version `1.19`
2. A version of [curl](https://curl.se/) installed

## Building and Running

1. To compile the server binary, run `go build .`
    1. A file called `simple-kv` should appear in the same directory
2. To run the server, type the following in the same directory as the `go build` command
    1. On MacOS/Linux run `./simple-kv`
    2. On Windows, run `simple-kv.exe`

The server should start with the following message:
```
INFO[0000] Started listening for requests on port 3000
```

## Interacting With the Server

To read and write keys to the server, the following endpoints are accessible at port `:3000`

* `PUT`, which takes in the arguments `key`, which corresponds to the unique key, and `value` which refers to the value that is associated with the key

For example, to insert the value `test` with the key `example`, you can write the following curl command

```bash
curl -X PUT -d key=example -d value=test http://localhost:3000
```

* `DELETE` just takes in a `key` and deletes that key from the database, along with its value

For example, to delete the value and key `example`, write with `curl`:

```bash
curl -X DELETE -d key=example http://localhost:3000
```

* `GET` which takes in a `key` and returns its value

For example, in curl, write

```bash
curl -X GET -d key=example http://localhost:3000
```
