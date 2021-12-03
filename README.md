# Go Client 

- It allows you to write very clean client codes. 
- You don't need to worry about retrying mechanisms, retrying intervals...
- Highly consistent when you have a failure and if your retry count is finally reached the end library provides you a dead-letter mechanism. So whenever you reached retry limit it will use the interface to send this letter to some database, message queue, file etc.
- Provides hands-on rate limiting for you. You only need to worry about the configurations