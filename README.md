# zapgpt

Integration between whatsapp and chatgpt

## Requirements

* serverless framework
* make
* AWS account
* Twillio account
* OpenAPI (chatgpt) account

## References

This project is a hands-on training from this Full Cycle course: [Integração do Chat GPT com WhatsApp + Abertura das matrículas](https://www.youtube.com/watch?v=01aejNssbA4)

## Run:

envs:

* MY_AWS_ACCESS_KEY_ID
* MY_AWS_SECRET_ACCESS_KEY
* GPT_API_KEY

build

`$ make build`

deploy

`$ serverless deploy`

After deploying, copy the endpoint and add it into your twillio account settings

## Troubleshooting

serverless deploy errors: 
IAM rules will be required when deploying to aws. Make sure to give the proper permissions to the user.