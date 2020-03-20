.PHONY:	edgecloud

edgecloud:
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	go build edgecloud.go
	@chmod 777 edgecloud


.PHONY: clean
clean:
	@rm -rf edgecloud
	@echo "[clean Done]"
