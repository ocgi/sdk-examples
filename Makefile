# Copyright 2020 THL A29 Limited, a Tencent company.
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

REGISTRY_NAME=hub.oa.com/dbyin
VERSION=0.1
OUTDIR=${PWD}/bin

all: fmt vet simpletcp simpletcp-image push

clean: ## make clean
	rm -rf ${OUTDIR}

fmt:
	@echo "run go fmt ..."
	@go fmt ./simple-tcp/...

vet:
	@echo "run go vet ..."
	@go vet ./simple-tcp/...

simpletcp: fmt vet
	@echo "build simple-tcp"
	cd simple-tcp && go build -o ${OUTDIR}/simple-tcp

simpletcp-image: simple-tcp
	docker build -t $(REGISTRY_NAME)/carrier-simpletcp:$(VERSION) -f simple-tcp/Dockerfile .

push: simpletcp-image
	docker push $(REGISTRY_NAME)/carrier-simpletcp:$(VERSION)
