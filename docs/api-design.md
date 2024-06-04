# API Design

The system will be designed as a collection of independent microservices that communicate via APIs. gRPC is used for communication between services for better performance and type safety while REST APIs are exposed to external clients. WebSockets on the other hand are used for real-time notifications.

The following microservices will be implemented:

## 1. User Service

Exposes APIs for user registration, login, profile management, and preference updates

API Endpoints:

- `POST /users`: Create a new user account.
- `GET /users/{id}`: Retrieve user profile by ID.
- `PUT /users/{id}`: Update user profile.
  - `PATCH /users/{id}/username`: Update user username.
  - `PATCH /users/{id}/password`: Update user password.
  - `PATCH /users/{id}/email`: Update user email.
  - `PATCH /users/{id}/bio`: Update user bio.
  - `PATCH /users/{id}/picture`: Update user profile picture.
  - `PATCH /users/{id}/preferences`: Update user preferences.
- `DELETE /users/{id}`: Delete user account.

- `POST /login`: Authenticate user and generate a JWT token.

- `POST /users/{id}/preferences`: Set user notification preferences.
- `GET /users/{id}/preferences`: Retrieve user notification preferences.
- `PUT /users/{id}/preferences`: Update user notification preferences.
  - `PATCH /users/{id}/preferences/email`: Update email notification preference.
  - `PATCH /users/{id}/preferences/push`: Update push notification preference.
- `DELETE /users/{id}/preferences`: Delete user notification preferences.

- `GET /users/{id}/followers`: Retrieve user followers.
- `GET /users/{id}/following`: Retrieve users followed by the user.

- `POST /users/{id}/follow`: Follow a user.
- `POST /users/{id}/unfollow`: Unfollow a user.

- `GET /users/{id}/feed`: Retrieve posts from users followed by the user.

gRPC Services:

- `GetUserByID`: Retrieve user profile by ID.
- `GetUserPreferences`: Retrieve user notification preferences.
- `UpdateFeed`: Update user feed with new posts.
- `IdentifyUser`: Identify user based on session token.

## 2. Post Service

Exposes APIs for creating posts, fetching user posts, and retrieving post details and comments.

API Endpoints:

- `POST /posts`: Create a new post.
- `GET /posts/{id}`: Retrieve a post by ID.
- `PUT /posts/{id}`: Update a post.
  - `PATCH /posts/{id}/content`: Update post content.
  - `PATCH /posts/{id}/image`: Update post image.
  - `PATCH /posts/{id}/tags`: Update post tags.
  - `PATCH /posts/{id}/visibility`: Update post visibility.
- `DELETE /posts/{id}`: Delete a post.

- `POST /posts/{id}/comments`: Add a comment to a post.
- `GET /posts/{id}/comments`: Retrieve comments on a post.
- `PUT /posts/{id}/comments/{comment_id}`: Update a comment.
- `DELETE /posts/{id}/comments/{comment_id}`: Delete a comment.

- `POST /posts/{id}/likes`: Like a post.
- `GET /posts/{id}/likes`: Retrieve likes on a post.
- `DELETE /posts/{id}/likes`: Unlike a post.

- `POST /posts/{id}/shares`: Share a post.
- `GET /posts/{id}/shares`: Retrieve shares of a post.
- `DELETE /posts/{id}/shares`: Unshare a post.

- `GET /users/{id}/posts`: Retrieve user posts.
- `GET /users/{id}/posts/{post_id}`: Retrieve a user post by ID.

## 3. Notification Service

Consumes messages from message brokers published by User Service and Post Service. Uses web sockets for real-time delivery.

API Endpoints:

- `GET /notifications`: Retrieve user notifications.
- `DELETE /notifications/{id}`: Delete a notification.
- `DELETE /notifications`: Delete all notifications.
- `POST /notifications/{id}/read`: Mark a notification as read.
- `POST /notifications/read-all`: Mark all notifications as read.

- `POST /notifications/settings`: Create notification settings.
- `GET /notifications/settings`: Retrieve notification settings.
- `PUT /notifications/settings`: Update notification settings.
  - `POST /notifications/settings/email`: Enable email notifications.
  - `POST /notifications/settings/push`: Enable push notifications.
  - `DELETE /notifications/settings/email`: Disable email notifications.
  - `DELETE /notifications/settings/push`: Disable push notifications.
- `DELETE /notifications/settings`: Delete notification settings.

WebSockets:

- `ws://notifications`: Real-time notification delivery using WebSockets for a subscribed user.

Cron Jobs:

- `sendDailyDigest`: Send a daily email digest of notifications to users.
- `sendWeeklyDigest`: Send a weekly email digest of notifications to users.
- `sendMonthlyDigest`: Send a monthly email digest of notifications to users.
