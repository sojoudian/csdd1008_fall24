# Save IP to MongoDB
```bash
git pull
git clone https://github.com/sojoudian/csdd1008_fall24.git


go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options


$ curl 'https://api.ipify.org?format=json'
docker run -d --name mongodb -p 27017:27017 -v mongodb:/data/db mongo:latest
```
