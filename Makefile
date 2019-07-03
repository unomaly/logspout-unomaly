NAME=logspout-unomaly
BUILD_DIR=build

# If you want to configure the Honeycomb Logspout adapter with environment
# variables, set them here. Otherwise you need to use the RoutesAPI.
HONEYCOMB_WRITE_KEY=
HONEYCOMB_DATASET=

# Builds a Docker image of LogSpout with the Honeycomb adapter included.
docker:
	docker build -t $(NAME) image

clean:
	rm -rf $(BUILD_DIR)
	docker rmi -f $(NAME)

clean-images:
	docker images | grep none | awk '{print $3}' | xargs docker rmi -f
