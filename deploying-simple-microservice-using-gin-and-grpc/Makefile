.PHONY: proto
proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ if ! which protoc-gen-go > /dev/null; then \
		echo "error: protoc-gen-go not installed" >&2; \
		exit 1; \
	fi
	@ if ! which protoc-gen-go-grpc > /dev/null; then \
		echo "error: protoc-gen-go-grpc not installed" >&2; \
		exit 1; \
	fi
	for file in $$(git ls-files '*.proto'); do \
		protoc -I $$(dirname $$file) \
		--go_out=:$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=:$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		$$file; \
	done

.PHONY: docker-build
docker-build:
	docker build -t jxlwqq/api-server -f cmd/api-server/Dockerfile .
	docker build -t jxlwqq/calculator -f cmd/calculator/Dockerfile .

.PHONY: kube-deploy
kube-deploy:
	kubectl apply -k ./kubernetes

.PHONY: kube-undeploy
kube-undeploy:
	kubectl delete -k ./kubernetes