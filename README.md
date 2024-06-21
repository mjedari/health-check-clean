## API Health Checker

### Introduction

The whole structure is based on clean architecture packaging is by technology. As a matter of Clean Code & Architecture,
OOP, SOLID principles there are many situations in code that can be revised but for the sake of time limitation deferred
to future. All repositories and cache mechanism decoupled to make the service scalable on more than one instance in the
future.

### Requirements

MySQL is required for storing endpoints. But for the caching tasks you can either use both `redis` or `memory`. For more than one
instance we recommend using `redis` which is configurable through config file.

### Migrations

```bash
make migration
```

### Quick Start

```bash
make start
```

### APIs

#### Create endpoint

```http request
POST http://localhost:8080/endpoint/
```

Validate all inputs and store it in repository.

```json lines
// Request Body
{
  "url": "http://localhost:8000",
  "method": "get",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": null,
  "interval": 3
}
```

#### All endpoints

```http request
GET http://localhost:8080/endpoint/
```

```json lines
// Response
[
  {
    "id": 3,
    "url": "http://localhost:8000",
    "method": "get",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": null,
    "interval": 3,
    "created_at": "2024-06-18T11:55:26.404Z",
    "updated_at": "2024-06-18T11:55:26.404Z"
  }
]
```

#### Delete endpoint

Deletes endpoint from repository and stops its task if has been started already.

```http request
DELETE http://localhost:8080/endpoint/{endpoint_id}
```

#### Start Watching

Start the task of watching the endpoint.

```http request
GET http://localhost:8080/endpoint/{endpoint_id}/start
```

```json lines
// Success response
{
  "message": "watching started!"
}
```

#### Stop Watching

Stop watching endpoint and remove the task from cache.

```http request
GET http://localhost:8080/endpoint/{endpoint_id}/stop
```

```json lines
// Success response
{
  "message": "watching stopped!"
}
```

### Test

```bash
make test
```

### Components

I decided to add some extra components to deliver high quality application. These components are:

#### Retry Pattern

On connecting to the redis for instance, application seeks the connection process for some times fulfill the job. But
this could be used also in database queries.
