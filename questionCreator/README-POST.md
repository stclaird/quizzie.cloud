```bash
curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "kubernetes etcd",
    "numincorrectans" : 1,
    "numcorrectans" : 2,
    "numquestions" : 10,
    "category" : "kubernetes",
    "subcategory" : "etc"
}'

curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "kubernetes etcd",
    "numincorrectans" : 1,
    "numcorrectans" : 1,
    "numquestions" : 10,
    "category" : "kubernetes",
    "subcategory" : "etc"
}'

curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "gcp gke",
    "numquestions" : 50,
    "category" : "gcp",
    "subcategory" : "gke"
}'

curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "gcp gce",
    "numquestions" : 50,
    "category" : "gcp",
    "subcategory" : "gce"
}'

curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "gcp pub/sub",
    "numquestions" : 50,
    "category" : "gcp",
    "subcategory" : "pub/sub"
}'


curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "gcp gcs",
    "numquestions" : 50,
    "category" : "gcp",
    "subcategory" : "gcs"
}'

curl -L -X POST http://localhost:5001/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "gcp cloud armour",
    "numquestions" : 50,
    "category" : "gcp",
    "subcategory" : "cloud-armour"
}'

```
