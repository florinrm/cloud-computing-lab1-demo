- cream imaginea folosind Dockerfile-ul din folderul curent:
```bash
docker build -t api-laborator-1-image .
```
- cream reteaua bridge:
```bash
docker network create -d bridge laborator1-db-network
```
- vedem reteaua creata:
```bash
docker network ls
```
- cream volumul:
```bash
docker volume create laborator1-db-persistent-volume
```
- vedem ca s-a creat volumul:
```bash
docker volume ls
```
- pornim containerul de baze de date:
```bash
docker run -v "$PWD"/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql -v laborator1-db-persistent-volume:/var/lib/postgresql/data --network=laborator1-db-network -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=books --name laborator1-db postgres
```
- pornim containerul de API (server):
```bash
docker run --network=laborator1-db-network -e PGUSER=admin -e PGPASSWORD=admin -e PGDATABASE=books -e PGHOST=laborator1-db -e PGPORT=5432 --name laborator1-api -p 5555:80 api-laborator-1-image
```
- serviciul ruleaza la urmatoarea adresa: http://localhost:5555/app/books