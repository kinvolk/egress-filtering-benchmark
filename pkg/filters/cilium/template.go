package cilium

const allowAllEgressOnHost = `
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: allow-all-from-host-{{ .Index}}
  labels:
    group: delete
spec:
  nodeSelector:
    matchLabels:
      nodetype: 'worker-benchmark'
  egress:
  - toEntities:
    - all
`
const denyPolicyForHosts = `
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: egress-filtering-benchmark-for-host-network-{{ .Index }}
  labels:
    group: delete
{{- if .Nets }}
spec:
  nodeSelector:
    matchLabels:
      nodetype: 'worker-benchmark'
  egressDeny:
  - toCIDR:
  {{- range $net := .Nets }}
    - {{ $net }}
  {{- end }}
{{- end }}
`
const allowAllEgressOnBenchmarkApp = `
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: allow-all-from-benchmark-app-{{ .Index}}
  labels:
    group: delete
spec:
  endpointSelector:
    matchLabels:
      app: benchmark
  egress:
  - toEntities:
    - all
`
const denyPolicyForBenchmarkApp = `
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: egress-filtering-benchmark-for-benchmark-app-{{ .Index }}
  labels:
    group: delete
{{- if .Nets }}
spec:
  endpointSelector:
    matchLabels:
      app: benchmark
  egressDeny:
  - toCIDR:
  {{- range $net := .Nets }}
    - {{ $net }}
  {{- end }}
{{- end }}
`
