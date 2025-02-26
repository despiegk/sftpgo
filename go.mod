module github.com/drakkan/sftpgo/v2

go 1.21

require (
	cloud.google.com/go/storage v1.33.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.8.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.2.0
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5
	github.com/alexedwards/argon2id v0.0.0-20230305115115-4b3c3280a736
	github.com/amoghe/go-crypt v0.0.0-20220222110647-20eada5f5964
	github.com/aws/aws-sdk-go-v2 v1.21.2
	github.com/aws/aws-sdk-go-v2/config v1.18.45
	github.com/aws/aws-sdk-go-v2/credentials v1.13.43
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.13
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.90
	github.com/aws/aws-sdk-go-v2/service/marketplacemetering v1.16.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.40.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.21.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.23.2
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/cockroachdb/cockroach-go/v2 v2.3.5
	github.com/coreos/go-oidc/v3 v3.6.0
	github.com/drakkan/webdav v0.0.0-20230227175313-32996838bcd8
	github.com/eikenb/pipeat v0.0.0-20210730190139-06b3e6902001
	github.com/fclairamb/ftpserverlib v0.22.0
	github.com/fclairamb/go-log v0.4.1
	github.com/go-acme/lego/v4 v4.14.2
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/jwtauth/v5 v5.1.1
	github.com/go-chi/render v1.0.3
	github.com/go-sql-driver/mysql v1.7.1
	github.com/golang/mock v1.6.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.3.1
	github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-plugin v1.5.2
	github.com/hashicorp/go-retryablehttp v0.7.4
	github.com/jackc/pgx/v5 v5.4.3
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/klauspost/compress v1.17.0
	github.com/lestrrat-go/jwx/v2 v2.0.13
	github.com/lithammer/shortuuid/v3 v3.0.7
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/mhale/smtpd v0.8.0
	github.com/minio/sio v0.3.1
	github.com/otiai10/copy v1.14.0
	github.com/pires/go-proxyproto v0.7.0
	github.com/pkg/sftp v1.13.6
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.17.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/cors v1.10.1
	github.com/rs/xid v1.5.0
	github.com/rs/zerolog v1.31.0
	github.com/sftpgo/sdk v0.1.6
	github.com/shirou/gopsutil/v3 v3.23.9
	github.com/spf13/afero v1.10.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.17.0
	github.com/stretchr/testify v1.8.4
	github.com/studio-b12/gowebdav v0.9.0
	github.com/subosito/gotenv v1.6.0
	github.com/unrolled/secure v1.13.0
	github.com/wagslane/go-password-validator v0.3.0
	github.com/wneessen/go-mail v0.4.1-0.20230823094700-0bd5390e370d
	github.com/yl2chen/cidranger v1.0.3-0.20210928021809-d1cb2c52f37a
	go.etcd.io/bbolt v1.3.7
	go.uber.org/automaxprocs v1.5.3
	gocloud.dev v0.34.0
	golang.org/x/crypto v0.26.0
	golang.org/x/net v0.25.0
	golang.org/x/oauth2 v0.13.0
	golang.org/x/sys v0.23.0
	golang.org/x/term v0.23.0
	golang.org/x/time v0.3.0
	google.golang.org/api v0.147.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	cloud.google.com/go v0.110.8 // indirect
	cloud.google.com/go/compute v1.23.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.43 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.45 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.17.3 // indirect
	github.com/aws/smithy-go v1.15.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.1 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.4 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/dns v1.1.56 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/grpc v1.58.3 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/jlaffaye/ftp => github.com/drakkan/ftp v0.0.0-20201114075148-9b9adce499a9
	github.com/robfig/cron/v3 => github.com/drakkan/cron/v3 v3.0.0-20230222140221-217a1e4d96c0
	golang.org/x/crypto => github.com/drakkan/crypto v0.0.0-20231010171726-2674169e524d
)
