DATABASE_NAME = gh-journal

.PHONY: db.dump
db.dump:
	rm -rf dump
	docker exec -ti ghjournal_mongo_1 mongodump -d $(DATABASE_NAME) -o /data/db/dump

.PHONY: db.restore
db.restore:
	docker exec -ti ghjournal_mongo_1 mongorestore -d $(DATABASE_NAME) --drop /data/db/dump/gh-journal
