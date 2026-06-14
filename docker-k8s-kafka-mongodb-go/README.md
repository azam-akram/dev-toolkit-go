# docker-k8s-kafka-mongodb-go

## HTTP Server

HTTP Sever is written in go.
Build docker image of http server by,
```
docker build . -t http-server-go
```

and run:
```
docker run --rm -p 8080:8080 --name http-server-go-container http-server-go
```

Output:
```
Request &{GET /healthcheck HTTP/1.1 1 1 map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Connection:[keep-alive] Postman-Token:[bed746f4-fdf3-442a-b9a7-e6fd16fde267] User-Agent:[PostmanRuntime/7.54.0]] {} <nil> 0 [] false localhost:8080 map[] map[] <nil> map[] 172.17.0.1:40782 /healthcheck <nil> <nil> <nil> / 0xc000128050 0xc0000b2180 [] map[]}Request &{POST /students HTTP/1.1 1 1 map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Connection:[keep-alive] Content-Length:[345] Content-Type:[application/json] Postman-Token:[545fcdbe-1e47-4853-abf7-e19c69bd39c1] User-Agent:[PostmanRuntime/7.54.0]] 0xc00006e000 <nil> 345 [] false localhost:8080 map[] map[] <nil> map[] 172.17.0.1:40782 /students <nil> <nil> <nil> / 0xc000070000 0xc0000b2180 [] map[]}Student Request: {Admin [{1 a 12} {2 b 9} {1 c 8}]}
Student Response: {v1 88ffd1b654da 2026-06-14 15:20:23.222955345 +0000 UTC m=+60.500957286 Admin [{1 a 12} {2 b 9} {1 c 8}]}
```

## Troubleshoot docker container

- Check running containers
```
docker ps
```

- Check all containers (including stopped ones):
``` 
docker ps -a
```

- Inspect container details
```
docker inspect http-server-go-container
```

- View container logs
```
docker logs <container-name>
```

Example:
```
docker logs http-server-go-container
```

- Follow logs live:
```
docker logs -f http-server-go-container
```

- Show last 100 lines:
```
docker logs --tail 100 http-server-go-container
```

- Remove a Stopped Container
```
docker rm <container_id>
```

Deletes a stopped container. Use -f to force-remove a running container.

- Remove All Stopped Containers
```
docker container prune
```

- Enter the container

For Alpine:
```
docker exec -it http-server-go-container sh
```
For Ubuntu/Debian images:
```
docker exec -it http-server-go-container bash
```

Inside, you can:
```
pwd
ls -la
ps
env
```

- Check why a container exited
```
docker ps -a
```

Example:
```
CONTAINER ID   IMAGE            STATUS
abc123         http-server-go   Exited (1) 10 seconds ago
```

Then:
```
docker logs abc123
```
or
```
docker inspect abc123
```
Look for:
```
"ExitCode": 1
```

- Run interactively

Very useful during development:
```
docker run --rm -it -p 8080:8080 http-server-go
```

You'll see your Go logs directly in the terminal.

- Verify port mapping
```
docker port http-server-go-container
```
Expected:
```
8080/tcp -> 0.0.0.0:8080
```

- Test from inside the container

Enter:
```
docker exec -it http-server-go-container sh
```
Then:
```
wget -qO- http://localhost:8080/healthcheck
```