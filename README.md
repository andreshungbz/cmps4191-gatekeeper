# CMPS4191 Gatekeeper

## Running the Application

### Docker Compose

```
docker compose up
```

### Manual Method

#### Prerequisites

- curl
- go
- golang-migrate
- make
- PostgreSQL

#### Database Setup

```
CREATE ROLE gatekeeper_user WITH LOGIN PASSWORD 'gatekeeper_password';
CREATE DATABASE gatekeeper;
ALTER DATABASE gatekeeper OWNER TO gatekeeper_user;
```

#### Application Setup

```
cp .envrc.example .envrc
make db/migrations/up
make run
```
