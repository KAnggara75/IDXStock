#!/bin/bash

search_dir=sql_dump
filename=


abort() {
	echo "$@"
	exit 1
}

connect_to_db() {
	for file in "$search_dir"/*.sql; do
		filename=$(echo "$file" | sed "s/\.\///g")
		filename=$(echo "$filename" | sed "s/sql_dump\///g")
		echo "Execute, $filename"
		mysql -h 127.0.0.1 -u user -p'password' -P 3306 database <"$(pwd)"/$search_dir/"$filename" >>insert.log
		echo "Done, $filename"
		sleep 1s
	done

}

main() {
	clear
	connect_to_db
}

main || abort "Connection Error!"
