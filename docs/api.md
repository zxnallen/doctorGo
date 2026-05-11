# API

## System

- `GET /`
- `GET /health`

## Auth

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

## Admin Auth

- `POST /api/v1/admin/auth/login`
- `POST /api/v1/admin/auth/refresh`
- `POST /api/v1/admin/auth/logout`

## News

- `GET /api/v1/news/categories`
- `GET /api/v1/news?page=1&size=10`
- `GET /api/v1/news/:id`

Public news endpoints use Redis cache when Redis is available.

## Admin News

- `GET /api/v1/admin/news`
- `POST /api/v1/admin/news`
- `GET /api/v1/admin/news/:id`
- `PUT /api/v1/admin/news/:id`
- `PATCH /api/v1/admin/news/:id/status`
- `DELETE /api/v1/admin/news/:id`

## Admin News Categories

- `GET /api/v1/admin/news-categories`
- `POST /api/v1/admin/news-categories`
- `PUT /api/v1/admin/news-categories/:id`
- `PATCH /api/v1/admin/news-categories/:id/status`
- `DELETE /api/v1/admin/news-categories/:id`

## Admin Operation Logs

- `GET /api/v1/admin/operation-logs?page=1&size=20&admin_id=1&action=create&resource=news`

Query fields: `page`, `size`, `admin_id`, `action`, `resource`.

## Admin RBAC

Protected admin endpoints require both admin JWT and matching permission code.

- `GET /api/v1/admin/rbac/roles`
- `POST /api/v1/admin/rbac/roles`
- `PUT /api/v1/admin/rbac/roles/:id`
- `PUT /api/v1/admin/rbac/roles/:id/permissions`
- `GET /api/v1/admin/rbac/permissions`
- `PUT /api/v1/admin/rbac/admins/:admin_id/roles`

## Upload

- `POST /api/v1/upload/signed-url`
- Header: `Authorization: Bearer <token>`
