.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmo -f seaandblue/webook:v.0.0.1 .
	@docker build -t seaandblue/webook:v.0.0.1 .