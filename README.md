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

```bash
uptime-sentinel-1  | [TASK 01] https://example.org/ -> OK (status 200), time=389.057125ms 
uptime-sentinel-1  | [TASK 04] https://example.org/ -> OK (status 200), time=393.012875ms 
uptime-sentinel-1  | [TASK 05] https://example.org/ -> OK (status 200), time=393.300709ms 
uptime-sentinel-1  | [TASK 02] https://example.org/ -> OK (status 200), time=404.892083ms 
uptime-sentinel-1  | [TASK 03] https://example.org/ -> OK (status 200), time=414.865708ms 
uptime-sentinel-1  | [TASK 06] https://example.org/ -> OK (status 200), time=136.057792ms 
uptime-sentinel-1  | [TASK 07] https://example.org/ -> OK (status 200), time=132.009167ms 
uptime-sentinel-1  | [TASK 08] https://example.org/ -> OK (status 200), time=147.424083ms 
```
