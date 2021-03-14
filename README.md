# Country Select

This is the server side for CountrySelect App.

## Pre-requsites

- Docker
- Go1.16

### To run

```sh
$ make setup

# Check all four containers are up
$ docker ps
# country-web - to get /countries
# country-es01 - to search countries
# country-pgweb - ui to view postgresql
# countrypg - postgres db
```

or

```sh
$ go build -o main
$ ./main
```
