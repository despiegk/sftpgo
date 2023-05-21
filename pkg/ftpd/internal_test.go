// Copyright (C) 2019-2022  Nicola Murino
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ftpd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/eikenb/pipeat"
	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/pires/go-proxyproto"
	"github.com/sftpgo/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drakkan/sftpgo/v2/pkg/common"
	"github.com/drakkan/sftpgo/v2/pkg/dataprovider"
	"github.com/drakkan/sftpgo/v2/pkg/vfs"
)

const (
	ftpsCert = `-----BEGIN CERTIFICATE-----
MIICHTCCAaKgAwIBAgIUHnqw7QnB1Bj9oUsNpdb+ZkFPOxMwCgYIKoZIzj0EAwIw
RTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDAyMDQwOTUzMDRaFw0zMDAyMDEw
OTUzMDRaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYD
VQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAARCjRMqJ85rzMC998X5z761nJ+xL3bkmGVqWvrJ51t5OxV0v25NsOgR82CA
NXUgvhVYs7vNFN+jxtb2aj6Xg+/2G/BNxkaFspIVCzgWkxiz7XE4lgUwX44FCXZM
3+JeUbKjUzBRMB0GA1UdDgQWBBRhLw+/o3+Z02MI/d4tmaMui9W16jAfBgNVHSME
GDAWgBRhLw+/o3+Z02MI/d4tmaMui9W16jAPBgNVHRMBAf8EBTADAQH/MAoGCCqG
SM49BAMCA2kAMGYCMQDqLt2lm8mE+tGgtjDmtFgdOcI72HSbRQ74D5rYTzgST1rY
/8wTi5xl8TiFUyLMUsICMQC5ViVxdXbhuG7gX6yEqSkMKZICHpO8hqFwOD/uaFVI
dV4vKmHUzwK/eIx+8Ay3neE=
-----END CERTIFICATE-----`
	ftpsKey = `-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDCfMNsN6miEE3rVyUPwElfiJSWaR5huPCzUenZOfJT04GAcQdWvEju3
UM2lmBLIXpGgBwYFK4EEACKhZANiAARCjRMqJ85rzMC998X5z761nJ+xL3bkmGVq
WvrJ51t5OxV0v25NsOgR82CANXUgvhVYs7vNFN+jxtb2aj6Xg+/2G/BNxkaFspIV
CzgWkxiz7XE4lgUwX44FCXZM3+JeUbI=
-----END EC PRIVATE KEY-----`
	caCRT = `-----BEGIN CERTIFICATE-----
MIIE5jCCAs6gAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhDZXJ0
QXV0aDAeFw0yMzAxMDMxMDIwNDdaFw0zMzAxMDMxMDMwNDZaMBMxETAPBgNVBAMT
CENlcnRBdXRoMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxq6Wl1Ih
hgvGdM8M2IVI7dwnv3yShJygZsnREQSEW0xeWJL5DtNeHCME5WByFUAlZKpePtW8
TNwln9DYDtgNSMiWwvO/wR0mXsyU8Ma4ZBMlX0oOkWo1Ff/M/u8YY9X78Vvwdt62
Yt7QmU5oUUW2HdAgh4AlhKJSjm3t0uDP5s54uvueL5bjChHwEb1ZGOtST9Zt86cj
YA/xtVHnDXCJbhohpzQI6dK96NegONZVDaxEohVCyYYOgI1I14Bxu0ZCMm5GjwoO
QohnUfEJ+BRgZqFpbsnYCE+PoVayVVFoLA+GMeqbQ2SHej1Pr1K0dbjUz6SAk8/+
DL7h8d+YAtflATsMtdsVJ4WzEfvZbSbiYKYmlVC6zk6ooXWadvQ5+aezVes9WMpH
YnAoScuKoeerRuKlsSU7u+XCmy/i7Hii5FwMrSvIL2GLtVE+tJFCTABA55OWZikt
ULMQfg3P2Hk3GFIE35M10mSjKQkGhz06WC5UQ7f2Xl9GzO6PqRSorzugewgMK6L4
SnN7XBFnUHHqx1bWNkUG8NPYB6Zs7UDHygemTWxqqxun43s501DNTSunCKIhwFbt
1ol5gOvYAFG+BXxnggBT815Mgz1Zht3S9CuprAgz0grNEwAYjRTm1PSaX3t8I1kv
oUUuMF6BzWLHJ66uZKOCsPs3ouGq+G3GfWUCAwEAAaNFMEMwDgYDVR0PAQH/BAQD
AgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFCj8lmcR7loB9zAP/feo
t1eIHWmIMA0GCSqGSIb3DQEBCwUAA4ICAQCu46fF0Tr2tZz1wkYt2Ty3OU77jcG9
zYU/7uxPqPC8OYIzQJrumXKLOkTYJXJ7k+7RQdsn/nbxdH1PbslNDD3Eid/sZsF/
dFdGR1ZYwXVQbpYwEd19CvJTALn9CyAZpMS8J2RJrmdScAeSSb0+nAGTYP7GvPm+
8ktOnrz3w8FtzTw+seuCW/DI/5UpfC9Jf+i/3XgxDozXWNW6YNOIw/CicyaqbBTk
5WFcJ0WJN+8qQurw8n+sOvQcNsuDTO7K3Tqu0wGTDUQKou7kiMX0UISRvd8roNOl
zvvokNQe4VgCGQA+Y2SxvSxVG1BaymYeNw/0Yxm7QiKSUI400V1iKIcpnIvIedJR
j2bGIlslVSV/P6zkRuF1srRVxTxSf1imEfs8J8mMhHB6DkOsP4Y93z5s6JZ0sPiM
eOb0CVKul/e1R0Kq23AdPf5eUv63RhfmokN1OsdarRKMFyHphWMxqGJXsSvRP+dl
3DaKeTDx/91OSWiMc+glHHKKJveMYQLeJ7GXmcxhuoBm6o4Coowgw8NFKMCtAsp0
ktvsQuhB3uFUterw/2ONsOChx7Ybu36Zk47TKBpktfxDQ578TVoZ7xWSAFqCPHvx
A5VSwAg7tdBvORfqQjhiJRnhwr50RaNQABTLS0l5Vsn2mitApPs7iKiIts2ieWsU
EsdgvPZR2e5IkA==
-----END CERTIFICATE-----`
	caKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	caCRL = `-----BEGIN X509 CRL-----
MIICpjCBjwIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhDZXJ0QXV0aBcN
MjMwMTAzMTAzMzI4WhcNMjUwMTAyMTAzMzI4WjAjMCECEHUfHtKUGlg/86yMN/aM
xxsXDTIzMDEwMzEwMzMyOFqgIzAhMB8GA1UdIwQYMBaAFCj8lmcR7loB9zAP/feo
t1eIHWmIMA0GCSqGSIb3DQEBCwUAA4ICAQAJf6MBMUc3xWKB6fy0VoPbXQjVTsL4
Yjm5lKaCtvcRiJ6onaITfJL6V3OCy/MAe94sHynvK3DyyYvxJ0ms7y+kmEtFzHwz
T+hBPHaEV/Ccamt+3zRZwndwEMomkQz5tBipwimOlsYXWqItjhXHcLLr84jWgqpD
JHcfDmLswCeJVqe8xyYSYCnWMjQ3sn0h+arjm53SdHTULlsjgKeX/ao2IJwt1Ddr
APYKZ/XBWq9vBq3l4l2Ufj16fUBY5NeHTjQcLLrkwmBwpSb0YS8+jBwmOwo1HwEF
MEwADBTHI2jT4ygzzKefVETfcSk4CuIQ1ve0qQL7KY5Fg5AXwbRycev6R0vEHR82
oOPAqg+dYgKtdkxK5QZrNLenloq6x0/3oEThwOg3J17+eCYjixBC+3PoUzLa+yfZ
xSQ/kkcRJExEhadw5I9TI7sEUk1RjDCl6AtHg53LQifokiLLfMRptOiN3a4NlLJ2
HNXfWUltRUnr6MCxk+G7U5Zaj1QtCN3Aldw+3xcJr7FOBU23VqRT22XbfW+B1gsr
4eNlw5Kk/PDF/WZeGbrwy7fvpSoFsDYI8lpVlzKVwLampIZVhnWsfaf7jk/pc4T0
6iZ+rnB6CP4P+LM34jKYiAtz+iufjEB6Ko0jN0ZWCznDGDVgMwnVynGJNKU+4bl8
vC4nIbS2OhclcA==
-----END X509 CRL-----`
	client1Crt = `-----BEGIN CERTIFICATE-----
MIIEIDCCAgigAwIBAgIQWwKNgBzKWM8ewyS4K78uOTANBgkqhkiG9w0BAQsFADAT
MREwDwYDVQQDEwhDZXJ0QXV0aDAeFw0yMzAxMDMxMDIzMTFaFw0zMzAxMDMxMDMw
NDVaMBIxEDAOBgNVBAMTB2NsaWVudDEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQC+jwuaG0FSpPtE5mZPBgdacAhXXa51/TIT18HTm+QnOYUcGRel3AuZ
OBWv3fOallW8iQX3i+M78cKeTMWOS5RGXCdDe866pYXEyUkZFRRSA/6573Dz5dJ/
DZOCsgW+91JlSkM1+FYE9cpt4qLkdAjSRXIoebcA64K60wqZr1Js+pQrH3leT9wu
33dM3KHkDHOeMj6X/V1me22htndD/DUlWmPc58jMFbcvxFG3oUBB9U65LJBwJNzr
XWVcli2QirZ0fLkC7Lo2FIYuN1qeU/8A/T4TTInZb/eW3Faqv4RuhjWPXFLqkdIP
4AzDxCNuhlWqyv9nfgegXAHOHpXZMDKxAgMBAAGjcTBvMA4GA1UdDwEB/wQEAwID
uDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYEFKloKnzI
w7YYnjm1sKU+LgvT5dU0MB8GA1UdIwQYMBaAFCj8lmcR7loB9zAP/feot1eIHWmI
MA0GCSqGSIb3DQEBCwUAA4ICAQAeja0rK7I14ibgis9vSPXAGmmKpagIjvZiYCE6
Ti/Rq6qbyQ6tKL08NxR2XPNjoXfxwGOGgboWR86S7WT93pz3HkftAjTfzUnxnXOx
S7dWfq+g0uY/3ql6IFDQBpGKHu/KN8/1Pvn39FiYSdCaM66bwyFukcvBXace+aC1
M6jzVsscxoCCjXhcZl++Tjpf6TzGMd8OFyArBQmOUCoFrTcMzLPKSAROAHp0k+Ju
HHNgLdgXPQCfAgRbWnqq2o2moApe7+gzMS+1X0eKhIXYS7csK8rFvGzjH/ANDo9A
8+AFcJh8YiIlEVI8nCb3ERdpAbu6G5xkfUDkqcWjCAhuokrbeFmU82RQOc3TQZc9
NMHfTkCOPhaIdPI/kC+fZkdz+5ftDCl/radSljeMX+/y0DVQUOtrQzyT1PBN0vCx
L+FCzj0fHJwdoDiFLxDLLN1pYWsxMnIichpg625CZM9r5i183yPErXxxQPydcDrX
Y6Ps7rGiU7eGILhAfQnS1XUDvH0gNfLUvO5uWm6yO4yUEDWkA/wOTnrc8Z5Waza+
kH+FmxnYpT1rMloTSoyiHIPvTq1nVJ8LILUODZAxW+ZHmccGgOpIN/DWuWunVRHG
tuaTSgU1xjWl2q/SeoS2DpiEKTIAZZQ5CTD819oc8SnLTzK0ISRpBXKg13AF2uJD
G9q7sA==
-----END CERTIFICATE-----`
	client1Key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAvo8LmhtBUqT7ROZmTwYHWnAIV12udf0yE9fB05vkJzmFHBkX
pdwLmTgVr93zmpZVvIkF94vjO/HCnkzFjkuURlwnQ3vOuqWFxMlJGRUUUgP+ue9w
8+XSfw2TgrIFvvdSZUpDNfhWBPXKbeKi5HQI0kVyKHm3AOuCutMKma9SbPqUKx95
Xk/cLt93TNyh5AxznjI+l/1dZnttobZ3Q/w1JVpj3OfIzBW3L8RRt6FAQfVOuSyQ
cCTc611lXJYtkIq2dHy5Auy6NhSGLjdanlP/AP0+E0yJ2W/3ltxWqr+EboY1j1xS
6pHSD+AMw8QjboZVqsr/Z34HoFwBzh6V2TAysQIDAQABAoIBAFaIHnycY81jnbZr
6Yl4813eAeuqXs61a0gXcazl3XTyab+YpWRrx9iL3009PKG2Iri6gDspCsbtwbKg
qhUzvOE2d53tWrLm9xelT8xUBiY4KjPEx0X51txbDeELdhCBvqjAUETxwB4Afyvm
/pE/H8JcRrqair+gMn0j2GxxcLyLQt8/DBaqbs50QDxYbLTrfZzXi3R5iAMmtGDM
ZuhBDJYjw/PdJnmWcCkeEFa731ZwHISvDFJtZ6kv0yU7guHzvDWOFlszFksv8HRI
s46i1AqvdLd3M/xVDWi2f5P3IuOK80v2xrTZAbJSc9Fo/oHhO+9mWoxnGF2JE2zO
cabYfAECgYEA/EIw0fvOLabhmsLTItq7p76Gt1kE2Rsn+KsRH+H4vE3iptHy1pks
j/aohQ+YeZM3KtsEz12dtPfUBQoTMfCxpiHVhhpX5JLyc6OAMWhZUItQX2X0lzeY
oPRdbEMRcfxOKjb3mY5T2h9TVUieuspE2wExYFRMaB8BT4pio86iWYUCgYEAwWKV
pP7w1+H6bpBucSh89Iq0inIgvHpFNz0bpAFTZV+VyydMzjfkY8k6IqL5ckr2aDsY
vk6XLClJi6I2qeQx/czIrd+xCWcSJPLTcjtKwv0T01ThNVq+ev1NBUqU03STyaJa
p14r4dIYxpZs13s+Mdkzr7R8uv4J5Y03AP90xj0CgYEA4j0W/ezBAE6QPbWHmNXl
wU7uEZgj8fcaBTqfVCHdbDzKDuVyzqZ3wfHtN9FB5Z9ztdrSWIxUec5e99oOVxbQ
rPfhQbF0rIpiKfY0bZtxpvwbLEQLdmelWo1vED6iccFf9RpxO+XbLGA14+IKgeoQ
kP5j40oXcLaF/WlWiCU1k+UCgYEAgVFcgn5dLfAWmLMKt678iEbs3hvdmkwlVwAN
KMoeK48Uy0pXiRtFJhldP+Y96tkIF8FVFYXWf5iIbtClv0wyxeaYV/VbHM+JCZ48
GYpevy+ff1WmWBh7giE6zQwHo7O0VES2XG+T5qmpGbtjw2DNwWXes2N9eUoB8jhR
jOBHBX0CgYEA6Ha3IdnpYyODII1W26gEPnBoUCk1ascsztAqDwikBgMY9605jxLi
t3L261iTtN4kTd26nPTsNaJlEnKfm7Oqg1P3zpYLmC2JoFVrOyAZVhyfLACBMV9g
Dy1qoA4qz5jjtwPQ0bsOpfE6/oXdIZZdgyi1CmVRMNF0z3KNs1LhLLU=
-----END RSA PRIVATE KEY-----`
	// client 2 crt is revoked
	client2Crt = `-----BEGIN CERTIFICATE-----
MIIEIDCCAgigAwIBAgIQdR8e0pQaWD/zrIw39ozHGzANBgkqhkiG9w0BAQsFADAT
MREwDwYDVQQDEwhDZXJ0QXV0aDAeFw0yMzAxMDMxMDIzMTRaFw0zMzAxMDMxMDMw
NDVaMBIxEDAOBgNVBAMTB2NsaWVudDIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQC/UZqUxeP15lhXBmPmpS5SdI470R75fxmN14FhYwTS3FsDoT+evqRg
II4Qo/wqbaGrk/BsbzB7ToVWqpkyZ58hYPdtLjKtBBHYsSCNCoKZEVJTz5JdW3sj
CKRsG3zPVhFjJcYW9pKsr/CGIIDWAfkuuwR+R/NHkUFSjEP5N9qMAc9wBvskxV84
YAJJykPD9rG8PjXHOKsfNUhH+/QfbqMkCeETJ1sp66o3ilql2aZ0m6K6x4gB7tM7
NZnM4eztLZbAnQVQhNBYCR6i7DGI2dujujPbpCqmSqSb42n+3a2o844k6EnU76HJ
RZwhd3ypy9CvTdkya5JbK+aKKo8fGFHbAgMBAAGjcTBvMA4GA1UdDwEB/wQEAwID
uDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYEFLItEDE1
gVYfe7JSax5YAjEW8tmzMB8GA1UdIwQYMBaAFCj8lmcR7loB9zAP/feot1eIHWmI
MA0GCSqGSIb3DQEBCwUAA4ICAQCZRStHCbwmhmH4tu7V5ammmQhn1TKcspV86cXz
JQ4ZM11RvGpRLTmYuRkl5XloMuvB8yYAE1ihhkYOhgU6zSAj33kUQSx6cHXWau7T
NjTchptKX+b17GR/yuFwIR3TugArBsnwyuUdts478gTY+MSgTOWWyOgWl3FujiEJ
GJ7EgKde4jURXv2qjp6ZtSVqMlAa3y8C3S8nLdyt9Rf8CcSjEy/t8t0JhoMYCvxg
o1k7QhMCfMYjjEIuEyDVOdCs2ExepG1zUVBP5h5239sXvLKrOZvgCZkslyTNd/m9
vv4yR5gLgCdt0Ol1uip0p910PJoSqX6nZNfeCx3+Kgyc7crl8PrsnUAVoPgLxpVm
FWF+KlUbh2KiYTuSi5cH0Ti9NtWT3Qi8d4WhmjFlu0zD3EJEUie0oYfHERiO9bVo
5EAzERSVhgQdxVOLgIc2Hbe1JYFf7idyqASRw6KdVkW6YIC/V/5efrJ1LZ5QNrdv
bmfJ5CznE6o1AH9JsQ8xMi+kmyn/It1uMWIwP/tYyjQ98dlOj2k9CHP2RzrvCCY9
yZNjs2QC5cleNdSpNMb2J2EUYTNAnaH3H8YdbT0scMHDvre1G7p4AjeuRJ9mW7VK
Dcqbw+VdSAPyAFdiCd9x8AU3sr28vYbPbPp+LsHQXIlYdnVR0hh2UKF5lR8iuqVx
y05cAQ==
-----END CERTIFICATE-----`
	client2Key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAv1GalMXj9eZYVwZj5qUuUnSOO9Ee+X8ZjdeBYWME0txbA6E/
nr6kYCCOEKP8Km2hq5PwbG8we06FVqqZMmefIWD3bS4yrQQR2LEgjQqCmRFSU8+S
XVt7IwikbBt8z1YRYyXGFvaSrK/whiCA1gH5LrsEfkfzR5FBUoxD+TfajAHPcAb7
JMVfOGACScpDw/axvD41xzirHzVIR/v0H26jJAnhEydbKeuqN4papdmmdJuiuseI
Ae7TOzWZzOHs7S2WwJ0FUITQWAkeouwxiNnbo7oz26Qqpkqkm+Np/t2tqPOOJOhJ
1O+hyUWcIXd8qcvQr03ZMmuSWyvmiiqPHxhR2wIDAQABAoIBAQCGAtE2uM8PJcRn
YPCFVNr3ovEmcTszJJZvxq632rY8RWHzTvXTalKVivg4K8WsqpJ+LuhP7CqXlM7N
gD5DElZi+RsXfS6+BoXBtYDJir0kHv/9+P3bKwM77QfPOgnY6b7QJlt1Jk5ja/Ic
4ZOdVFCJLTLeieOdE+AfxGSwozEQs9N3wBjPi6i5Rarc6i8HbuSemp/KfXrSR/Sh
EFajk0l3nFVgr3VOLGsV/ieT6EW42p6ZA1ZBEi4sr4hN49zU2Vpj+lXBl/RhVGgM
6cSYJkOP98eD2t9cjHyZFqSw18/UqTNonMfoT2uvSNni9/jAouzkt7SwaPAqQpjE
BfiJnK9RAoGBAMNdK6AhS+9ouQEP5v+9ubQ3AIEMYb+3uVR+1veCBqq49J0Z3tgk
7Ts5eflsYnddmtys+8CMnAvh+1EedK+X/MQyalQAUHke+llt94N+tHpSPDw/ZHOy
koyLFg6efQr+626x6o33jqu+/9fu7Szxv41tmnCfh9hxGXda3aiWHsUdAoGBAPqz
BQVWI7NOJmsiSB0OoDs+x3fqjp31IxEw63t+lDtSTTzxCU53sfie9vPsKdPFqR+d
yNa5i5M8YDuaEjbN3hpuOpRWbfg2aPVyx3TNPp8bHNNUuJkCQ4Z2b0Imlv4Sycl+
CCMMXvysIAomxkAZ3Q3BsSAZd2n+qvLvMt2jGZlXAoGAa/AhN1LOMpMojBauKSQ4
4wH0jFg79YHbqnx95rf3WQHhXJ87iS41yCAEbTNd39dexYfpfEPzv3j2sqXiEFYn
+HpmVszpqVHdPeXM9+DcdCzVTPA1XtsNrwr1f9Q/AAFCMKGqFw/syqU3k6VVcxyK
GeixiIILuyEZ0eDpUMjIbV0CgYBbwvLvhRwEIXLGfAHRQO09QjlYly4kevme7T0E
Msym+fTzfXZelkk6K1VQ6vxUW2EQBXzhu4BvIAZJSpeoH6pQGlCuwwP1elTool6H
TijBq/bdE4GN39o/eVI38FAMJ2xcqBjqWzjZW1dO3+poxA65XlAq46dl0KVZzlvb
7DsOeQKBgQCW8iELrECLQ9xhPbzqdNEOcI4wxEI8oDNLvUar/VnMrSUBxi/jo3j2
08IOKMKqSl+BX77ftgazhyL+hEgxlZuPKeqUuOWcNxuAs0vK6Gc5+Y9UpQEq78nH
uaPG3o9EBDf5eFKi76o+pVtqxrwhY88M/Yw0ykEA6Nf7RCo2ucdemg==
-----END RSA PRIVATE KEY-----`
)

var (
	configDir = filepath.Join(".", "..", "..")
)

type mockFTPClientContext struct {
	lastDataChannel ftpserver.DataChannel
	remoteIP        string
	localIP         string
}

func (cc mockFTPClientContext) Path() string {
	return ""
}

func (cc mockFTPClientContext) SetPath(name string) {}

func (cc mockFTPClientContext) SetListPath(name string) {}

func (cc mockFTPClientContext) SetDebug(debug bool) {}

func (cc mockFTPClientContext) Debug() bool {
	return false
}

func (cc mockFTPClientContext) ID() uint32 {
	return 1
}

func (cc mockFTPClientContext) RemoteAddr() net.Addr {
	ip := "127.0.0.1"
	if cc.remoteIP != "" {
		ip = cc.remoteIP
	}
	return &net.IPAddr{IP: net.ParseIP(ip)}
}

func (cc mockFTPClientContext) LocalAddr() net.Addr {
	ip := "127.0.0.1"
	if cc.localIP != "" {
		ip = cc.localIP
	}
	return &net.IPAddr{IP: net.ParseIP(ip)}
}

func (cc mockFTPClientContext) GetClientVersion() string {
	return "mock version"
}

func (cc mockFTPClientContext) Close() error {
	return nil
}

func (cc mockFTPClientContext) HasTLSForControl() bool {
	return false
}

func (cc mockFTPClientContext) HasTLSForTransfers() bool {
	return false
}

func (cc mockFTPClientContext) SetTLSRequirement(requirement ftpserver.TLSRequirement) error {
	return nil
}

func (cc mockFTPClientContext) GetLastCommand() string {
	return ""
}

func (cc mockFTPClientContext) GetLastDataChannel() ftpserver.DataChannel {
	return cc.lastDataChannel
}

// MockOsFs mockable OsFs
type MockOsFs struct {
	vfs.Fs
	err                     error
	statErr                 error
	isAtomicUploadSupported bool
}

// Name returns the name for the Fs implementation
func (fs MockOsFs) Name() string {
	return "mockOsFs"
}

// IsUploadResumeSupported returns true if resuming uploads is supported
func (MockOsFs) IsUploadResumeSupported() bool {
	return false
}

// IsAtomicUploadSupported returns true if atomic upload is supported
func (fs MockOsFs) IsAtomicUploadSupported() bool {
	return fs.isAtomicUploadSupported
}

// Stat returns a FileInfo describing the named file
func (fs MockOsFs) Stat(name string) (os.FileInfo, error) {
	if fs.statErr != nil {
		return nil, fs.statErr
	}
	return os.Stat(name)
}

// Lstat returns a FileInfo describing the named file
func (fs MockOsFs) Lstat(name string) (os.FileInfo, error) {
	if fs.statErr != nil {
		return nil, fs.statErr
	}
	return os.Lstat(name)
}

// Remove removes the named file or (empty) directory.
func (fs MockOsFs) Remove(name string, isDir bool) error {
	if fs.err != nil {
		return fs.err
	}
	return os.Remove(name)
}

// Rename renames (moves) source to target
func (fs MockOsFs) Rename(source, target string) error {
	if fs.err != nil {
		return fs.err
	}
	return os.Rename(source, target)
}

func newMockOsFs(err, statErr error, atomicUpload bool, connectionID, rootDir string) vfs.Fs {
	return &MockOsFs{
		Fs:                      vfs.NewOsFs(connectionID, rootDir, ""),
		err:                     err,
		statErr:                 statErr,
		isAtomicUploadSupported: atomicUpload,
	}
}

func TestInitialization(t *testing.T) {
	oldMgr := certMgr
	certMgr = nil

	binding := Binding{
		Port: 2121,
	}
	c := &Configuration{
		Bindings:           []Binding{binding},
		CertificateFile:    "acert",
		CertificateKeyFile: "akey",
	}
	assert.False(t, binding.HasProxy())
	assert.Equal(t, "Disabled", binding.GetTLSDescription())
	err := c.Initialize(configDir)
	assert.Error(t, err)
	c.CertificateFile = ""
	c.CertificateKeyFile = ""
	c.BannerFile = "afile"
	server := NewServer(c, configDir, binding, 0)
	assert.Equal(t, "", server.initialMsg)
	_, err = server.GetTLSConfig()
	assert.Error(t, err)

	binding.TLSMode = 1
	server = NewServer(c, configDir, binding, 0)
	_, err = server.GetSettings()
	assert.Error(t, err)

	binding.PassiveConnectionsSecurity = 100
	binding.ActiveConnectionsSecurity = 100
	server = NewServer(c, configDir, binding, 0)
	_, err = server.GetSettings()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid passive_connections_security")
	}
	binding.PassiveConnectionsSecurity = 1
	server = NewServer(c, configDir, binding, 0)
	_, err = server.GetSettings()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid active_connections_security")
	}
	binding = Binding{
		Port:           2121,
		ForcePassiveIP: "192.168.1",
	}
	server = NewServer(c, configDir, binding, 0)
	_, err = server.GetSettings()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not valid")
	}

	binding.ForcePassiveIP = "::ffff:192.168.89.9"
	err = binding.checkPassiveIP()
	assert.NoError(t, err)
	assert.Equal(t, "192.168.89.9", binding.ForcePassiveIP)

	binding.ForcePassiveIP = "::1"
	err = binding.checkPassiveIP()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not a valid IPv4 address")
	}

	err = ReloadCertificateMgr()
	assert.NoError(t, err)

	certMgr = oldMgr

	binding = Binding{
		Port:           2121,
		ClientAuthType: 1,
	}
	server = NewServer(c, configDir, binding, 0)
	cfg, err := server.GetTLSConfig()
	assert.NoError(t, err)
	assert.Equal(t, tls.RequireAndVerifyClientCert, cfg.ClientAuth)
}

func TestServerGetSettings(t *testing.T) {
	oldConfig := common.Config

	binding := Binding{
		Port:             2121,
		ApplyProxyConfig: true,
	}
	c := &Configuration{
		Bindings: []Binding{binding},
		PassivePortRange: PortRange{
			Start: 10000,
			End:   11000,
		},
	}
	assert.False(t, binding.HasProxy())
	server := NewServer(c, configDir, binding, 0)
	settings, err := server.GetSettings()
	assert.NoError(t, err)
	assert.Equal(t, 10000, settings.PassiveTransferPortRange.Start)
	assert.Equal(t, 11000, settings.PassiveTransferPortRange.End)

	common.Config.ProxyProtocol = 1
	common.Config.ProxyAllowed = []string{"invalid"}
	assert.True(t, binding.HasProxy())
	_, err = server.GetSettings()
	assert.Error(t, err)
	server.binding.Port = 8021
	_, err = server.GetSettings()
	assert.Error(t, err)

	assert.Equal(t, "Plain and explicit", binding.GetTLSDescription())

	binding.TLSMode = 1
	assert.Equal(t, "Explicit required", binding.GetTLSDescription())

	binding.TLSMode = 2
	assert.Equal(t, "Implicit", binding.GetTLSDescription())

	certPath := filepath.Join(os.TempDir(), "test.crt")
	keyPath := filepath.Join(os.TempDir(), "test.key")
	err = os.WriteFile(certPath, []byte(ftpsCert), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(keyPath, []byte(ftpsKey), os.ModePerm)
	assert.NoError(t, err)

	common.Config.ProxyAllowed = nil
	c.CertificateFile = certPath
	c.CertificateKeyFile = keyPath
	server = NewServer(c, configDir, binding, 0)
	server.binding.Port = 9021
	settings, err = server.GetSettings()
	assert.NoError(t, err)
	assert.NotNil(t, settings.Listener)

	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	listener, err = server.WrapPassiveListener(listener)
	assert.NoError(t, err)

	_, ok := listener.(*proxyproto.Listener)
	assert.True(t, ok)

	err = os.Remove(certPath)
	assert.NoError(t, err)
	err = os.Remove(keyPath)
	assert.NoError(t, err)

	common.Config = oldConfig
}

func TestUserInvalidParams(t *testing.T) {
	u := dataprovider.User{
		BaseUser: sdk.BaseUser{
			HomeDir: "invalid",
		},
	}
	binding := Binding{
		Port: 2121,
	}
	c := &Configuration{
		Bindings: []Binding{binding},
		PassivePortRange: PortRange{
			Start: 10000,
			End:   11000,
		},
	}
	server := NewServer(c, configDir, binding, 3)
	_, err := server.validateUser(u, mockFTPClientContext{}, dataprovider.LoginMethodPassword)
	assert.Error(t, err)

	u.Username = "a"
	u.HomeDir = filepath.Clean(os.TempDir())
	subDir := "subdir"
	mappedPath1 := filepath.Join(os.TempDir(), "vdir1")
	vdirPath1 := "/vdir1"
	mappedPath2 := filepath.Join(os.TempDir(), "vdir1", subDir)
	vdirPath2 := "/vdir2"
	u.VirtualFolders = append(u.VirtualFolders, vfs.VirtualFolder{
		BaseVirtualFolder: vfs.BaseVirtualFolder{
			MappedPath: mappedPath1,
		},
		VirtualPath: vdirPath1,
	})
	u.VirtualFolders = append(u.VirtualFolders, vfs.VirtualFolder{
		BaseVirtualFolder: vfs.BaseVirtualFolder{
			MappedPath: mappedPath2,
		},
		VirtualPath: vdirPath2,
	})
	_, err = server.validateUser(u, mockFTPClientContext{}, dataprovider.LoginMethodPassword)
	assert.Error(t, err)
	u.VirtualFolders = nil
	_, err = server.validateUser(u, mockFTPClientContext{}, dataprovider.LoginMethodPassword)
	assert.Error(t, err)
}

func TestFTPMode(t *testing.T) {
	connection := &Connection{
		BaseConnection: common.NewBaseConnection("", common.ProtocolFTP, "", "", dataprovider.User{}),
	}
	assert.Empty(t, connection.getFTPMode())
	connection.clientContext = mockFTPClientContext{lastDataChannel: ftpserver.DataChannelActive}
	assert.Equal(t, "active", connection.getFTPMode())
	connection.clientContext = mockFTPClientContext{lastDataChannel: ftpserver.DataChannelPassive}
	assert.Equal(t, "passive", connection.getFTPMode())
	connection.clientContext = mockFTPClientContext{lastDataChannel: 0}
	assert.Empty(t, connection.getFTPMode())
}

func TestClientVersion(t *testing.T) {
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("2_%v", mockCC.ID())
	user := dataprovider.User{}
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	err := common.Connections.Add(connection)
	assert.NoError(t, err)
	stats := common.Connections.GetStats()
	if assert.Len(t, stats, 1) {
		assert.Equal(t, "mock version", stats[0].ClientVersion)
		common.Connections.Remove(connection.GetID())
	}
	assert.Len(t, common.Connections.GetStats(), 0)
}

func TestDriverMethodsNotImplemented(t *testing.T) {
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("2_%v", mockCC.ID())
	user := dataprovider.User{}
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	_, err := connection.Create("")
	assert.EqualError(t, err, errNotImplemented.Error())
	err = connection.MkdirAll("", os.ModePerm)
	assert.EqualError(t, err, errNotImplemented.Error())
	_, err = connection.Open("")
	assert.EqualError(t, err, errNotImplemented.Error())
	_, err = connection.OpenFile("", 0, os.ModePerm)
	assert.EqualError(t, err, errNotImplemented.Error())
	err = connection.RemoveAll("")
	assert.EqualError(t, err, errNotImplemented.Error())
	assert.Equal(t, connection.GetID(), connection.Name())
}

func TestResolvePathErrors(t *testing.T) {
	user := dataprovider.User{
		BaseUser: sdk.BaseUser{
			HomeDir: "invalid",
		},
	}
	user.Permissions = make(map[string][]string)
	user.Permissions["/"] = []string{dataprovider.PermAny}
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("%v", mockCC.ID())
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	err := connection.Mkdir("", os.ModePerm)
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	err = connection.Remove("")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	err = connection.RemoveDir("")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	err = connection.Rename("", "")
	assert.ErrorIs(t, err, common.ErrOpUnsupported)
	err = connection.Symlink("", "")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	_, err = connection.Stat("")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	err = connection.Chmod("", os.ModePerm)
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	err = connection.Chtimes("", time.Now(), time.Now())
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	_, err = connection.ReadDir("")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	_, err = connection.GetHandle("", 0, 0)
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
	_, err = connection.GetAvailableSpace("")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrGenericFailure.Error())
	}
}

func TestUploadFileStatError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("this test is not available on Windows")
	}
	user := dataprovider.User{
		BaseUser: sdk.BaseUser{
			Username: "user",
			HomeDir:  filepath.Clean(os.TempDir()),
		},
	}
	user.Permissions = make(map[string][]string)
	user.Permissions["/"] = []string{dataprovider.PermAny}
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("%v", mockCC.ID())
	fs := vfs.NewOsFs(connID, user.HomeDir, "")
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	testFile := filepath.Join(user.HomeDir, "test", "testfile")
	err := os.MkdirAll(filepath.Dir(testFile), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(testFile, []byte("data"), os.ModePerm)
	assert.NoError(t, err)
	err = os.Chmod(filepath.Dir(testFile), 0001)
	assert.NoError(t, err)
	_, err = connection.uploadFile(fs, testFile, "test", 0)
	assert.Error(t, err)
	err = os.Chmod(filepath.Dir(testFile), os.ModePerm)
	assert.NoError(t, err)
	err = os.RemoveAll(filepath.Dir(testFile))
	assert.NoError(t, err)
}

func TestAVBLErrors(t *testing.T) {
	user := dataprovider.User{
		BaseUser: sdk.BaseUser{
			Username: "user",
			HomeDir:  filepath.Clean(os.TempDir()),
		},
	}
	user.Permissions = make(map[string][]string)
	user.Permissions["/"] = []string{dataprovider.PermAny}
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("%v", mockCC.ID())
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	_, err := connection.GetAvailableSpace("/")
	assert.NoError(t, err)
	_, err = connection.GetAvailableSpace("/missing-path")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, fs.ErrNotExist))
}

func TestUploadOverwriteErrors(t *testing.T) {
	user := dataprovider.User{
		BaseUser: sdk.BaseUser{
			Username: "user",
			HomeDir:  filepath.Clean(os.TempDir()),
		},
	}
	user.Permissions = make(map[string][]string)
	user.Permissions["/"] = []string{dataprovider.PermAny}
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("%v", mockCC.ID())
	fs := newMockOsFs(nil, nil, false, connID, user.GetHomeDir())
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	flags := 0
	flags |= os.O_APPEND
	_, err := connection.handleFTPUploadToExistingFile(fs, flags, "", "", 0, "")
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrOpUnsupported.Error())
	}

	f, err := os.CreateTemp("", "temp")
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	flags = 0
	flags |= os.O_CREATE
	flags |= os.O_TRUNC
	tr, err := connection.handleFTPUploadToExistingFile(fs, flags, f.Name(), f.Name(), 123, f.Name())
	if assert.NoError(t, err) {
		transfer := tr.(*transfer)
		transfers := connection.GetTransfers()
		if assert.Equal(t, 1, len(transfers)) {
			assert.Equal(t, transfers[0].ID, transfer.GetID())
			assert.Equal(t, int64(123), transfer.InitialSize)
			err = transfer.Close()
			assert.NoError(t, err)
			assert.Equal(t, 0, len(connection.GetTransfers()))
		}
	}
	err = os.Remove(f.Name())
	assert.NoError(t, err)

	_, err = connection.handleFTPUploadToExistingFile(fs, os.O_TRUNC, filepath.Join(os.TempDir(), "sub", "file"),
		filepath.Join(os.TempDir(), "sub", "file1"), 0, "/sub/file1")
	assert.Error(t, err)
	fs = vfs.NewOsFs(connID, user.GetHomeDir(), "")
	_, err = connection.handleFTPUploadToExistingFile(fs, 0, "missing1", "missing2", 0, "missing")
	assert.Error(t, err)
}

func TestTransferErrors(t *testing.T) {
	testfile := "testfile"
	file, err := os.Create(testfile)
	assert.NoError(t, err)
	user := dataprovider.User{
		BaseUser: sdk.BaseUser{
			Username: "user",
			HomeDir:  filepath.Clean(os.TempDir()),
		},
	}
	user.Permissions = make(map[string][]string)
	user.Permissions["/"] = []string{dataprovider.PermAny}
	mockCC := mockFTPClientContext{}
	connID := fmt.Sprintf("%v", mockCC.ID())
	fs := newMockOsFs(nil, nil, false, connID, user.GetHomeDir())
	connection := &Connection{
		BaseConnection: common.NewBaseConnection(connID, common.ProtocolFTP, "", "", user),
		clientContext:  mockCC,
	}
	baseTransfer := common.NewBaseTransfer(file, connection.BaseConnection, nil, file.Name(), file.Name(), testfile,
		common.TransferDownload, 0, 0, 0, 0, false, fs, dataprovider.TransferQuota{})
	tr := newTransfer(baseTransfer, nil, nil, 0)
	err = tr.Close()
	assert.NoError(t, err)
	_, err = tr.Seek(10, 0)
	assert.Error(t, err)
	buf := make([]byte, 64)
	_, err = tr.Read(buf)
	assert.Error(t, err)
	err = tr.Close()
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrTransferClosed.Error())
	}
	assert.Len(t, connection.GetTransfers(), 0)

	r, _, err := pipeat.Pipe()
	assert.NoError(t, err)
	baseTransfer = common.NewBaseTransfer(nil, connection.BaseConnection, nil, testfile, testfile, testfile,
		common.TransferUpload, 0, 0, 0, 0, false, fs, dataprovider.TransferQuota{})
	tr = newTransfer(baseTransfer, nil, r, 10)
	pos, err := tr.Seek(10, 0)
	assert.NoError(t, err)
	assert.Equal(t, pos, tr.expectedOffset)
	err = tr.closeIO()
	assert.NoError(t, err)

	r, w, err := pipeat.Pipe()
	assert.NoError(t, err)
	pipeWriter := vfs.NewPipeWriter(w)
	baseTransfer = common.NewBaseTransfer(nil, connection.BaseConnection, nil, testfile, testfile, testfile,
		common.TransferUpload, 0, 0, 0, 0, false, fs, dataprovider.TransferQuota{})
	tr = newTransfer(baseTransfer, pipeWriter, nil, 0)

	err = r.Close()
	assert.NoError(t, err)
	errFake := fmt.Errorf("fake upload error")
	go func() {
		time.Sleep(100 * time.Millisecond)
		pipeWriter.Done(errFake)
	}()
	err = tr.closeIO()
	assert.EqualError(t, err, errFake.Error())
	_, err = tr.Seek(1, 0)
	if assert.Error(t, err) {
		assert.EqualError(t, err, common.ErrOpUnsupported.Error())
	}
	err = os.Remove(testfile)
	assert.NoError(t, err)
}

func TestVerifyTLSConnection(t *testing.T) {
	oldCertMgr := certMgr

	caCrlPath := filepath.Join(os.TempDir(), "testcrl.crt")
	certPath := filepath.Join(os.TempDir(), "test.crt")
	keyPath := filepath.Join(os.TempDir(), "test.key")
	err := os.WriteFile(caCrlPath, []byte(caCRL), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(certPath, []byte(ftpsCert), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(keyPath, []byte(ftpsKey), os.ModePerm)
	assert.NoError(t, err)
	keyPairs := []common.TLSKeyPair{
		{
			Cert: certPath,
			Key:  keyPath,
			ID:   common.DefaultTLSKeyPaidID,
		},
	}
	certMgr, err = common.NewCertManager(keyPairs, "", "ftp_test")
	assert.NoError(t, err)

	certMgr.SetCARevocationLists([]string{caCrlPath})
	err = certMgr.LoadCRLs()
	assert.NoError(t, err)

	crt, err := tls.X509KeyPair([]byte(client1Crt), []byte(client1Key))
	assert.NoError(t, err)
	x509crt, err := x509.ParseCertificate(crt.Certificate[0])
	assert.NoError(t, err)

	server := Server{}
	state := tls.ConnectionState{
		PeerCertificates: []*x509.Certificate{x509crt},
	}

	err = server.verifyTLSConnection(state)
	assert.Error(t, err) // no verified certification chain

	crt, err = tls.X509KeyPair([]byte(caCRT), []byte(caKey))
	assert.NoError(t, err)

	x509CAcrt, err := x509.ParseCertificate(crt.Certificate[0])
	assert.NoError(t, err)

	state.VerifiedChains = append(state.VerifiedChains, []*x509.Certificate{x509crt, x509CAcrt})
	err = server.verifyTLSConnection(state)
	assert.NoError(t, err)

	crt, err = tls.X509KeyPair([]byte(client2Crt), []byte(client2Key))
	assert.NoError(t, err)
	x509crtRevoked, err := x509.ParseCertificate(crt.Certificate[0])
	assert.NoError(t, err)

	state.VerifiedChains = append(state.VerifiedChains, []*x509.Certificate{x509crtRevoked, x509CAcrt})
	state.PeerCertificates = []*x509.Certificate{x509crtRevoked}
	err = server.verifyTLSConnection(state)
	assert.EqualError(t, err, common.ErrCrtRevoked.Error())

	err = os.Remove(caCrlPath)
	assert.NoError(t, err)
	err = os.Remove(certPath)
	assert.NoError(t, err)
	err = os.Remove(keyPath)
	assert.NoError(t, err)

	certMgr = oldCertMgr
}

func TestCiphers(t *testing.T) {
	b := Binding{
		TLSCipherSuites: []string{},
	}
	b.setCiphers()
	require.Nil(t, b.ciphers)
	b.TLSCipherSuites = []string{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384"}
	b.setCiphers()
	require.Len(t, b.ciphers, 2)
	require.Equal(t, []uint16{tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384}, b.ciphers)
}

func TestPassiveIPResolver(t *testing.T) {
	b := Binding{
		PassiveIPOverrides: []PassiveIPOverride{
			{},
		},
	}
	err := b.checkPassiveIP()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "passive IP networks override cannot be empty")
	b = Binding{
		PassiveIPOverrides: []PassiveIPOverride{
			{
				IP: "invalid ip",
			},
		},
	}
	err = b.checkPassiveIP()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not valid")

	b = Binding{
		PassiveIPOverrides: []PassiveIPOverride{
			{
				IP:       "192.168.1.1",
				Networks: []string{"192.168.1.0/24", "invalid cidr"},
			},
		},
	}
	err = b.checkPassiveIP()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid passive IP networks override")
	b = Binding{
		ForcePassiveIP: "192.168.2.1",
		PassiveIPOverrides: []PassiveIPOverride{
			{
				IP:       "::ffff:192.168.1.1",
				Networks: []string{"192.168.1.0/24"},
			},
		},
	}
	err = b.checkPassiveIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, b.PassiveIPOverrides[0].GetNetworksAsString())
	assert.Equal(t, "192.168.1.1", b.PassiveIPOverrides[0].IP)
	require.Len(t, b.PassiveIPOverrides[0].parsedNetworks, 1)
	ip := net.ParseIP("192.168.1.2")
	assert.True(t, b.PassiveIPOverrides[0].parsedNetworks[0](ip))
	ip = net.ParseIP("192.168.0.2")
	assert.False(t, b.PassiveIPOverrides[0].parsedNetworks[0](ip))

	mockCC := mockFTPClientContext{
		remoteIP: "192.168.1.10",
		localIP:  "192.168.1.3",
	}
	passiveIP, err := b.passiveIPResolver(mockCC)
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", passiveIP)
	b.PassiveIPOverrides[0].IP = ""
	passiveIP, err = b.passiveIPResolver(mockCC)
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.3", passiveIP)
	mockCC.remoteIP = "172.16.2.3"
	passiveIP, err = b.passiveIPResolver(mockCC)
	assert.NoError(t, err)
	assert.Equal(t, b.ForcePassiveIP, passiveIP)
}

func TestRelativePath(t *testing.T) {
	rel := getPathRelativeTo("/testpath", "/testpath")
	assert.Empty(t, rel)
	rel = getPathRelativeTo("/", "/")
	assert.Empty(t, rel)
	rel = getPathRelativeTo("/", "/dir/sub")
	assert.Equal(t, "dir/sub", rel)
	rel = getPathRelativeTo("./", "/dir/sub")
	assert.Equal(t, "/dir/sub", rel)
	rel = getPathRelativeTo("/sub", "/dir/sub")
	assert.Equal(t, "../dir/sub", rel)
	rel = getPathRelativeTo("/dir", "/dir/sub")
	assert.Equal(t, "sub", rel)
	rel = getPathRelativeTo("/dir/sub", "/dir")
	assert.Equal(t, "../", rel)
	rel = getPathRelativeTo("dir", "/dir1")
	assert.Equal(t, "/dir1", rel)
	rel = getPathRelativeTo("", "/dir2")
	assert.Equal(t, "dir2", rel)
	rel = getPathRelativeTo(".", "/dir2")
	assert.Equal(t, "/dir2", rel)
	rel = getPathRelativeTo("/dir3", "dir3")
	assert.Equal(t, "dir3", rel)
}
