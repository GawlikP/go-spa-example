## GO SPA EXAMPLE

This project is a project that demos a crud application written using <b>Go</b> and <b>Vue</b> JavaScript framework as frontend.
The forntend part will be hosted as spa on the go server.

## Runing the project

To run the project we need to only set the environment variables as described in `.env.example` file.
Then we can use the `docker-compose.yml` file to run the docker compose. For example

```
    cp .env.exmaple .env
    docker-compose up
```

After that by default our app should be accessible on `http://localhost:8080` url.
