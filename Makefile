
build:
	go build -o dist/server cmd/*.go

run:
	go run cmd/*.go

export:
	go run cmd/*.go export

deploy: export
	wrangler pages deploy dist --project-name go-weebsite --branch main --commit-dirty=true

