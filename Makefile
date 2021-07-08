APP=short-url
IMAGE_TAG=$(shell git log --pretty=format:"%ad_%h" -1 --date=short)
dev:
	go run main.go

build:
	go build -ldflags "-s -w"
	upx -9 $(APP)
	docker build -t 9d77v/$(APP):$(IMAGE_TAG) .
	docker push 9d77v/$(APP):$(IMAGE_TAG)
	rm -r $(APP)