run:
	go run cmd/web/main.go

test:
	go run gotest.tools/gotestsum@latest --hide-summary=skipped ./...

test-e2e:
	npm run test:e2e

test-e2e-headed:
	npm run test:e2e:headed

install-e2e:
	npm install

setup-e2e: install-e2e
	npx cypress install
