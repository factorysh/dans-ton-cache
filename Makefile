test:
	go test -timeout 30s -cover github.com/factorysh/dans-ton-cache/disk
	go test -timeout 30s -cover github.com/factorysh/dans-ton-cache/cache
