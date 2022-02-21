# Listed below are:
# 1. Explanation of how two files need to be changed
# 2. The original files, each one separated by a '---' string
# 3. Files updated with the described changes, with a '---' string in-between each file
# 4. A '####' string, indicating the end of the document


## Description of issues:
`@file1:grafana/base/datasource.yaml`
`@file2:grafana/base/grafana-route.yaml`

The amount of CPU in @file1 needs to be increased to 512M
The PVC storage amount in @file2 should be decreased to 10Gi


## Original files:
# @file1
apiVersion: integreatly.org/v1alpha1
kind: GrafanaDataSource
metadata:
  name: datasource
spec:
  name: prometheus-grafanadatasource.yaml
  datasources:
    - name: Prometheus
    - access: proxy
      editable: true
      isDefault: true
      jsonData:
        httpHeaderName1: 'Authorization'
        timeInterval: 5s
        tlsSkipVerify: true
      name: Prometheus
      secureJsonData:
        httpHeaderValue1: 'Bearer ${BEARER_TOKEN}'
      type: prometheus
      url: 'https://thanos-querier.openshift-monitoring.svc.cluster.local:9091'
---
# @file2
kind: Route
apiVersion: route.openshift.io/v1
metadata:
  name: grafana
  annotations:
    kubernetes.io/tls-acme: "true"
spec:
  host: grafana.operate-first.cloud
  to:
    kind: Service
    name: grafana-service
  port:
    targetPort: 3000

## Updated files:
