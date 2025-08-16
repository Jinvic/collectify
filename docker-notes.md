# Build and run the application in a Docker container

## 1. Build the Docker image

- This uses the multi-stage Dockerfile to build both frontend and backend, resulting in a small, optimized final image.
- The image is tagged as 'collectify-app'.

```bash
docker build -t collectify-app .
```

## 2. Run the Docker container

- This command runs the 'collectify-app' image in detached mode (-d)
- maps port 8080 of the container to port 8080 on the host (-p 8080:8080)
- and names the running container 'collectify-container' (--name collectify-container)
- The application inside the container will be accessible at <http://localhost:8080>

```bash
docker run -d -p 8080:8080 --name collectify-container collectify-app
```

## --- Other useful Docker commands ---

### Stop the running container

```bash
docker stop collectify-container
```

### Start the stopped container

```bash
docker start collectify-container
```

### Remove the container (must be stopped first)

```bash
docker rm collectify-container
```

### Remove the image

```bash
docker rmi collectify-app
```

### View logs from the running container

```bash
docker logs collectify-container
```

### Execute a command inside the running container (e.g., open a shell)

```bash
docker exec -it collectify-container sh
```
