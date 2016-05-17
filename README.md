## Blackflow Hub 
The service is a simple application store or application hub for Blackflow IOT project .
### Service configuration 

The service can be configured either by setting up environment variables or by specify parameters as command line 

ENV variables : 

* BFH_BIND_ADDR or -addr : Server bind address
* BFH_DB_CONN_STR or -db_conn : Mongo db connection string.Default = localhost. Example mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
* BFH_DB_NAME or -db_name : Database name
* BFH_FS_LOCATION or -fs_location : File store location
* BFH_JWT_SECRET or -jwt_secret : Jwt secret 

### Docker 
#### Build 
docker rmi alivinco/blackflowhub; docker build -t alivinco/blackflowhub .

#### Run
docker run --name blackflowhub -d -t -p 5050:5050 --link mongo -e BFH_DB_CONN_STR="mongo:27017" -e BFH_JWT_SECRET="_REPLACE_WITH_YOUR_SECRET"  alivinco/blackflowhub