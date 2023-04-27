### Localstack issue

A small repo showing a dynamoDB issue we've seen when trying to upgrade from older 1.2 images to 1.4 or even 2.0


`docker-compose -f docker-compose_1.2_works.yml up --build`

The app spins up, creates the table, seeds the data. All is well.

`docker-compose -f docker-compose_2.0_error.yml up --build`

The app spins up, creates the table, localstack shows a 200

`AWS dynamodb.CreateTable => 200`

The app/localstack then throw a 400

`ResourceNotFoundException: Cannot do operations on a non-existent table`

Trying to list the table, it doesn't exist:

`aws dynamodb list-tables --endpoint-url http://localhost:4566 --region us-west-2 --output text`

Rerunning the code again, it seems to think the table exists somewhere:

`docker restart $container-id-for-endpoint`

`ResourceInUseException: Cannot create preexisting table`

Waiting a little makes this go away in 1.4/2.0, but we don't have to wait in 1.2
