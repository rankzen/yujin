templ:
	@templ generate

tailwind:
	@npx tailwindcss -i ./static/css/tailwind.css -o ./static/css/output.css

install:
	go install github.com/a-h/templ/cmd/templ@latest
	npm install -D tailwindcss

.PHONY: live/templ
live/templ:
	@templ generate --watch --open-browser=false

live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

.PHONY: live/tailwind
live/tailwind:
	@npx tailwindcss -i ./static/css/tailwind.css -o ./static/css/output.css --minify --watch

live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"

.PHONY: live
live:
	make -j2 live/templ live/tailwind
