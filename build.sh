#!/bin/bash

lsof -ti tcp:8999 | xargs kill > /dev/null
go build