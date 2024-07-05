# stori
Technical Challenge for Stori


## AWS
create a aws accoun or configure your credentials in env variables

if runnin local, make sure to Configure CLI AWS Account Credentials ~/aws/creds

if runing prod, load env vars with

AWS_ACCOUNT_KEY
AWS_ACCOUNT_SECRET
AWS_REGION


## Configure DynamoDB at aws

Configura a DynamoDB table called 

"transactions"


## Deploy Account Txs App into Lambda

Configure CLI AWS Account Credentials 

Configura Lambda Execution Role

``` aws iam create-role --role-name lambda-ex --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}'  ```


Agregar politicas necesarias 
aws iam attach-role-policy --role-name lambda-ex --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
aws iam attach-role-policy --role-name lambda-ex --policy-arn arn:aws:iam::aws:policy/AmazonElasticFileSystemClientFullAccess


Build your programs 

```GOOS=linux GOARCH=amd64 go build -o txs cmd/transactions/main.go```

Zip for Lambda 

``` zip function.zip txs ```

Create the lambda function 

aws lambda create-function --function-name process-account-tx \
--zip-file fileb://function.zip --handler main \
--runtime go1.x --role arn:aws:iam::YOUR_ACCOUNT_ID:role/lambda-execution-role 


<!-- --file-system-configs Arn=arn:aws:elasticfilesystem:YOUR_REGION:YOUR_ACCOUNT_ID:access-point/your-efs-access-point,LocalMountPath=/mnt/efs -->


## Docker Build - Docker Compose run 

change the example.docker-compose.yml to docker-compose.yml

Complete all the missing ENV variables for the app service