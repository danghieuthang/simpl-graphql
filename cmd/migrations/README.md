This migrate using golang-migrate https://github.com/golang-migrate/migrate

### Setup and installation
####Window: 
- scoop https://scoop.sh/
- Install golang migrate:
```
$ scoop install migrate
```

#### Linux:
```
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
```
#### Mac:
```
brew install golang-migrate
```

### Create migration
```
migrate create -ext sql -dir {Path to migrations} -seq {name of migration}
```
### Run migration
#### Up: 
```
migrate -path {Path to migrations} -database "{connection string}" -verbose up
```
Example:
```
migrate -path migration -database "postgresql://postgres:1@localhost:5432/golang?sslmode=disable" -verbose up
```
#### Rollback migrations: 
```
migrate -path {Path to migrations} -database "{connection string}" -verbose down {Number of commit to down}
```
Example:
```
migrate -path migration -database "postgresql://postgres:1@localhost:5432/golang?sslmode=disable" -verbose down
```
#### Forcing database version:
```
migrate -path {Path to migrations} -database "{connection string}" force {VERSION
```
Example:
```
migrate -path migration -database "postgresql://postgres:1@localhost:5432/golang?sslmode=disable" -verbose force 2
```