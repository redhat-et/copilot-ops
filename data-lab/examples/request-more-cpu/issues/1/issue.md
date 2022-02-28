`@kustomization:cluster-scope/base/core/namespaces/training-model/kustomization.yaml`
`@resourcequota:cluster-scope/base/core/namespaces/training-model/resourcequota.yaml`

The learning model is lagging during testing. We are attempting to add more CPU limits to resolve not having GPU. `@resourcequota` needs to increase the allotted `limits.cpu` and `requests.cpu` count to 64, as well as increasing the `limits.memory` and `requests.memory` to 128Gi.
