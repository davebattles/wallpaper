help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: help build install

build: ## build daily-wallpaper
	CGO_ENABLED=0 go build

install: ## install daily-wallpaper
	CGO_ENABLED=0 go install
	cp daily-wallpaper.service _tmp_daily-wallpaper.service
	sed -i "s|GOPATH|${GOPATH}|g" _tmp_daily-wallpaper.service
	sed -i "s|UNSPLASH_CLIENT_ID|${UNSPLASH_CLIENT_ID}|g" _tmp_daily-wallpaper.service
	sudo mv _tmp_daily-wallpaper.service /etc/systemd/system/daily-wallpaper.service
	sudo systemctl daemon-reload
	sudo systemctl enable daily-wallpaper
