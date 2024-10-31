# Response to scalingo

## Execution

```
docker compose up
```

Application will be then running on port `5000`.
you can then send request to access the data stored

## Test

heartbeat function
```
$ curl localhost:5000/ping
{ "status": "pong" }
```

Get for everything
```
$ curl localhost:5000/getAll
{ ...
"way too long"
... }
```

Get to get only some element
```
$ curl localhost:5000/getFiltered
{ ...
"still way too long "
... }
```

## What is still missing

1- An automated test base
2- Adapt the logger to give more context