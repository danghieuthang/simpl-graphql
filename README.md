# Simpl-GraphQL


This project employs a practical GraphQL API written in the Golang programming language.


# Technical Stack
- GraphQL with gqlgen
- Gorm
- PostgreSQL
- Log in with go.uber.org/zap
- Manage the environment with godotenv
- Generic Repository Pattern
- Performance tracking
- Simple Change Tracker
- Manage migration database with go-migrate
    - See [migrate guide](https://github.com/danghieuthang/simpl-graphql/blob/main/cmd/migrations/README.md)

# Features
- CRUD with GraphQL
- Simple authentication
- Try to implement a dynamically built SQL query from a GraphQL query.
- Caching with memory caching and redis
- Build a complex audit trail


# Project organization
```
├───assets
├───cmd
│   ├───assets
│   ├───config
│   ├───env
│   ├───logs
│   └───migrations
├───graph
│   ├───model
│   ├───resolver
│   └───schema
├───internal
│   ├───config
│   ├───constant
│   │   ├───app_error
│   │   └───enum
│   └───service
│       ├───role
│       └───user
├───pkg
│   ├───audit
│   ├───database
│   ├───entity
│   ├───file
│   ├───logger
│   ├───middleware
│   │   └───auth
│   ├───repository
│   └───utils
└───ui
    ├───.vscode
    └───src
        ├───app
        │   ├───app-footer
        │   ├───app-header
        │   ├───graphql
        │   ├───user-detail
        │   └───user-list
        └───assets
```

