# Linkrr API Documentation

Comprehensive guide to Linkrr’s REST API: authentication, URL shortening, redirects, analytics, and user endpoints.

---

## Table of Contents
- [Getting Started](#getting-started)
- [Auth](#auth)
  - [Sign Up](#sign-up)
  - [Log In](#log-in)
  - [Renew Access Token](#renew-access-token)
  - [Forgot Password](#forgot-password)
  - [Reset Password](#reset-password)
  - [Revoke Refresh Token](#revoke-refresh-token)
- [Shortener](#shortener)
  - [Create Short URL](#create-short-url)
  - [Create Alias](#create-alias)
- [Redirect](#redirect)
- [Analytics](#analytics)
  - [Global Analytics](#global-analytics)
  - [URL Analytics](#url-analytics)
  - [Alias Analytics](#alias-analytics)
- [Users](#users)
  - [Get My Links](#get-my-links)
- [Errors](#errors)
- [Notes](#notes)

---

## Getting Started
- Base URL: `http://localhost:8080` (adjust to your deployment)
- Authentication: Bearer token in `Authorization: Bearer <access_token>` unless otherwise stated.
- Content-Type: `application/json` for request bodies.

## Auth

### Sign Up
- Method: POST
- Path: `/auth/signup`
- Body:
```json
{
  "email": "user@example.com",
  "password": "your-strong-password",
  "name": "Jane Doe"
}
```
- Response: 201 Created
```json
{
  "message": "signup successful",
  "userId": "..."
}
```

### Log In
- Method: POST
- Path: `/auth/login`
- Body:
```json
{
  "email": "user@example.com",
  "password": "your-strong-password"
}
```
- Response:
```json
{
  "accessToken": "...",
  "refreshToken": "..."
}
```

### Renew Access Token
- Method: POST
- Path: `/auth/renew`
- Body:
```json
{
  "refreshToken": "..."
}
```
- Response:
```json
{
  "accessToken": "..."
}
```

### Forgot Password
- Method: POST
- Path: `/auth/forgot-password`
- Body:
```json
{
  "email": "user@example.com"
}
```
- Response: 200 OK `{ "message": "email sent" }`

### Reset Password
- Method: POST
- Path: `/auth/reset-password`
- Body:
```json
{
  "token": "reset-token",
  "newPassword": "new-strong-password"
}
```
- Response: 200 OK `{ "message": "password updated" }`

### Revoke Refresh Token
- Method: POST
- Path: `/auth/revoke`
- Body:
```json
{
  "refreshToken": "..."
}
```
- Response: 200 OK `{ "message": "revoked" }`

## Shortener

### Create Short URL
- Method: POST
- Path: `/shortener/url`
- Auth: Required
- Body:
```json
{
  "originalUrl": "https://example.com/very/long/path",
  "expiresAt": "2025-12-31T23:59:59Z" // optional
}
```
- Response:
```json
{
  "shortUrl": "https://lnk.rr/abc123",
  "id": "...",
  "alias": null
}
```

### Create Alias
- Method: POST
- Path: `/shortener/alias`
- Auth: Required
- Body:
```json
{
  "alias": "my-link",
  "targetUrl": "https://example.com/landing"
}
```
- Response:
```json
{
  "alias": "my-link",
  "shortUrl": "https://lnk.rr/my-link"
}
```

## Redirect
- Method: GET
- Path: `/:codeOrAlias`
- Description: Redirects to the original URL and records analytics.
- Response: 302 Found → location header to target URL.

## Analytics

### Global Analytics
- Method: GET
- Path: `/analytics/global`
- Auth: Required
- Query:
  - `from`: ISO timestamp
  - `to`: ISO timestamp
- Response:
```json
{
  "totalClicks": 1234,
  "uniqueVisitors": 987,
  "byCountry": { "US": 500, "NG": 120 },
  "byReferrer": { "google": 300, "twitter": 100 }
}
```

### URL Analytics
- Method: GET
- Path: `/analytics/url/:urlId`
- Auth: Required
- Query: `from`, `to` (optional)
- Response similar to Global with per-URL stats.

### Alias Analytics
- Method: GET
- Path: `/analytics/alias/:alias`
- Auth: Required
- Query: `from`, `to` (optional)
- Response similar to Global with per-alias stats.

## Users

### Get My Links
- Method: GET
- Path: `/users/me/links`
- Auth: Required
- Response:
```json
[
  {
    "id": "...",
    "shortUrl": "https://lnk.rr/abc123",
    "originalUrl": "https://example.com",
    "alias": null,
    "createdAt": "..."
  }
]
```

## Errors
- Standard JSON error format:
```json
{
  "error": "Bad Request",
  "message": "validation failed: ..."
}
```
- Common codes: 400, 401, 403, 404, 429, 500.

## Notes
- Rate limiting and analytics aggregation are handled internally.
- Redirects collect click data (IP, user agent, referrer) as configured.
- Actual paths and shapes may vary slightly with deployment; inspect `Routes.go` and `internal/handlers/*` for current routing.
