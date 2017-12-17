# Workshop - backend

A collection of endpoints to power ~~ workshop ~~

## Endpoints

### Workshops

`GET /workshops`

Returns a json list of workshop events scheduled for after the current UTC date.

`POST /workshops`

Accepts a json workshop struct to add a workshop to the db

`PUT /workshops`

Updates an existing workshop

### Events

`GET /events`

Returns a json list of events scheduled for after the current UTC date.

`POST /events`

Accepts a json event struct to add a workshop to the db

`PUT /events`

Updates an existing workshop


