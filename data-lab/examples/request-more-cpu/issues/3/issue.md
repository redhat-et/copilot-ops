`@resources:cluster-scope/base/core/namespaces/training-model/resourcequota.yaml`

While the model is definitely happy with an increased amount of resources, the cluster isn't. We need to decrease the amount of resources in `@resources` for both `.limits` and `.requests` by 3/4ths. 