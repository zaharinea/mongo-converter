Convert data dict from:
```
{
    "_id": {"$oid": "5cff9d32d2f8c149030bfb34"},
    "modified_at": {"$date": "2019-06-11T12:23:14.496Z"},
    "nested_dict": {"_id": {"$oid": "5e33d21491ca5b913a5c6df5"}},
}
```
to:
```
{
    "_id": "5cff9d32d2f8c149030bfb34",
    "modified_at": "2019-06-11T12:23:14.496Z",
    "nested_dict": {"_id": "5e33d21491ca5b913a5c6df5"},
}
```

## example usage:
```
mongoexport "mongodb://localhost:27017" -d database -c collection --type=json -o mongo_dump.json
./mongo-converter -in=mongo_dump.json -unwrap-numbers=0 > mongo_dump.converted.json
```
or
```
mongoexport "mongodb://localhost:27017" -d database -c collection --type=json | ./mongo-converter -unwrap-numbers=0 > mongo_dump.converted.json
```
