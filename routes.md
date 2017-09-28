# github.com/dhax/go-base

GoBase REST API.

## Routes

<details>
<summary>`/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/***
	- _GET_
		- [SPAHandler.func1](/api/api.go#L101)

</details>
<details>
<summary>`/admin/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/admin/***
	- [RequiresRole.func1](/auth/authorizer.go#L11)
	- **/**
		- _GET_
			- [(*API).Router.func1](/api/admin/api.go#L42)

</details>
<details>
<summary>`/admin/*/accounts/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/admin/***
	- [RequiresRole.func1](/auth/authorizer.go#L11)
	- **/accounts/***
		- **/**
			- _GET_
				- [(*AccountResource).(github.com/dhax/go-base/api/admin.list)-fm](/api/admin/accounts.go#L50)
			- _POST_
				- [(*AccountResource).(github.com/dhax/go-base/api/admin.create)-fm](/api/admin/accounts.go#L51)

</details>
<details>
<summary>`/admin/*/accounts/*/{accountID}/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/admin/***
	- [RequiresRole.func1](/auth/authorizer.go#L11)
	- **/accounts/***
		- **/{accountID}/***
			- [(*AccountResource).(github.com/dhax/go-base/api/admin.accountCtx)-fm](/api/admin/accounts.go#L53)
			- **/**
				- _PUT_
					- [(*AccountResource).(github.com/dhax/go-base/api/admin.update)-fm](/api/admin/accounts.go#L55)
				- _DELETE_
					- [(*AccountResource).(github.com/dhax/go-base/api/admin.delete)-fm](/api/admin/accounts.go#L56)
				- _GET_
					- [(*AccountResource).(github.com/dhax/go-base/api/admin.get)-fm](/api/admin/accounts.go#L54)

</details>
<details>
<summary>`/api/*/account/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/api/***
	- **/account/***
		- [(*AccountResource).(github.com/dhax/go-base/api/app.accountCtx)-fm](/api/app/account.go#L48)
		- **/**
			- _PUT_
				- [(*AccountResource).(github.com/dhax/go-base/api/app.update)-fm](/api/app/account.go#L50)
			- _DELETE_
				- [(*AccountResource).(github.com/dhax/go-base/api/app.delete)-fm](/api/app/account.go#L51)
			- _GET_
				- [(*AccountResource).(github.com/dhax/go-base/api/app.get)-fm](/api/app/account.go#L49)

</details>
<details>
<summary>`/api/*/account/*/profile`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/api/***
	- **/account/***
		- [(*AccountResource).(github.com/dhax/go-base/api/app.accountCtx)-fm](/api/app/account.go#L48)
		- **/profile**
			- _PUT_
				- [(*AccountResource).(github.com/dhax/go-base/api/app.updateProfile)-fm](/api/app/account.go#L56)

</details>
<details>
<summary>`/api/*/account/*/token/{tokenID}/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/api/***
	- **/account/***
		- [(*AccountResource).(github.com/dhax/go-base/api/app.accountCtx)-fm](/api/app/account.go#L48)
		- **/token/{tokenID}/***
			- **/**
				- _PUT_
					- [(*AccountResource).(github.com/dhax/go-base/api/app.updateToken)-fm](/api/app/account.go#L53)
				- _DELETE_
					- [(*AccountResource).(github.com/dhax/go-base/api/app.deleteToken)-fm](/api/app/account.go#L54)

</details>
<details>
<summary>`/auth/*/login`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/auth/***
	- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
	- **/login**
		- _POST_
			- [(*Resource).(github.com/dhax/go-base/auth.login)-fm](/auth/api.go#L67)

</details>
<details>
<summary>`/auth/*/logout`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/auth/***
	- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
	- **/logout**
		- _POST_
			- [Verifier.func1](/vendor/github.com/go-chi/jwtauth/jwtauth.go#L70)
			- [AuthenticateRefreshJWT](/auth/authenticator.go#L66)
			- [(*Resource).(github.com/dhax/go-base/auth.logout)-fm](/auth/api.go#L73)

</details>
<details>
<summary>`/auth/*/refresh`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/auth/***
	- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
	- **/refresh**
		- _POST_
			- [Verifier.func1](/vendor/github.com/go-chi/jwtauth/jwtauth.go#L70)
			- [AuthenticateRefreshJWT](/auth/authenticator.go#L66)
			- [(*Resource).(github.com/dhax/go-base/auth.refresh)-fm](/auth/api.go#L72)

</details>
<details>
<summary>`/auth/*/token`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/auth/***
	- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
	- **/token**
		- _POST_
			- [(*Resource).(github.com/dhax/go-base/auth.token)-fm](/auth/api.go#L68)

</details>
<details>
<summary>`/ping`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [DefaultCompress](/vendor/github.com/go-chi/chi/middleware/compress.go#L38)
- [Timeout.func1](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [RequestLogger.func1](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType.func1](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/ping**
	- _GET_
		- [NewAPI.func2](/api/api.go#L73)

</details>

Total # of routes: 12
