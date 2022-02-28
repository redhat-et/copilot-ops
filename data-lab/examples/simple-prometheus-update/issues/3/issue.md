`@grafanafile:grafana/base/datasource.yaml`

@grafanafile needs to add a new httpHeaderValue which sets content-type to `application/json` to the .secureJsonData entry in the name=prometheus datasource in .spec.datasourcces. 