`@resources:cluster-scope/base/core/namespaces/training-model/resourcequota.yaml`

We need to decrease the amount of CPU & Memory in `@resources` for both `.limits` and `.requests` by 0.75, and increase the amount of storage to 100Gi.