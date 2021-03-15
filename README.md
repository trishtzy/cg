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

Once it's ready, go into es01 bash.

```sh
$ docker exec -it country-es01 bash
$ vi config/elasticsearch.yml # or use your preferred editor
```

### API

- GET localhost:8080/countries

- Open "http://localhost:8081/" in browser to view database UI

- localhost:9200 for Elasticsearch API
