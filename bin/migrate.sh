#!/bin/bash

./bin/migrate -database postgres://app:password@localhost:5444/app?sslmode=disable "$@"