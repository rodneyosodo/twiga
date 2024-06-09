# API Design

The system will be designed as a collection of independent microservices that communicate via APIs. gRPC is used for communication between services for better performance and type safety while REST APIs are exposed to external clients. WebSockets on the other hand are used for real-time notifications.

The following microservices will be implemented:

## 1. User Service

Exposes APIs for user registration, login, profile management, and preference updates

API Endpoints:

- `POST /users`: Create a new user account.
- `GET /users`: Retrieve user profile.
- `PATCH /users/password`: Update user password.
- `GET /users/{id}`: Retrieve user profile by ID.
- `PUT /users/{id}`: Update user profile.
  - `PATCH /users/{id}/username`: Update user username.
  - `PATCH /users/{id}/email`: Update user email.
  - `PATCH /users/{id}/bio`: Update user bio.
  - `PATCH /users/{id}/picture`: Update user profile picture.
  - `PATCH /users/{id}/preferences`: Update user preferences.
- `DELETE /users/{id}`: Delete user account.

- `POST /token/issue`: Issue a new access token.
- `POST /token/refresh`: Refresh an access token.

- `POST /users/preferences`: Set user notification preferences.
- `GET /users/preferences`: Retrieve user notification preferences.
- `PUT /users/preferences`: Update user notification preferences.
  - `PATCH /users/preferences/email`: Update email notification preference.
  - `PATCH /users/preferences/push`: Update push notification preference.
- `GET /users/{id}/followers`: Retrieve user followers.
- `GET /users/{id}/following`: Retrieve users followed by the user.

- `POST /users/{id}/follow`: Follow a user.
- `DELETE /users/{id}/unfollow`: Unfollow a user.

- `GET /users/feed`: Retrieve posts from users followed by the user.

- `GET /version`: Retrieve service version.
- `GET /metrics`: Retrieve service metrics.

gRPC Services:

- `GetUserByID`: Retrieve user profile by ID.
- `GetUserPreferences`: Retrieve user notification preferences.
- `GetUserFollowers`: Retrieve user followers.
- `CreateFeed`: Create a feed for a user.
- `IdentifyUser`: Identify user based on session token.

## 2. Post Service

Exposes APIs for creating posts, fetching user posts, and retrieving post details and comments.

API Endpoints:

- `POST /posts`: Create a new post.
- `GET /posts`: Retrieve posts.
- `GET /posts/{id}`: Retrieve a post by ID.
- `PUT /posts/{id}`: Update a post.
  - `PATCH /posts/{id}/content`: Update post content.
  - `PATCH /posts/{id}/tags`: Update post tags.
  - `PATCH /posts/{id}/image`: Update post image.
  - `PATCH /posts/{id}/tags`: Update post tags.
  - `PATCH /posts/{id}/visibility`: Update post visibility.
- `DELETE /posts/{id}`: Delete a post.

- `POST /posts/{id}/comments`: Add a comment to a post.
- `GET /posts/{id}/comments`: Retrieve comments on a post.
- `GET /posts/comments/{comment_id}`: Retrieve a comment by ID.
- `PUT /posts/comments/{comment_id}`: Update a comment.
- `DELETE /posts/comments/{comment_id}`: Delete a comment.

- `POST /posts/{id}/like`: Like a post.
- `GET /posts/{id}/likes`: Retrieve likes on a post.
- `DELETE /posts/{id}/unlike`: Unlike a post.

- `POST /posts/{id}/share`: Share a post.
- `GET /posts/{id}/shares`: Retrieve shares of a post.
- `DELETE /posts/{id}/unshare`: Unshare a post.

- `GET /version`: Retrieve service version.
- `GET /metrics`: Retrieve service metrics.

## 3. Notification Service

Consumes messages from message brokers published by User Service and Post Service. Uses web sockets for real-time delivery.

API Endpoints:

- `GET /notifications`: Retrieve user notifications.
- `GET /notifications/{id}`: Retrieve a notification by ID.
- `POST /notifications/{id}/read`: Mark a notification as read.
- `POST /notifications/read`: Mark all notifications as read.
- `DELETE /notifications/{id}`: Delete a notification.

- `GET /version`: Retrieve service version.
- `GET /metrics`: Retrieve service metrics.

WebSockets:

- `ws://ws`: Real-time notification delivery using WebSockets for a subscribed user.

Cron Jobs:

- `sendDailyDigest`: Send a daily email digest of notifications to users.
- `sendWeeklyDigest`: Send a weekly email digest of notifications to users.
- `sendMonthlyDigest`: Send a monthly email digest of notifications to users.
