FROM quay.io/openshift/origin-operator-registry

COPY olm-catalog manifests/automationbroker
RUN ./initializer -o bundles.db

ENTRYPOINT ["/registry-server"]
CMD ["--database", "bundles.db"]
