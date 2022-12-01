BASE_DIR 									?= ${PWD}
KIND_K8S_VERSION 					?= "1.23.6"
KIND_CLUSTER_NAME 				?= "monitor"
DOCKER_IMG_NAME						?= "monitor-url-prometheus"

deploy-prometheus-monitor:
	@ kubectl create -f ${BASE_DIR}/deploy/prometheus
.PHONY: deploy-prometheus-monitor

deploy-app:
	@ kubectl create ns apps && \
		kubectl create -f ${BASE_DIR}/deploy/app
.PHONY: deploy-app

docker-build:
	@ docker build -f ${BASE_DIR}/app/Dockerfile -t ${DOCKER_IMG_NAME}:latest ${BASE_DIR}/app && \
		docker image ls | grep "${DOCKER_IMG_NAME}" && \
		kind load docker-image ${DOCKER_IMG_NAME} --name ${KIND_CLUSTER_NAME} ## Kind known issue: https://kind.sigs.k8s.io/docs/user/known-issues/#unable-to-pull-images
.PHONY: docker-build

kube-prometheus-install:
	@ rm -rf ${BASE_DIR}/tmp && \
		mkdir ${BASE_DIR}/tmp && \
		git clone --depth 1 https://github.com/prometheus-operator/kube-prometheus.git -b release-0.10 ${BASE_DIR}/tmp/ && \
		kubectl create -f ${BASE_DIR}/tmp/manifests/setup/ && \
		kubectl create -f ${BASE_DIR}/tmp/manifests/
.PHONY: kube-prometheus-install

kind-create-cluster: kind-delete-cluster
	@	kind create cluster \
		--image kindest/node:v${KIND_K8S_VERSION} \
		--config ${BASE_DIR}/kind.yaml \
		--name ${KIND_CLUSTER_NAME}
.PHONY: kind-create-cluster

kind-delete-cluster:
	@	kind delete cluster --name ${KIND_CLUSTER_NAME}
.PHONY: kind-delete-cluster