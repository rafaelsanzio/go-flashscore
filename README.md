# Flashscore ⚽️

An API to handle soccer tournaments

## Acknowledgements 📚

- [Golang](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Kafka](https://kafka.apache.org/)
- [JWT](https://jwt.io/)

## Run Locally ▶️

Clone the project

```bash
  git clone git@github.com:rafaelsanzio/go-flashscore.git
```

Go to the project directory

```bash
  cd go-flashscore
```

### Setting environment variables

To run this project, you will need to add the following environment variables to your .env file

`APP_PORT`
`APP_AUTH_PORT`

`MONGO_PASSWORD`
`MONGO_USERNAME`
`MONGO_DATABASE`
`MONGO_PORT`
`MONGO_URI`

`SECRET_API_KEY`

`KAFKA_ADDRESS_1`
`KAFKA_ADDRESS_2`
`KAFKA_ADDRESS_3`

Start the server

```bash
  docker-compose up
```

## Running Tests 🧪

To run tests, run the following command

```bash
  make test
```

## API's Documentation 📑

- [Auth](https://github.com/rafaelsanzio/go-flashscore/tree/main/cmd/auth)
- [API](https://github.com/rafaelsanzio/go-flashscore/tree/main/cmd/api)
