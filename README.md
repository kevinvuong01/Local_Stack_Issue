### Localstack Issue

This is a small demo app to demonstrate the problem I'm seeing.

In `docker-compose.yml` we specify two services:

* `endpoint` the web server
* `janus-localstack` the local version of DynamoDB

You can start both with `make up` and stop them with `make down`.

Note: `janus-localstack` is our interal version of [localstack](https://github.com/localstack/localstack) that does some provisioning of AWS resources specified in `.janus/deploy.yml`.

---
## THE PROBLEM:

This app is supposed to seed a DynamoDB table with a couple records when the web server starts up. The DynamoDB table gets created from lines 16-20 in `.janus/deploy.yml`. 

That works perfectly when we use this version of `janus-localstack`:  
`latest@sha256:e6d306a60d24ab56980b284428a97b5faf3b7f333423c57c285ce8c7a9d53fea`


However, the above is an outdated version of `janus-localstack`. When we tell `docker-compose.yml` to use a newer version of `janus-localstack` (on lines 29-30 of `docker-compose.yml`) the DynamoDB seeding fails because the Dynamo table hasn't yet been created, even though the `janus-localstack` service reports itself as being healthy. You can see it panic in the log output.

So in sum:

**This version of `janus-localstack` works as expected:**  
`latest@sha256:e6d306a60d24ab56980b284428a97b5faf3b7f333423c57c285ce8c7a9d53fea`

**And this one does not:**  
`latest@sha256:88d5657058e4e6b2980a02432ff20c5685e764bb2fbb4a92596408c87a0d55c9`


#### Things to consider: 
* It's possible that it's unrealistic to expect a table to be available when the `janus-localstack` container reports itself as healthy. 
* It's possible/likely that this was caused by something in the upstream [localstack](https://github.com/localstack/localstack) repo, in which case it's not our bug.


---

#### Helpful things
The `janus-localstack` code that's responsible for provisioning the DynamoDB tables is here:  
https://github.com/cambiahealth/janus-localstack/blob/master/janus-bootstrap/janus-localstack-provisioner/dynamodb.js


**List Dynamo tables**: `aws --endpoint-url http://localhost:4566 dynamodb list-tables`

**NoSQL Workbench:**
https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/workbench.html
