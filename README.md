# Simpl-GraphQL

Simpl-GraphQL is a practical GraphQL API written in Golang. It's designed to provide a robust and efficient way to create GraphQL APIs with features like CRUD operations, simple authentication, dynamic SQL query building from GraphQL queries, caching, and a complex audit trail.


## Getting Started

These instructions will guide you through the process of setting up the project on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed on your local machine:

- Go 1.16 or later
- PostgresSQL

### Installation

1. **Clone the repository**

   Start by cloning the repository to your local machine

   ```bash
   git clone https://github.com/danghieuthang/simpl-graphql
   ```

2. **Navigate to the project directory**

   ```bash
   cd simpl-graphql
   ```

3. **Install the dependencies**

   Use the `go get` command to install the necessary dependencies

   ```bash
   go get .
   ```

4. **Generate gqlgen**

   Run the gqlgen code generation tool with the `go generate` command

   ```bash
   go generate ./...
   ```

5. **Run the server**

   Finally, you can start the server using the `go run` command

   ```bash
   go run cmd/server.go
   ```

Now, you should have the server running locally on your machine. You can access it at `http://localhost:8080`.

## Technical Stack
- GraphQL with gqlgen
- Gorm
- PostgreSQL
- Log in with go.uber.org/zap
- Manage the environment with godotenv
- Generic Repository Pattern
- Performance tracking
- Simple Change Tracker
- Manage migration database with go-migrate
- - Manage migration database with go-migrate
    - See [migrate guide](https://github.com/danghieuthang/simpl-graphql/blob/main/cmd/migrations/README.md)

## Features
- CRUD with GraphQL
- Simple authentication
- Dynamically built SQL query from a GraphQL query
- Caching with memory caching and redis
- Complex audit trail
- Concurrency update control
## Project Organization
The project is organized into several directories and files. Here is a brief overview:

```
- assets
- cmd
    - assets
    - config
    - env
    - logs
    - migrations
- graph
    - model
    - resolver
    - schema
- internal
    - config
    - constant
        - app_error
        - enum
    - service
        - role
        - user
- pkg
    - audit
    - database
    - entity
    - file
    - logger
    - middleware
        - auth
    - repository
    - utils
- ui
    - .vscode
    - src
        - app
            - app-footer
            - app-header
            - graphql
            - user-detail
            - user-list
        - assets
```

## Contributing
We welcome contributions from the community. Please read our contribution guidelines before submitting a pull request.

## License
This project is licensed under the MIT License. See the `LICENSE` file in the project root for more details.

## Contact
If you have any questions or suggestions, feel free to reach out to us at dhthang1998@gmail.com