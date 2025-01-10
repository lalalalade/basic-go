.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f lalalalade/webook:v.0.0.1 .
	@docker build -t lalalalade/webook-test:v.0.0.1 .