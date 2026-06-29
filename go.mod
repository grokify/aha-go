module github.com/grokify/aha-go

go 1.26.0

require (
	github.com/Khan/genqlient v0.8.1
	github.com/go-faster/errors v0.7.1
	github.com/go-faster/jx v1.2.0

	// Browser automation for template creation (no API available).
	// go-rod v0.116.2 (July 2024) requires fetchup v0.2.4; newer fetchup versions
	// have breaking API changes. go-rod appears unmaintained since 2024.
	// TODO: Consider migrating to chromedp with a high-level wrapper.
	// See: docs/specs/ROADMAP.md for migration plan.
	github.com/go-rod/rod v0.116.2
	github.com/grokify/goauth v0.23.30
	github.com/grokify/prism-roadmap v0.14.1
	github.com/ogen-go/ogen v1.22.0
	github.com/spf13/cobra v1.10.2
	github.com/ysmood/fetchup v0.2.4 // indirect; pinned for go-rod v0.116.2 compatibility
	go.opentelemetry.io/otel v1.44.0
	go.opentelemetry.io/otel/metric v1.44.0
	go.opentelemetry.io/otel/trace v1.44.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go/auth v0.20.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dlclark/regexp2 v1.12.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.16 // indirect
	github.com/googleapis/gax-go/v2 v2.22.0 // indirect
	github.com/grokify/mogo v0.74.6 // indirect
	github.com/grokify/priority-frameworks v0.2.0 // indirect
	github.com/grokify/prism-core v0.3.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.6.1 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/vektah/gqlparser/v2 v2.5.35 // indirect
	github.com/ysmood/goob v0.4.0 // indirect
	github.com/ysmood/got v0.41.0 // indirect
	github.com/ysmood/gson v0.7.3 // indirect
	github.com/ysmood/leakless v0.9.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.69.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/api v0.285.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260618152121-87f3d3e198d3 // indirect
	google.golang.org/grpc v1.81.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
