# Go Restful API Boilerplate

[![GoDoc Badge]][godoc] [![GoReportCard Badge]][goreportcard]

Easily extendible RESTful API boilerplate aiming to follow idiomatic go and best practice.

The goal of this boiler is to have a solid and structured foundation to build upon on.

Any feedback and pull requests are welcome and highly appreciated. Feel free to open issues just for comments and discussions.

## Features

The following feature set is a minimal selection of typical Web API requirements:

- Configuration using [viper](https://github.com/spf13/viper)
- CLI features using [cobra](https://github.com/spf13/cobra)
- PostgreSQL support including migrations using [go-pg](https://github.com/go-pg/pg)
- Structured logging with [Logrus](https://github.com/sirupsen/logrus)
- Routing with [chi router](https://github.com/go-chi/chi) and middleware
- JWT Authentication using [lestrrat-go/jwx](https://github.com/lestrrat-go/jwx) with example passwordless email authentication
- Request data validation using [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
- HTML emails with [go-mail](https://github.com/go-mail/mail)

## Start Application

- Clone this repository
- Create a postgres database and set environment variables for your database accordingly if not using same as default
- Run the application to see available commands: `go run main.go`
- First initialize the database running all migrations found in ./database/migrate at once with command _migrate_: `go run main.go migrate`
- Run the application with command _serve_: `go run main.go serve`

Or just use the provided docker-compose file. After first start attach to the server container and run `./main migrate` to populate the database.

## API Routes

### Authentication

For passwordless login following routes are available:

| Path          | Method | Required JSON | Header                                | Description                                                             |
| ------------- | ------ | ------------- | ------------------------------------- | ----------------------------------------------------------------------- |
| /auth/login   | POST   | email         |                                       | the email you want to login with (see below)                            |
| /auth/token   | POST   | token         |                                       | the token you received via email (or printed to stdout if smtp not set) |
| /auth/refresh | POST   |               | Authorization: "Bearer refresh_token" | refresh JWTs                                                            |
| /auth/logout  | POST   |               | Authorizaiton: "Bearer refresh_token" | logout from this device                                                 |

### Example API

Besides /auth/_ the API provides two main routes /api/_ and /admin/\*, as an example to separate application and administration context. The latter requires to be logged in as administrator by providing the respective JWT in Authorization Header.

Check [routes.md](routes.md) for a generated overview of the provided API routes.

### Client API Access and CORS

The server is configured to serve a Progressive Web App (PWA) client from _./public_ folder (this repo only serves an example index.html, see below for a demo PWA client to put here). In this case enabling CORS is not required, because the client is served from the same host as the api.

If you want to access the api from a client that is serverd from a different host, including e.g. a development live reloading server with below demo client, you must enable CORS on the server first by setting environment variable _ENABLE_CORS=true_ on the server to allow api connections from clients serverd by other hosts.

#### Demo client application

For demonstration of the login and account management features this API serves a demo [Vue.js](https://vuejs.org) PWA. The client's source code can be found [here](https://github.com/dhax/go-base-vue). Build and put it into the api's _./public_ folder, or use the live development server (requires CORS enabled).

Outgoing emails containing the login token will be print to stdout if no valid email smtp settings are provided by environment variables (see table below). If _EMAIL_SMTP_HOST_ is set but the host can not be reached the application will exit immediately at start.

Use one of the following bootstrapped users for login:

- admin@boot.io (has access to admin panel)
- user@boot.io

A deployed version can also be found on [Heroku](https://govue.herokuapp.com)

### Environment Variables

By default viper will look at $HOME/.go-base.yaml for a config file. Setting your config as Environment Variables is recommended as by 12-Factor App.

| Name                    | Type          | Default                     | Description                                                                                                                                                                                               |
| ----------------------- | ------------- | --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| PORT                    | string        | localhost:3000              | http address (accepts also port number only for heroku compability)                                                                                                                                       |
| LOG_LEVEL               | string        | debug                       | log level                                                                                                                                                                                                 |
| LOG_TEXTLOGGING         | bool          | false                       | defaults to json logging                                                                                                                                                                                  |
| DB_NETWORK              | string        | tcp                         | database 'tcp' or 'unix' connection                                                                                                                                                                       |
| DB_ADDR                 | string        | localhost:5432              | database tcp address or unix socket                                                                                                                                                                       |
| DB_USER                 | string        | postgres                    | database user name                                                                                                                                                                                        |
| DB_PASSWORD             | string        | postgres                    | database user password                                                                                                                                                                                    |
| DB_DATABASE             | string        | postgres                    | database shema name                                                                                                                                                                                       |
| AUTH_LOGIN_URL          | string        | http://localhost:3000/login | client login url as sent in login token email                                                                                                                                                             |
| AUTH_LOGIN_TOKEN_LENGTH | int           | 8                           | length of login token                                                                                                                                                                                     |
| AUTH_LOGIN_TOKEN_EXPIRY | time.Duration | 11m                         | login token expiry                                                                                                                                                                                        |
| AUTH_JWT_SECRET         | string        | random                      | jwt sign and verify key - value "random" creates random 32 char secret at startup (and automatically invalidates existing tokens on app restarts, so during dev you might want to set a fixed value here) |
| AUTH_JWT_EXPIRY         | time.Duration | 15m                         | jwt access token expiry                                                                                                                                                                                   |
| AUTH_JWT_REFRESH_EXPIRY | time.Duration | 1h                          | jwt refresh token expiry                                                                                                                                                                                  |
| EMAIL_SMTP_HOST         | string        |                             | email smtp host (if set and connection can't be established then app exits)                                                                                                                               |
| EMAIL_SMTP_PORT         | int           |                             | email smtp port                                                                                                                                                                                           |
| EMAIL_SMTP_USER         | string        |                             | email smtp username                                                                                                                                                                                       |
| EMAIL_SMTP_PASSWORD     | string        |                             | email smtp password                                                                                                                                                                                       |
| EMAIL_FROM_ADDRESS      | string        |                             | from address used in sending emails                                                                                                                                                                       |
| EMAIL_FROM_NAME         | string        |                             | from name used in sending emails                                                                                                                                                                          |
| ENABLE_CORS             | bool          | false                       | enable CORS requests                                                                                                                                                                                      |

### Testing

Package auth/pwdless contains example api tests using a mocked database.

[godoc]: https://godoc.org/github.com/dhax/go-base
[godoc badge]: https://godoc.org/github.com/dhax/go-base?status.svg
[goreportcard]: https://goreportcard.com/report/github.com/dhax/go-base
[goreportcard badge]: https://goreportcard.com/badge/github.com/dhax/go-base
