# Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: build_api build_frontend clean run-fe run-be run-iam install-deps

all: build_frontend build_backend
	@echo ""
	@echo "Build complete!"
	@echo ""

build_frontend:
	@echo "Building frontend..."
	cd frontend/server && npm run build
	cd frontend/demoapp && npm run build

build_backend:
	@echo "Building backend..."
	go build -o bin/iam-server cmd/server/main.go
	go build -o bin/demo-server cmd/demoapp/main.go

run-fe:
	cd frontend/server && npm run dev

run-be: build_frontend build_backend
	./bin/iam-server

run-iam: build_frontend build_backend
	(trap 'kill 0' SIGINT; ./bin/iam-server & cd frontend/server && npm run dev)

run-demo-fe:
	cd frontend/demoapp && npm run dev

run-demo-be: build_frontend build_backend
	./bin/demo-server

run-demoapp: build_frontend build_backend
	(trap 'kill 0' SIGINT; ./bin/demo-server & cd frontend/demoapp && npm run dev)

install-deps:
	@echo "Installing frontend dependencies..."
	cd frontend/server && npm install
	cd frontend/demoapp && npm install
	@echo "Installing backend dependencies..."
	go mod tidy
	@echo "Dependencies installed!"

clean:
	@echo "Cleaning build artifacts..."
	rm -fr frontend/server/dist
	rm -fr cmd/server/dist
	rm -fr frontend/demoapp/dist
	rm -fr cmd/demoapp/dist
	@echo "Clean complete!"

