curl -XPOST -d  '{"name": "new"}' http://localhost:8081/users
curl -XGET http://localhost:8081/users
curl -XGET http://localhost:8081/users/1/relationships
curl -XPUT -d '{"state": "disliked"}' http://localhost:8081/users/2/relationships/1
curl -XPUT -d '{"state": "liked"}' http://localhost:8081/users/2/relationships/1
