# uptime-sentinel

Small concurrent URL checker written in Go.  


## Project structure

```text
.
├─ main.go
├─ go.mod
├─ Dockerfile
└─ docker-compose.yml
```

## How to run

1. Build the image:

   ```bash
   docker compose build
   ```

2. Run the service:

   ```bash
   docker compose up
   ```

   You will see output similar to:

   ```text
   [TASK 01] https://example.org -> OK (status 200), time=120.345ms
   [TASK 02] https://google.com -> OK (status 200), time=230.789ms
   [TASK 03] https://httpstat.us/200 -> OK (status 200), time=150.123ms
   [TASK 04] https://httpstat.us/404 -> NOT OK (status 404), time=140.456ms
   ```
