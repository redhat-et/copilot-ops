import KubernetesObject from "@dot-i/k8s-operator";

export interface CompletionResource extends KubernetesObject {
  metadata: any;
  spec: CompletionSpec;
  status: CompletionStatus;
}

export interface CompletionSpec {
  userPrompt: string;
}

export interface CompletionStatus {
  observedGeneration?: number;
  completion?: string;
}
