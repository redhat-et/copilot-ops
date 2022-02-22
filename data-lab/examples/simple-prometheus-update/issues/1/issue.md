`@file1:grafana/base/datasource.yaml`
`@file2:grafana/base/grafana-route.yaml`

@file1 and @file2 must specify the `grafana-datasource` namespace.
@file1 needs to increase the 'Prometheus' data source to a time interval of 20s.
Disable tls-acme for @file2.
