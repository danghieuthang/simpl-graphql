# prac-golang


This project employs a practical API written in the Golang programming language.


# Technical Stack
- GraphQL with graphql-go
- Gorm
- PostgreSQL
- Log in with go.uber.org/zap
- Manage the environment with godotenv
- Simple Repository Pattern
- Performance tracking


# Features
- CRUD with GraphQL
- Simple authentication
- Try to implement a dynamically built SQL query from a GraphQL query.
- Caching with memory caching and redis
- Build a complex audit trail


# Project organization
```
├───.vscode
├───controller
│   └───user
├───database
├───entity
├───gql
├───initializers
├───pkg
│   ├───jwt
│   └───log
├───repository
│   ├───rolerepo
│   └───userrepo
└───utils
```