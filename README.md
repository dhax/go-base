
## Go Restful API Boilerplate

Easily extendible RESTful API boilerplate aiming to follow idiomatic go and best practice.

### Features
* PostgreSQL support including migrations using [go-pg](https://github.com/go-pg/pg)
* Structured logging with [Logrus](https://github.com/sirupsen/logrus)
* Routing with [chi router](https://github.com/go-chi/chi)
* JWT Authentication using [jwt-go](https://github.com/dgrijalva/jwt-go) with passwordless email authentication (could be easily extended to use passwords instead)
* Configuration using [viper](https://github.com/spf13/viper)
* CLI features using [cobra](https://github.com/spf13/cobra)
* [dep](https://github.com/golang/dep) for dependency management



### Environment Variables

Name | Type | Default | Description
---|---|---|---
PORT | int | 3000 | http port
LOG_LEVEL | string | debug | log level
LOG_TEXTLOGGING | bool | false | defaults to json logging
DATABASE_URL | string | postgres://postgres:postgres@localhost:5432/gobase?sslmode=disable | PostgreSQL connection string
AUTH_LOGIN_URL | string | http://localhost:3000/login | client login url as sent in login token email
AUTH_LOGIN_TOKEN_LENGTH | int | 8 | length of login token
AUTH_LOGIN_TOKEN_EXPIRY | int | 11 | login token expiry in minutes
AUTH_JWT_SECRET | string | random | jwt sign and verify key - value "random" sets random 32 char secret at startup
AUTH_JWT_EXPIRY | int | 15 | jwt access token expiry in minutes
AUTH_JWT_REFRESH_EXPIRY | int | 60 | jwt refresh token expiry in minutes
EMAIL_SMTP_HOST | string || email smtp host
EMAIL_SMTP_PORT | int || email smtp port
EMAIL_SMTP_USER | string || email smtp username
EMAIL_SMTP_PASSWORD | string || email smtp password
EMAIL_FROM_ADDRESS | string || from address used in sending emails
EMAIL_FROM_NAME | string || from name used in sending emails

