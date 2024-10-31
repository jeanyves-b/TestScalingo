# Response to scalingo

## Execution

```
docker compose up
```

Application will be then running on port `5000`
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

Get to get only some element -- work in progress
```
$ curl localhost:5000/getFiltered
{ ...
"still way too long "
... }
```