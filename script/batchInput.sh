#!/bin/bash

URL="http://127.0.0.1:3000/api/stocks/idx"
DATA_FOLDER="./stock_summary"
TOKEN=""

for file in "$DATA_FOLDER"/Stock*Summary-*.xlsx; do
    curl -o /dev/null -s -w "%{http_code} " --location $URL \
    --header "Authorization: Bearer $TOKEN" \
    --form "file=@$file"
    echo "-$file"
    sleep 5s
done
