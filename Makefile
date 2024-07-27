include Makefile.env

## run: Compile and run server
run: go-compile start

## start: Start in development mode. Auto-starts when code changes.
start: start-favorites

## stop: Stop development mode. GO_FAVORITES
stop: stop-favorites

start-favorites: stop-favorites
	@echo "  >  $(PROJECTNAME) is available at $(ADDRESS)"
	@-$(GOBIN)/favorites & echo $$! > $(PID_GO_FAVORITES)
	@cat $(PID_GO_FAVORITES) | sed "/^/s/^/  \>  PID: /"

stop-favorites:
	@echo "  >  stop by $(PID_GO_FAVORITES)"
	@-touch $(PID_GO_FAVORITES)
	@-kill `cat $(PID_GO_FAVORITES)` 2> /dev/null || true
	@-rm $(PID_GO_FAVORITES)

restart-favorites: stop-favorites start-favorites

## build: Build and the binary compile server
build: go-build-favorites

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

go-compile: go-build-favorites

go-build-favorites:
	@echo "  >  Building GO_FAVORITES binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) cd ./cmd/favorites && go build -o $(GOBIN)/favorites $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

.PHONY: go-update-deps
go-update-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-swag:
	swag init -g $(MAIN_GO) --output $(DOCS_DIR)
	sed -i 's/"localhost:8080",/env.GetConfig().Address(),/' $(DOCS_GO)
	goimports -w $(DOCS_GO)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

cert:
	@cd cert; openssl req -x509 -newkey rsa:1024 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Tech School/OU=Education/CN=localhost/emailAddress=none@gmail.com"
	@echo "CA's self-signed certificate"
	@cd cert; openssl x509 -in ca-cert.pem -noout -text
	@cd cert; openssl req -newkey rsa:1024 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Tech School/OU=Education/CN=localhost/emailAddress=none@gmail.com"
	@cd cert; openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
	@echo "Server's signed certificate"
	@cd cert; openssl x509 -in server-cert.pem -noout -text

##################
# Implicit targets
##################

# This rulle is used to generate the message source files based
# on the *.proto files.
%.pb.go: %.proto
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./$<

%_search_service_mock_test.go: %_search_service.go
	@mockgen -source=./$< -package=services > ./$@

$(SEARCH_SERVICE_DIR)/repo_mock_domain_test.go: $(REPO_DIR)/repo.go
	@mockgen -source=./$< -package=services > ./$@

$(SEARCH_SERVICE_DIR)/transactional_mock_domain_test.go: $(REPO_DIR)/transactional.go
	@mockgen -source=./$< -package=services > ./$@

$(CONTROLLERS_DIR)/api_favorites_service_mock_test.go: $(SEARCH_SERVICE_DIR)/api_favorites_service.go
	@mockgen -source=./$< -package=controllers > ./$@

####################################
# Major source code-generate targets
####################################
generate: $(PROTO_PB_GO) $(SEARCH_SERVICE_MOCKS) $(REPO_MOCKS) $(FAVORITES_MOCK)
	@echo "  >  Done generating source files based on *.proto and Mock files."

test:
	@echo "  > Test Iteration ..."
	go vet -vettool=$(which statictest) ./...
	cd cmd/favorites && ./favoritestest -test.v -test.run=^TestGophermart$$ -favorites-binary-path=./favorites -favorites-host=localhost -favorites-port=$(GO_FAVORITES_PORT) -favorites-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable" -accrual-binary-path=./accrual -accrual-host=localhost -accrual-port=$(ACCRUAL_PORT) -accrual-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"

.PHONY: cert help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
