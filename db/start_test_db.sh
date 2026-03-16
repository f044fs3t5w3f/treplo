#!/bin/sh
docker run --name treplo-db -p 5432:5432 -e POSTGRES_USER=treplo -e POSTGRES_DB=treplo -e POSTGRES_PASSWORD=1234  -d postgres 