#!/bin/bash

search_dir=sql_dump
filename=

abort() {
	echo $platform
	echo "$@"
	exit 1
}

read_all_file() {
	echo -e "INSERT INTO history (code, date,previous,open_price,first_trade,high,low,close,\`change\`,volume,value,frequency,index_individual,offer,offer_volume,bid,bid_volume,listed_shares,tradeble_shares,weight_for_index,foreign_sell,foreign_buy,delisting_date,non_regular_volume,non_regular_value,non_regular_frequency) VALUES" >IHSG.sql

	for file in "$search_dir"/*.sql; do
		filename=$(echo "$file" | sed "s/\.\///g")
		filename=$(echo "$filename" | sed "s/$search_dir\///g")
		echo "Collect from, $filename"
		sed -i='' -e '1d' $search_dir/$filename
		sed -i='' -e '$s/$/,/' $search_dir/$filename
		cat $search_dir/$filename >>IHSG.sql
		rm $search_dir/$filename=
	done
	sed -i='' -e '$s/,$//' IHSG.sql
	rm IHSG.sql=
}

main() {
	clear
	read_all_file
	git restore .
}

main || abort "Compose Error!"
