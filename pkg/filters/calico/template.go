package calico

const globalNetworkSetTmpl = `
apiVersion: crd.projectcalico.org/v1
kind: GlobalNetworkSet
metadata:
  name: egress-filtering-benchmark-{{ .Index }}
  labels:
    type: egress-deny-list
{{- if .Nets }}
spec:
  nets:
  {{- range $net := .Nets }}
  - {{ $net }}
  {{- end }}
{{- end }}
`

const globalNetworkPolicyTmpl = `
apiVersion: crd.projectcalico.org/v1
kind: GlobalNetworkPolicy
metadata:
  name: 00-egress-filtering-benchmark
spec:
  selector: host-endpoint == 'ingress' && nodetype == 'worker'
  order: 0
  #doNotTrack: true
  #applyOnForward: true
  types:
  - Egress
  egress:
  - action: Deny
    protocol: TCP
    destination:
      selector: type == 'egress-deny-list'
  - action: Deny
    protocol: UDP
    destination:
      selector: type == 'egress-deny-list'
  - action: Deny
    protocol: ICMP
    destination:
      selector: type == 'egress-deny-list'
`

const gnpTmplForWorkloads = `
apiVersion: crd.projectcalico.org/v1
kind: GlobalNetworkPolicy
metadata:
  name: 00-egress-filtering-benchmark-for-non-namespaced-resources
spec:
  selector: global()
  order: 0
    #doNotTrack: true
    #applyOnForward: true
  types:
  - Egress
  egress:
  - action: Deny
    protocol: TCP
    destination:
      selector: type == 'egress-deny-list'
  - action: Deny
    protocol: UDP
    destination:
      selector: type == 'egress-deny-list'
  - action: Deny
    protocol: ICMP
    destination:
      selector: type == 'egress-deny-list'
`
