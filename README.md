# RealWorld API - Gin + GORM Implementation

This is a backend implementation of the [RealWorld](https://github.com/gothinkster/realworld) API spec using **Go (Gin)** and **GORM**.

## Project Structure

```
.
├── main.go
├── database/
│   └── mysql.go
├── models/
│   └── article.go
│   └── comment.go
│   └── tag.go
│   └── user.go
│   └── favorite.go
│   └── follow.go
├── repositories/
│   └── article_repository.go
│   └── comment_repository.go
│   └── tag_repository.go
│   └── favorite_repository.go
│   └── follow_repository.go
│   └── user_repository.go
├── services/
│   └── article_service.go
│   └── comment_service.go
│   └── tag_service.go
│   └── user_service.go
│   └── profile_service.go
├── handlers/
│   └── article_handler.go
│   └── comment_handler.go
│   └── tag_handler.go
│   └── profile_handler.go
│   └── user_handler.go
├── middlewares/
│   └── auth_middleware.go
├── utils/
│   └── utils.go
├── dto/
│   └── article.go
└── go.mod
```

## Setup

```bash
go mod tidy
go run main.go
```

> ⚠️ Configure your database in `database/db.go`.

## API Endpoints

### Authentication

| Method | Endpoint           | Description      |
| ------ | ------------------ | ---------------- |
| POST   | `/api/users/login` | Login            |
| POST   | `/api/users`       | Register         |
| GET    | `/api/user`        | Get Current User |
| PUT    | `/api/user`        | Update User      |

---

### Articles

| Method | Endpoint              | Description        |
| ------ | --------------------- | ------------------ |
| POST   | `/api/articles`       | Create Article     |
| GET    | `/api/articles`       | Get Articles       |
| GET    | `/api/articles/:slug` | Get Single Article |
| PUT    | `/api/articles/:slug` | Update Article     |
| DELETE | `/api/articles/:slug` | Delete Article     |

#### Example Create Article Request:

```json
{
  "article": {
    "title": "How to train your dragon",
    "description": "Ever wonder how?",
    "body": "You have to believe",
    "tagList": ["reactjs", "angularjs", "dragons"]
  }
}
```

---

### Comments

| Method | Endpoint                                  | Description    |
| ------ | ----------------------------------------- | -------------- |
| POST   | `/api/articles/:slug/comments`            | Add Comment    |
| GET    | `/api/articles/:slug/comments`            | Get Comments   |
| DELETE | `/api/articles/:slug/comments/:commentId` | Delete Comment |

Example Comment Request:

```json
{
  "comment": {
    "body": "His name was my name too."
  }
}
```

---

### Favorite Article

| Method | Endpoint                       | Description           |
| ------ | ------------------------------ | --------------------- |
| POST   | `/api/articles/:slug/favorite` | Favorite an Article   |
| DELETE | `/api/articles/:slug/favorite` | Unfavorite an Article |

---

### Tags

| Method | Endpoint    | Description      |
| ------ | ----------- | ---------------- |
| GET    | `/api/tags` | Get List of Tags |

Example Response:

```json
{
  "tags": ["reactjs", "angularjs", "dragons"]
}
```

---

## Authentication

All protected routes require a token in the `Authorization` header:

```
Authorization: Token <your_token>
```

---

## References

* [RealWorld API Spec](https://realworld-docs.netlify.app/)
* [Gin Web Framework](https://gin-gonic.com/)
* [GORM ORM](https://gorm.io/)
