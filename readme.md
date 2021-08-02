# Event Tracker Challenge

## Project Description
Event tracking microservice that allows for events to be recorded, and then retrieved in various ways. 

## Endpoints

### User (/api/v1 subroute)
These endpoints are intended for user usage. The admin accounts are also authorized to use these.

#### POST
- /events/{name}
    - Allows the user to create a new event occurrence. The name of the event is a parameter (*name*) in the URL. If the event was already registered before, a new post registers new occurrences.
    - The request body (in JSON format) can include the following parameters:
      - "count": the event occurrences count.
      - "date": the date and hour in which those occurrences happened, must be in the format "YYYY-MM-DD HH:mm:ss".
    - Example:  POST {base_url}/api/v1/events/login1 with an empty body: creates a single 'login1' event occurrence, at the current time.

#### GET
- /events
  - Returns the total list of registered events, including the count and date of occurrence, summing up the count by dates.
    - Optional query parameters:
      - "start_date" and "end_date": These determine a date range for the results, must be in the format "YYYY-MM-DD".
- /event_history
  - Returns a history of all the registered events, and the total count for each one.
- /event_frequencies/{name}/hist
  - Returns a png image with a histogram showing the distribution of a given event (the *name* parameter in the URL) in the database, along the 24 hours of a day.

### Admin (/admin/v1 subroute)
These endpoints are intended for admin usage, and involve more _dangerous_ operations.

#### GET
- /events/{name}
  - Returns all the recorded occurrences of a given event (the *name* parameter in the URL), summing up the count by dates.
- /event_frequencies/{name}
  - Returns the total count of occurrences of a given event (the *name* parameter in the URL) and its hourly distribution.
- /event_frequencies
  - Returns the total count of occurrences of all the registered events and their hourly distribution.

#### DELETE
- /events/{name}
  - Deletes all the occurrences of a given event (the *name* parameter in the URL).

### Health (/health subroute)

- /ping
  - Just a simple ping check.

## Authorization 

Minimal API Key Authorization is required to use the API. There are two levels of authorization: user and admin.
The API key is sent via a 'x-api-key' header.

## Database

The database used is a SQLite3 database.