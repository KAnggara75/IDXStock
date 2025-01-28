#!/bin/bash

search_dir=csv2sql
filename=

abort() {
	echo "$@"
	exit 1
}

copy_file() {
	cp ./dataset/Saham/Semua/*.csv ./csv2sql/
}

read_all_file() {
	for file in "$search_dir"/*.csv; do
		filename=$(echo "$file" | sed "s/\.\///g")
		old_filename=$(echo "$filename" | sed "s/$search_dir\///g")
		# shellcheck disable=SC2001
		new_filename=$(echo "$old_filename" | sed "s/\.csv//g")
		echo "Formating, $old_filename to $new_filename.sql"
		# Replace blank with NULL
		sed -i='' -e 's/,,/,NULL,/g' "$filename"
		# Remove T00:00:00 if exist
		sed -i='' -e "s/T00:00:00//g" "$filename"
		# add quote 'YYYY-MM-DD'
		sed -i='' -e 's/\([0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}\)/'\''\1'\''/' "$filename"
		# add (
		sed -i='' -e "1!s/^/('$new_filename',/" "$filename"
		# add ),
		sed -i='' -e '1!s/$/),/' "$filename"
		# Remove coma symbol on last line
		sed -i='' -e '$s/,$//' "$filename"
		# add INSERT (column... in first Line
		sed -i='' -e '1s/^/INSERT INTO history (code, /' "$filename"
		# add ) VALUES in first Line
		sed -i='' -e '1s/$/) VALUES/' "$filename"
		# fix change to `change`
		# shellcheck disable=SC2016
		sed -i='' -e 's/change/`change`/g' "$filename"

		mv $search_dir/"$old_filename" ./sql_dump/"$new_filename".sql
	done
}

main() {
	clear
	pwd
	mkdir -p csv2sql
	rm -rf sql_dump
	mkdir -p sql_dump

	copy_file
	read_all_file

	rm -rf csv2sql
}

main || abort "Compose Error!"
