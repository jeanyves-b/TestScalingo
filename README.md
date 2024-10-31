# Response to scalingo

## Execution

```bash
docker compose up
```

Application will be then running on port `5000`.
you can then send request to access the data stored

## Available API calls

### heartbeat function

```bash
$ curl localhost:5000/ping
{ "status": "pong" }
```

### /getAll

Get method for every repo in the database

```bash
$ curl localhost:5000/getAll
{ "status": "OK/ERROR" ,
...
"way too long"
... }
```

### /getFilterd

Get method to get only some elements of the database

```bash
$ curl localhost:5000/getFiltered
{ "status": "OK/ERROR" ,
...
"still way too long "
... }
```

### /getAllUpdated

To force an update of the list and get everything, you can use /getAllUpdated

```bash
$ curl localhost:5000/getAllUpdated
{ "status": "OK/ERROR" ,
... }
```

### /getFilterdUpdated

To force an update of the list and get only some element, you can use /getFilteredUpdated

```bash
$ curl localhost:5000/getFilteredUpdated
{ "status": "OK/ERROR" ,
...}
```

### Test scripts

stresstest.sh : a small stresstest that launches the docker and try to launch as many request as it cansend to the server in 20 seconds and print the result at the end of the script.

## Architectural descisions

You can find the architecture documents in the architecture folder of this repo

### Using a struct containing a httpClient for the github call

As mentioned on the documentation in order to optimise performances, we can create a server that will be reused for every call we make to the same API. This is covered with the struct GitHubClient wich contains a http client, the url to call for the API and the last response given.

The API i called every 10 seconds to update the stored data or the update can be forced with a call to /getFilteredUpdated or /getAllUpdated. (The value of 10s is arbitrary and can be changed in gitHubAPICall.go line 24)

### Storing the github informations

In order to avoid making a DDOS on the github API we store the received data and update it regularily.

#### Storing the informations as a map of [string]interface{}

I chose to store the informations given by github inside a map of [string]interface{} to be the most efficient possible on getting the informations out of the structure and to ensure that the structure is as easy to maintain as possible, as scalable as possible, and adapts automaticaly to every possible changes in the github API.

It also provide a simple way to filter the returned information it we want to implement that in the futur

## What is still missing

- An automated test base
- Adapt the logger to give more context
- script a way to generate the README to include puml diagrams
