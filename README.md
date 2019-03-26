# Introduction

aws-mock-secretsmanager is a lightweight AWS secretsmanager implemented in golang, mostly useful in tests.

#### Supported Operations

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

#### Example Usage

``docker run -d -v /path/to/certs:/data -p 8080:8080 mgmuhilan/aws-mock-secretsmanager:0.2.0 ``