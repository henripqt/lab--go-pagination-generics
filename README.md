# LAB - GO 1.18Beta1 TYPE PARAMETERS & PAGINATION

This repository is a test/showcase of how can the new go type parameters feature help us achieve a more scalable code when dealing with request for paginated resources in database.

## Dependencies

- go 1.18beta1
- Docker
- docker-compose
- go-migrate

## Setup

### 1. Setup database

```sh
docker-compose up -d && \
make migrate-up
```

## Running the app

### 1. No type parameters

```sh
make run
```

### 2. Type parameters branch

```sh
git checkout feat/implement_type_parameters
make run
```