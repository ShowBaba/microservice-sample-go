#! /bin/bash
while getopts e: flag
do
    case "${flag}" in
        e) env=${OPTARG};;
    esac
done

if [ -z "$env" ]; then env="test"; fi

echo "Environment: $env";
echo "[$env] Building auth-service ..."
GOOS=linux GOARCH=amd64 go build -o cicd/$env/bin/app main.go


if [ -f "cicd/$env/bin/app" ]; then
    cp .env cicd/$env/bin/.env
    cd cicd/$env
    docker build --no-cache --tag auth-service .
    cd ../../
fi
