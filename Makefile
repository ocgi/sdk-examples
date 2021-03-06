# Copyright 2021 The OCGI Authors.
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

REGISTRY_NAME=ocgi
VERSION=0.1
OUTDIR=${PWD}/bin

all: fmt vet simple-tcp simple-tcp-image push

clean: ## make clean
	rm -rf ${OUTDIR}

fmt:
	@echo "run go fmt ..."
	@go fmt ./simple-tcp/...

vet:
	@echo "run go vet ..."
	@go vet ./simple-tcp/...

simple-tcp: fmt vet
	@echo "build simple-tcp"
	cd simple-tcp && GOOS=linux go build -ldflags "-X 'main.Version=$(VERSION)'" -o ${OUTDIR}/simple-tcp

simple-tcp-image: simple-tcp
	docker build -t $(REGISTRY_NAME)/simple-tcp:$(VERSION) -f simple-tcp/Dockerfile .

push: simple-tcp-image
	docker push $(REGISTRY_NAME)/simple-tcp:$(VERSION)
