# go-fiber-api-template
starter pack for building backend with go and fiber, with jwt auth


# authentication

there are few steps for authentication
steps:


## check if user exists in our db, by it's phone number:
route: /api/auth/checkusername



## if user exists:
get the password and send tokens
route: /api/auth/token




## if user not exists:
send a verification_id to ui, and a pin number to user's phone, get the pin and verification_id from ui
route: /api/auth/verify




## complete registration
get the user info and verification_id from ui
route: /api/auth/completeregistration





# authorization
there are several roles that can be apply to users, the roles defenition are under user model.
for access a route to a role group simpley apply a middware in routes. (see contacts route for example)


# routes and handlers
there are 5 main generic handler for handeling data, get all items, get single item, delete, update and create an item.
get items support pagination and filtering. 
in this senario you just need to create your model and and routes with authorization middleware, see contacts and users model for example.
#
