# Golang Image resizer

This is a simple program that run in AWS Lambda it performs the following:
 - Is invokable via a API call
 - Accepts utp 3 parameters on the URL:
    - url : eg ?url=https://{urlTo}/{image.png,jpg,jpeg}
    - width : eg &width=100
    - height : eg &height=100
 exmaple:
 ```bash
    https://{apigwUrl}/prod/convert?url=https://cdn.kukiel.dev/cows.jpg&width=500&height=500 
```
 
 Width or height is required, if one of the other is missing the image will be sized according
 to the value that was present and resized to the correct aspect ratio.
 
```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── resize                      <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
└── template.yaml               <-- SAM template
```

## Requirements

* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target 

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run
 
```shell
make
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

Then browse to:
```bash
 https://127.0.0.1:3000/convert?url=https://cdn.kukiel.dev/cows.jpg&width=500&height=500 
```

### Deployment
Run:
```bash 
sam deploy --guided
```
And follow the steps, for subsequent deployments:
```bash 
make deploy
```