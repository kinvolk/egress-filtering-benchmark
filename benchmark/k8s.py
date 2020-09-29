from string import Template

benchmark_pod_tmpl = '''
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: default
  name: benchmark-sa
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: benchmark-privileged-psp
  namespace: default
roleRef:
  kind: ClusterRole
  name: privileged-psp
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: ServiceAccount
  name: benchmark-sa
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  # "namespace" omitted since ClusterRoles are not namespaced
  name: benchmark-calico-resource-creator
rules:
- apiGroups: ["crd.projectcalico.org"]
  #resources: ["*"]
  resources: ["globalnetworkpolicies", "globalnetworksets","hostendpoints"]
  verbs: ["get", "watch", "list", "create", "delete","update"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: benchmark-calico-resource-creator
roleRef:
  kind: ClusterRole
  name: benchmark-calico-resource-creator
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: ServiceAccount
  name: benchmark-sa
  namespace: default
---
apiVersion: v1
kind: Pod
metadata:
  name: egress-benchmark
  labels:
    app: benchmark
spec:
  hostNetwork: true
  serviceAccountName: benchmark-sa
  nodeSelector:
    nodetype: worker-benchmark
  containers:
  - name: benchmark
    image: quay.io/imran/benchmark:v0.20
    imagePullPolicy: Always
    securityContext:
      privileged: true
    command:
      - benchmark
      - --count=$count
      - --iface=$iface
      - --seed=$seed
      - --ipnets="$ipnets"
      - --filter=$filter
    env:
    - name: BENCHMARK_COMMAND
      value: "$cmd"
  restartPolicy: Never
'''

