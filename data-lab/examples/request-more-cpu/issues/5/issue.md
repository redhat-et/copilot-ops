`@resources:cluster-scope/base/core/namespaces/training-model/resourcequota.yaml`

Woohoo! We finally obtained a GPU! We can reduce the amount of CPUs required to just 8 in `limits` and `requests`; however, we need to introduce vRAM limits of 32Gi to `limits` & `requests`. 