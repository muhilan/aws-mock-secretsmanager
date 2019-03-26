# Introduction

aws-mock-secretsmanager is a lightweight AWS secretsmanager implemented in golang, mostly useful in tests.

####Supported Operations

get-secret-value

#### Adding data

Mount a folder containing certs/secrets and they will be loaded recursively

##### Consider the following files in a folder
`` ->  signing-cert.pem    ``

`` ->  mysecret    ``
 
`` ->  encryption.private.key    ``

`` ->  encryption-cert    ``

###### The following would be the Secret Ids
``signing-cert``

``mysecret``

`` encryption/private ``

`` encryption-cert``

