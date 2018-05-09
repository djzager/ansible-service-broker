#!/bin/bash

PROJECT_ROOT=$(dirname "${BASH_SOURCE}")/..
BROKER_IMAGE=${BROKER_IMAGE:-"docker.io/ansibleplaybookbundle/origin-ansible-service-broker:latest"}
APB_NAME=${APB_NAME:-"automation-broker-apb"}
APB_IMAGE=${APB_IMAGE:-"docker.io/automation-broker/automation-broker-apb:latest"}
ACTION=${ACTION:-"provision"}

if which kubectl; then
    CMD=kubectl
else
    CMD=oc
fi

# sed magic to make it possible to reuse install.yaml to deploy and wait for the broker
ARGS="[ \"${ACTION}\", \"-e create_broker_namespace=true\", \"-e wait_for_broker=true\", \"-e broker_image=${BROKER_IMAGE}\" ]"
APB_YAML=$(sed "s%\(image:\).*%\1 ${APB_IMAGE}%; s%\(args:\).*%\1 ${ARGS}%" ${PROJECT_ROOT}/apb/install.yaml)

echo "${APB_YAML}"
echo "${APB_YAML}" | ${CMD} create -f -
sleep 5

${CMD} logs -n ${APB_NAME} "${APB_NAME}" -f
EXIT_CODE=$(${CMD} get pod -n ${APB_NAME} "${APB_NAME}" -o go-template="{{ range .status.containerStatuses }}{{.state.terminated.exitCode}}{{ end }}")

echo "${APB_YAML}" | ${CMD} delete -f -
if [ -n "${EXIT_CODE}" ]; then
    exit ${EXIT_CODE}
else
    exit 0
fi
