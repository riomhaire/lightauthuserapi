.DEFAULT_GOAL := everything

dependencies:
	@echo Downloading Dependencies
	@go get ./...

build: dependencies
	@echo Compiling Apps
	@echo   --- lightauthuserapi 
	@go build github.com/riomhaire/lightauthuserapi/frameworks/application/lightauthuserapi
	@go install github.com/riomhaire/lightauthuserapi/frameworks/application/lightauthuserapi
	@echo Done Compiling Apps

test:
	@echo Running Unit Tests
	@go test ./...

profile:
	@echo Profiling Code
	@go get -u github.com/haya14busa/goverage 
	@goverage -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

clean:
	@echo Cleaning
	@go clean
	@rm -f lightauthuserapi
	@rm -f coverage-*.html
	@find . -name "debug.test" -exec rm -f {} \;

everything: clean build test profile  
	@echo Done
