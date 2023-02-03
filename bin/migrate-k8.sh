#!/bin/bash

./bin/migrate -database postgres://postgres:@localhost:5432/app?sslmode=disable "$@"