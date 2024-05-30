/*
Package awsx is intended to host a few struct/abstraction to be able to use common struct in aws clients
Managed interfaces :

Awslogger : allows the abstraction of logging.Logger to be usable in aws packages

awsx is also in charge of the basic initialisatoin of the aws configuration. The NewConfig is used in awsx & s3x to configure the common parts (auth, ...) of the configuration.
*/
package awsx
