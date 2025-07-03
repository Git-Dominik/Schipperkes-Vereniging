# Schipperkes setup

# Local/dev
1. Copy the .env.sample and name it .env
2. Fill the variables in the .env
3. Run the code base in the root directory
```go run .```

# Prod (docker)
1. Build docker image (assumes you are in the project dir)
```docker build -t schipperkes .```
2. Run the docker compose
```docker compose up -d```

