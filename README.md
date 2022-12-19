## BASIC MICROSERVICE PROJECT IN GO

This project is a basic blog application backend, my attempt at implementing a microservice application with different communication protocols

### Here's the proposed architecture 
![microservice-architecture image](microservice-architecture.png)

[link](https://drive.google.com/file/d/1xaSWEzuC7NARDynK8X6u38MIRKt9ptMt/view?usp=sharing)


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
2. Login
  - verify user info
  - return jwt token
3. Forgot Password
  - send verification OTP to mail
  - validate otp and update password

## AVATER-GENERATOR (ASYNC)
1. Generate a custom avater based on first letters of user first and last names.
2. Upload generated avater to AWS S3
3. Update user record with avater url

## LOGGER (ASYNC)
1. save logs to datebase (mongo)

## NOTIFICATION (ASYNC)
1. Push Notification
2. Email
3. SMS