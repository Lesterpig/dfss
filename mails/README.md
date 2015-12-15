### DFSS - Mails lib ###

This library is designed to wrap the smtp library of go.

## Initiating a connection to a server ##

To start a connection to a server, create a CustomClient via NewCustomClient
This takes :
- A sender (ex : qdauchy@insa-rennes.fr)
- A host (ex : mailhost.insa-rennes.fr)
- A port (ex : 587)
- A user (ex : qdauchy)
- A password

This requires the server to have TLS

## Using the connection ##

The connection that has been created can then be used to send one or several mails

Using Send requires :
- A slice of receivers
- A subject
- A message
- A (possibly empty) slice of extensions
- A (possibly empty) slice of filenames. This slice must be of the same length as the extensions one.

## Closing the connection ##

Finally, close the connection using Close.

## Example ###

Refer to the doc's to see the library in practice

## Testing the library ##

The testing file uses the following variables to set up the tests :
DFSS_TEST_MAIL_SENDER
DFSS_TEST_MAIL_HOST
DFSS_TEST_MAIL_PORT
DFSS_TEST_MAIL_USER
DFSS_TEST_MAIL_PASSWORD
DFSS_TEST_MAIL_RCPT1
DFSS_TEST_MAIL_RCPT2
