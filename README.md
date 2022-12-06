## BASIC MICROSERVICE PROJECT IN GO

This project is a basic blog application backend, my attempt at implementing a microservice application with different communication protocols

### Here's a link to the proposed architecture [link](https://lucid.app/lucidchart/f7ca1b27-1270-4bd3-97c5-4df8510cc775/edit?viewport_loc=43%2C-128%2C3096%2C1632%2C0_0&invitationId=inv_394bdf6c-3b0f-4251-ab3c-112bd78c32d2)

### FEATURES:
1. Authentication
2. Create blog post
3. Update blog post
4. Delete blog post
5. Fetch blog posts
6. Comment on blog posts
7. Replies on comments
8. Delete comments
9. Update comments


### SERVICES
1. Auth
2. Gateway
3. Avater generator
4. Upload
5. Logging
6. Notification


## AUTH (SYNC)
1. Signup
  - connects to database (sql)
  - verify user email
  - store user information
  - send user id to avater-generator service queue
1. Login
  - verify user info
  - return jwt token
1. Forgot Password
  - send verification OTP to mail
  - validate otp and update password

## AVATER-GENERATOR (ASYNC)
1. Subscribe to avater-generator service
2. Generate a custom avater based on first letters of user first and last names.
3. Upload generated avater to AWS S3
4. Update user record with avater url

## LOGGER (ASYNC)
1. connects to database (mongo)
2. save logs to datebase

## NOTIFICATION (ASYNC)
1. Push Notifications
2. Email
3. SMS