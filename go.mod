module rancherlabs/cattle-drive

go 1.24.0

replace (
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20240205102821-ed248439462a
	k8s.io/api => k8s.io/api v0.28.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.28.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.28.4
	k8s.io/apiserver => k8s.io/apiserver v0.28.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.28.4
	k8s.io/client-go => github.com/rancher/client-go v1.28.6-rancher1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.28.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.28.4
	k8s.io/code-generator => k8s.io/code-generator v0.28.4
	k8s.io/component-base => k8s.io/component-base v0.28.4
	k8s.io/cri-api => k8s.io/cri-api v0.28.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.28.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.28.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.28.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.28.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.28.4
	k8s.io/kubectl => k8s.io/kubectl v0.28.4
	k8s.io/kubelet => k8s.io/kubelet v0.28.4
	k8s.io/kubernetes => k8s.io/kubernetes v1.28.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.28.4
	k8s.io/metrics => k8s.io/metrics v0.28.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.28.4
)

require (
	github.com/briandowns/spinner v1.23.0
	github.com/charmbracelet/lipgloss v0.9.1
	github.com/rancher/rancher/pkg/apis v0.0.0-20230915232223-a9ea4ce4a5ba
	github.com/rancher/wrangler v1.1.1
	github.com/sirupsen/logrus v1.9.3
	k8s.io/client-go v12.0.0+incompatible
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/harmonica v0.2.0 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rancher/norman v0.0.0-20230831160711-5de27f66385d // indirect
	github.com/rivo/uniseg v0.4.6 // indirect
	github.com/sahilm/fuzzy v0.1.1-0.20230530133925-c48e322e2a8f // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/charmbracelet/bubbles v0.18.0
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.10.2 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rancher/aks-operator v1.3.0-rc1 // indirect
	github.com/rancher/eks-operator v1.4.0-rc1 // indirect
	github.com/rancher/fleet/pkg/apis v0.0.0-20231017140638-93432f288e79 // indirect
	github.com/rancher/gke-operator v1.3.0-rc2 // indirect
	github.com/rancher/lasso v0.0.0-20230830164424-d684fdeb6f29
	github.com/rancher/rke v1.5.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/urfave/cli/v2 v2.27.1
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/oauth2 v0.15.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.28.6
	k8s.io/apiextensions-apiserver v0.27.5 // indirect
	k8s.io/apimachinery v0.28.6
	k8s.io/apiserver v0.28.4 // indirect
	k8s.io/component-base v0.28.4 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-aggregator v0.25.4 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	k8s.io/kubernetes v1.27.9 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	sigs.k8s.io/cli-utils v0.27.0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
