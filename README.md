# Serwer

Expose your computer files to the network. 

## Compiling & Using it

Currently, use the `make compile` to generate a binary. To start the server simply 
call `serwer -directory "<your absolute directory>"` to start a server on the port
8095.

## Security

To use the authentication you have to specify the user and password in a .csv file
of your choice and provide the file to the program when starting it with the 
`-security <file.csv>` flag. The password **must be hashed using SHA512** as the 
server  will hash the sent password from the frontend when comparing passwords.
