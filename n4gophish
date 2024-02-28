#!/bin/bash

CONTAINER_NAME="n4-gophish"
IMAGE_NAME="n4-gophish:latest"


build() {
    docker build -t $IMAGE_NAME .
}

start() {
	GOPHISH_INITIAL_ADMIN_PASSWORD="password"
	GOPHISH_ADMIN_PASSWORD_SHOULD_BE_RESET=false


	# Check if the container is running
	if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
	    echo "Container $CONTAINER_NAME is already running."
	    exit 0
	else
	    echo "Container $CONTAINER_NAME is not running. Building and running it now..."
	fi

	# Build the Docker image if not already built
	if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${IMAGE_NAME}$"; then
	    echo "Building Docker image: $IMAGE_NAME"
	    docker build -t $IMAGE_NAME .
	fi

	# Run the Docker container
	docker run -d --name $CONTAINER_NAME \
               -p 3333:3333 -p 8080:80 \
               -e "GOPHISH_INITIAL_ADMIN_PASSWORD=$GOPHISH_INITIAL_ADMIN_PASSWORD" \
               -e "GOPHISH_ADMIN_PASSWORD_SHOULD_BE_RESET=$GOPHISH_ADMIN_PASSWORD_SHOULD_BE_RESET" \
               --restart=always $IMAGE_NAME
}

logs() {
	# access container logs with follow option
	docker logs $CONTAINER_NAME -f
}

shell() {
	# get a bash shell on container
	docker exec -it $CONTAINER_NAME /bin/bash
}

stop() {
	# stop container
	if docker ps -q --filter name=n4-gophish | grep -q "."; then
        docker stop $CONTAINER_NAME
        echo "Container $CONTAINER_NAME stopped."
    else
        echo "Container $CONTAINER_NAME is not running."
    fi
}

rm () {
	 # Check if the container is running
    if docker ps -q --filter name=n4-gophish | grep -q "."; then
        # Container is running, stop and remove it
        docker stop n4-gophish
        docker rm n4-gophish
        echo "Container $CONTAINER_NAME stopped and removed."
    else
        # Container is not running, remove it
        docker rm n4-gophish
        echo "Container $CONTAINER_NAME removed."
    fi
}

# Main function
main() {
    case "$1" in
        build)
            build
            ;;
		start)
			start
		;;
		logs)
            logs
            ;;
		shell)
            shell
            ;;
        stop)
            stop
            ;;
		rm)
            rm
            ;;
        *)
            echo "Usage: $0 [build|start|logs|shell|stop|rm]"
            exit 1
            ;;
    esac
}

main "$@"
