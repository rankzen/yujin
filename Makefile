APP := yujin
BUILD_DATE := `date +%FT%T%z`

CSS_SRC := internal/ui/static/css/tailwind.css
CSS_OUT := internal/ui/static/css/styles.css

# .PHONY: \
# 	yujin \
# 	live \
# 	live/templ \
# 	live/server \
# 	live/tailwind \
# 	live/sync_assets

yujin: build-templ build-tailwind
	@ go build -o $(APP) main.go

build-templ:
	@ templ generate

build-tailwind:
	@ npx tailwindcss -i $(CSS_SRC) -o $(CSS_OUT)

live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=true

live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
		--build.cmd "go build -o tmp/bin/yujin main.go" \
		--build.bin "./tmp/bin/yujin ui" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

live/tailwind:
	npx tailwindcss -i ./tailwind.css -o ./pkg/server/ui/static/css/styles.css --watch

live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"

live:
	make -j4 live/templ live/server live/tailwind live/sync_assets
