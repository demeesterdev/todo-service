# Cloud-native TODO application. (CNATA)

> This application is used as a use case for training in cloud native tooling.
> As such the architecture might not be optimal or secure.

CNATA is a work tracking application working with TODO's

TODO's are defined as a documents with a title and a description.  
The state of a TODO can be open or closed.
Each TODO can be deleted as well but they will be put into trash so they can be retrieved.

Each user will only be able to see their own todo's unless shared with other users.

## Components

 - authorization (used for sign in/up) 
 - Todo service (used to retrieve TODO's)
 - frontend/api service (used as single point interaction for client)


