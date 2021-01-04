#!/bin/bash
set -xe
mkdir -p /etc/timetable/resources

cp -r assets/ /etc/timetable/resources/

go install

habit-tracker

echo "1- mongo client to work with. 2- helper functions for habit-tracker 3- rest endpoints to seed data"
echo "4- rest endpoints for CRUD operations on db. 5- javascript/jquery changes for the UI. 6- client side changes for calling rest endpoints. 7- logic to show 7 days data"