LOCAL_PRIVATE_REPO=127.0.0.1:5000
VERSION=v1.0.0
GO111MODULE=off
# Image URL to use all building/pushing image targets
IMG ?= ctripcloud/namespace-delete-checker:${VERSION}

# Install CRDs into a cluster
install:
	kubectl apply -f yaml/

uninstall:
	kubectl delete -f yaml/

#create-secrets:
#    kubectl create secret generic certs --from-file=key.pem=/etc/kubernetes/pki/cert.pem --from-file=cert.pem=/etc/kubernetes/pki/key.pem --dry-run -o yaml | kubectl -n default apply -f -

# Build the docker image
docker-build:
	docker build . -t ${IMG}
	# sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./yaml/statefulset.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

docker-push-local-private:
	docker tag $(IMG) $(LOCAL_PRIVATE_REPO)/$(IMG)
	docker push $(LOCAL_PRIVATE_REPO)/$(IMG)