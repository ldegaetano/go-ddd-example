# Golang ddd example

  The api was designed following a Domain Driven Design. The idea is to show a way to organize the code of an api. 

## Features and Roadmap
### v0.1

* Service as a transparent cache for retrive prices.
* Prices storage.
* Prices cache.
* Prices handler.
* Http server.
* Add docker.

## Notes

 Run tests with docker.
````
   docker-compose -f docker-test.yml build
   docker-compose -f docker-test.yml run test
````

Run the api with docker.
```
    docker-compose build
    docker-compose up
```

### Get Prices

Request: 
````
 curl --location --request GET 'localhost:8080/api/items/prices?items_codes=p1,p2'
````

Response :
- Status 200
`````
{
    "items": [
        {
            "item_code": "p1",
            "item_price": 5
        },
        {
            "item_code": "p2",
            "item_price": 4.5
        }
    ]
}
`````

- error:
````
{
    "code": 1,
    "message": "Items not found: p3."
}
````

### Set Price

Request: 
````
 curl --location --request POST 'localhost:8080/api/items/prices' \
    --header 'Content-Type: application/json' \
    --data-raw '{
	    "item_code": "p2",
	    "item_price": 4
    }'
````

Response :
- Status 204 No content

- error:
````
{
    "code": 0,
    "message": "Internal server error."
}
````