FROM mongo

COPY migrate-data.json /migrate-data.json
CMD mongoimport --host mongodb --db hotel --collection account --type json --file /migrate-data.json --jsonArray