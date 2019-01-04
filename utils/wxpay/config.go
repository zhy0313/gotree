package wxpay

const (
	// 即时用车微信支付账号
	CALL_WECHAT_ACCOUNT = `{"AppId":"wxe38d5ae955d4362e","MchId":"1232303202","Key":"qwertyuiiuytrewqqwertyui12345678"}`
	// 小程序微信支付账号
	SP_WECHAT_ACCOUNT = `{"AppId":"wxe38d5ae955d4362e","MchId":"1232303202","Key":"qwertyuiiuytrewqqwertyui12345678"}`
)

// 证书相关 注意 ***所有证书必须顶头写***
// 即时用车微信账号证书
var CallCertPem = []byte(`
-----BEGIN CERTIFICATE-----
MIIEcTCCA9qgAwIBAgIDExbHMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE2MDMwMzEwNDY0OVoX
DTI2MDMwMTEwNDY0OVowgaExCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv
bmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL
EwVNTVBheTE2MDQGA1UEAxQt5YyX5Lqs5YGH5pel6Ziz5YWJ546v55CD5peF6KGM
56S+5pyJ6ZmQ5YWs5Y+4MREwDwYDVQQEEwgxMDEyNjI5NTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALUA9gaylzdu1+A21JLpMKXd2PmRpB3mNp+ah6Hc
+CBuVm9W/wq1Ct40VzDAe3+QphC9SniHSbc0zw2RAc+s8iF4nYo1PNw5lh8vSNgf
sJpv1Al+1EK/VtpoGlfMl+h036K3hcDDcns2sN+xLLYIs4hXl1t4EYfIPPd5iOSs
6MaPwSyppQOTcMu7QKFVI0ottFYOp0MbiPwDFhKgwOoYNWGU77h+AstHYdkdS9A9
LXRI/sl2hLbzKvfohRR0U7eZOB0HSq1k9d+qJbHVdCFKzsGUA48Dxf+yeNiW4c0A
069StyCz/xoq/ZTDJlXYp43kljWRLRgiUTMB2Vuio2tZ77kCAwEAAaOCAUYwggFC
MAkGA1UdEwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2Vy
dGlmaWNhdGUiMB0GA1UdDgQWBBRWpVfDvPwSxpe5pumzCt4VszrI9jCBvwYDVR0j
BIG3MIG0gBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMC
Q04xEjAQBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNV
BAoTB1RlbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEf
MB0GCSqGSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1Ud
DwEB/wQEAwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUF
AAOBgQAjQEgX5ILIC6ylBjGO1EGeAUp/9v+9DOfZUyfeJnq7MsMxcRISMgDudbXl
XaQb4ZHyRQvauh4hy/0/O1x8rEHaVhQytDEPbtbR+ouUWQ1qRy9pWsY/3hZA7cxt
PNRyoiSWqxWF/B7JUR02ms2470xQXJVlj1+Lfatddv7vxn6aqg==
-----END CERTIFICATE-----
	`)

// 即时用车微信账号证书Key
var CallKeyPem = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC1APYGspc3btfg
NtSS6TCl3dj5kaQd5jafmoeh3PggblZvVv8KtQreNFcwwHt/kKYQvUp4h0m3NM8N
kQHPrPIheJ2KNTzcOZYfL0jYH7Cab9QJftRCv1baaBpXzJfodN+it4XAw3J7NrDf
sSy2CLOIV5dbeBGHyDz3eYjkrOjGj8EsqaUDk3DLu0ChVSNKLbRWDqdDG4j8AxYS
oMDqGDVhlO+4fgLLR2HZHUvQPS10SP7JdoS28yr36IUUdFO3mTgdB0qtZPXfqiWx
1XQhSs7BlAOPA8X/snjYluHNANOvUrcgs/8aKv2UwyZV2KeN5JY1kS0YIlEzAdlb
oqNrWe+5AgMBAAECggEAMppMLd8r63FvpN1vXIsY8KYvDMaszTcZOlGnbZcrP1XZ
kFMQBlxK54hWf+iwHh/AbZmeAkhAUZnP8QkDKp49KyJfWM49b4wh7iH4CYQCiYqO
CwGMMgTwMSs9giJOrcxO4VzRZey+SrglhsQvrcyM9KcYR5gvvng6oy01GklE4o3B
jUXLSBQHBDBqzbDy4JPf0kK3mo2GPtHJk6uYZEUgEyVZTyDa595pjy9Dw8d4/99p
epz7aJy82ew1KID9/Tmm2pl96pvS6nfEdeKYcxiq/Aov54IVWcPfj5EsgHYf0VdE
SS63EgdNXkX+IvAhEKn1JDk1Ql9j1SASUCMidXLqIQKBgQDf7yODi0DD/+niNveA
HHrlFLxySUX6VqmkubX6ZXslzBINeSu9bS06bwY/ZAypXKeu1C4/sm19S9HqfRN0
e+7HDL1HdvbH6+kk5vZ8T3N+fxlAOnYsUBI1OKvfocqHC208DLSLzQGhRpI47QP9
gZC16K/m5Z6g9txAJ2GRNrKkxQKBgQDO7BpI9Rrh40VvxH2Kyh8d+RrZxm2tGcQS
2LDt6/WOcs/rqHJLgVmKSbDXkJymh3VdV9093plxVKm3gQVKdAjy8kWXENk94GD5
ADPNkmQuTBN7EefSP3erjureJ0Nfnz4R+O78zl831/6fGIrVr7uqOZ8NiBl35DwS
3yu2kT4WZQKBgGbCnG9u5YeL1k4Cn0zgxNx+yYNAcKZSQoLe3c1L6FkN7nLUWegR
Q6H+9MT+KnlFlYU6xQZh4LCQrIGIZ/caMBaTmABFbTWM4m4WtqGQ7BTuSi4ZJcVr
8Q8PNH/pBME30yatReSpbMgPVGZfDWe1nyx63M1+LW78GVIvQCydBxlpAoGBAKVS
O/n4Yp8BXvPacFdX/56J7Sr7f5sij+Zi3JFqyYkjL/3fWln7IZf8Il9IOfBPH7UR
Q0FwPPYwJ1zmp1yB8rhwWqtEmdz3DWNEBx+Ci6n1vEbC2o2/iZQ3Hm2ZvxmB+CyR
0BeJpsfOOa/RAvORcQWi/fHowDhq0JhfV+SIjKuFAoGAdPcvh3KDmpD9YqE8v8FJ
fDvw67nA3fb1VU9ZJ2naxCyy6acfdEiBw7qHUs86Fy+Zbb60rnT8Cp++zgC7mgRD
nKJda/lQuJ/PkBKrotXbaiSGtt98wWOT2L28LOR9HQpjdaThoKcJwGdNqf0oEHoD
x0DNGxeORlFewgyegOa92lI=
-----END PRIVATE KEY-----
	`)

// 即时用车微信账号root证书
var CallRootPem = []byte(`
-----BEGIN CERTIFICATE-----
MIIDIDCCAomgAwIBAgIENd70zzANBgkqhkiG9w0BAQUFADBOMQswCQYDVQQGEwJV
UzEQMA4GA1UEChMHRXF1aWZheDEtMCsGA1UECxMkRXF1aWZheCBTZWN1cmUgQ2Vy
dGlmaWNhdGUgQXV0aG9yaXR5MB4XDTk4MDgyMjE2NDE1MVoXDTE4MDgyMjE2NDE1
MVowTjELMAkGA1UEBhMCVVMxEDAOBgNVBAoTB0VxdWlmYXgxLTArBgNVBAsTJEVx
dWlmYXggU2VjdXJlIENlcnRpZmljYXRlIEF1dGhvcml0eTCBnzANBgkqhkiG9w0B
AQEFAAOBjQAwgYkCgYEAwV2xWGcIYu6gmi0fCG2RFGiYCh7+2gRvE4RiIcPRfM6f
BeC4AfBONOziipUEZKzxa1NfBbPLZ4C/QgKO/t0BCezhABRP/PvwDN1Dulsr4R+A
cJkVV5MW8Q+XarfCaCMczE1ZMKxRHjuvK9buY0V7xdlfUNLjUA86iOe/FP3gx7kC
AwEAAaOCAQkwggEFMHAGA1UdHwRpMGcwZaBjoGGkXzBdMQswCQYDVQQGEwJVUzEQ
MA4GA1UEChMHRXF1aWZheDEtMCsGA1UECxMkRXF1aWZheCBTZWN1cmUgQ2VydGlm
aWNhdGUgQXV0aG9yaXR5MQ0wCwYDVQQDEwRDUkwxMBoGA1UdEAQTMBGBDzIwMTgw
ODIyMTY0MTUxWjALBgNVHQ8EBAMCAQYwHwYDVR0jBBgwFoAUSOZo+SvSspXXR9gj
IBBPM5iQn9QwHQYDVR0OBBYEFEjmaPkr0rKV10fYIyAQTzOYkJ/UMAwGA1UdEwQF
MAMBAf8wGgYJKoZIhvZ9B0EABA0wCxsFVjMuMGMDAgbAMA0GCSqGSIb3DQEBBQUA
A4GBAFjOKer89961zgK5F7WF0bnj4JXMJTENAKaSbn+2kmOeUJXRmm/kEd5jhW6Y
7qj/WsjTVbJmcVfewCHrPSqnI0kBBIZCe/zuf6IWUrVnZ9NA2zsmWLIodz2uFHdh
1voqZiegDfqnc1zqcPGUIWVEX/r87yloqaKHee9570+sB3c4
-----END CERTIFICATE-----%
	`)

// 小程序微信账号证书
var SpCertPem = []byte(`
-----BEGIN CERTIFICATE-----
MIIEcTCCA9qgAwIBAgIDExbHMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE2MDMwMzEwNDY0OVoX
DTI2MDMwMTEwNDY0OVowgaExCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv
bmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL
EwVNTVBheTE2MDQGA1UEAxQt5YyX5Lqs5YGH5pel6Ziz5YWJ546v55CD5peF6KGM
56S+5pyJ6ZmQ5YWs5Y+4MREwDwYDVQQEEwgxMDEyNjI5NTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALUA9gaylzdu1+A21JLpMKXd2PmRpB3mNp+ah6Hc
+CBuVm9W/wq1Ct40VzDAe3+QphC9SniHSbc0zw2RAc+s8iF4nYo1PNw5lh8vSNgf
sJpv1Al+1EK/VtpoGlfMl+h036K3hcDDcns2sN+xLLYIs4hXl1t4EYfIPPd5iOSs
6MaPwSyppQOTcMu7QKFVI0ottFYOp0MbiPwDFhKgwOoYNWGU77h+AstHYdkdS9A9
LXRI/sl2hLbzKvfohRR0U7eZOB0HSq1k9d+qJbHVdCFKzsGUA48Dxf+yeNiW4c0A
069StyCz/xoq/ZTDJlXYp43kljWRLRgiUTMB2Vuio2tZ77kCAwEAAaOCAUYwggFC
MAkGA1UdEwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2Vy
dGlmaWNhdGUiMB0GA1UdDgQWBBRWpVfDvPwSxpe5pumzCt4VszrI9jCBvwYDVR0j
BIG3MIG0gBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMC
Q04xEjAQBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNV
BAoTB1RlbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEf
MB0GCSqGSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1Ud
DwEB/wQEAwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUF
AAOBgQAjQEgX5ILIC6ylBjGO1EGeAUp/9v+9DOfZUyfeJnq7MsMxcRISMgDudbXl
XaQb4ZHyRQvauh4hy/0/O1x8rEHaVhQytDEPbtbR+ouUWQ1qRy9pWsY/3hZA7cxt
PNRyoiSWqxWF/B7JUR02ms2470xQXJVlj1+Lfatddv7vxn6aqg==
-----END CERTIFICATE-----
	`)

// 小程序微信账号证书Key
var SpKeyPem = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC1APYGspc3btfg
NtSS6TCl3dj5kaQd5jafmoeh3PggblZvVv8KtQreNFcwwHt/kKYQvUp4h0m3NM8N
kQHPrPIheJ2KNTzcOZYfL0jYH7Cab9QJftRCv1baaBpXzJfodN+it4XAw3J7NrDf
sSy2CLOIV5dbeBGHyDz3eYjkrOjGj8EsqaUDk3DLu0ChVSNKLbRWDqdDG4j8AxYS
oMDqGDVhlO+4fgLLR2HZHUvQPS10SP7JdoS28yr36IUUdFO3mTgdB0qtZPXfqiWx
1XQhSs7BlAOPA8X/snjYluHNANOvUrcgs/8aKv2UwyZV2KeN5JY1kS0YIlEzAdlb
oqNrWe+5AgMBAAECggEAMppMLd8r63FvpN1vXIsY8KYvDMaszTcZOlGnbZcrP1XZ
kFMQBlxK54hWf+iwHh/AbZmeAkhAUZnP8QkDKp49KyJfWM49b4wh7iH4CYQCiYqO
CwGMMgTwMSs9giJOrcxO4VzRZey+SrglhsQvrcyM9KcYR5gvvng6oy01GklE4o3B
jUXLSBQHBDBqzbDy4JPf0kK3mo2GPtHJk6uYZEUgEyVZTyDa595pjy9Dw8d4/99p
epz7aJy82ew1KID9/Tmm2pl96pvS6nfEdeKYcxiq/Aov54IVWcPfj5EsgHYf0VdE
SS63EgdNXkX+IvAhEKn1JDk1Ql9j1SASUCMidXLqIQKBgQDf7yODi0DD/+niNveA
HHrlFLxySUX6VqmkubX6ZXslzBINeSu9bS06bwY/ZAypXKeu1C4/sm19S9HqfRN0
e+7HDL1HdvbH6+kk5vZ8T3N+fxlAOnYsUBI1OKvfocqHC208DLSLzQGhRpI47QP9
gZC16K/m5Z6g9txAJ2GRNrKkxQKBgQDO7BpI9Rrh40VvxH2Kyh8d+RrZxm2tGcQS
2LDt6/WOcs/rqHJLgVmKSbDXkJymh3VdV9093plxVKm3gQVKdAjy8kWXENk94GD5
ADPNkmQuTBN7EefSP3erjureJ0Nfnz4R+O78zl831/6fGIrVr7uqOZ8NiBl35DwS
3yu2kT4WZQKBgGbCnG9u5YeL1k4Cn0zgxNx+yYNAcKZSQoLe3c1L6FkN7nLUWegR
Q6H+9MT+KnlFlYU6xQZh4LCQrIGIZ/caMBaTmABFbTWM4m4WtqGQ7BTuSi4ZJcVr
8Q8PNH/pBME30yatReSpbMgPVGZfDWe1nyx63M1+LW78GVIvQCydBxlpAoGBAKVS
O/n4Yp8BXvPacFdX/56J7Sr7f5sij+Zi3JFqyYkjL/3fWln7IZf8Il9IOfBPH7UR
Q0FwPPYwJ1zmp1yB8rhwWqtEmdz3DWNEBx+Ci6n1vEbC2o2/iZQ3Hm2ZvxmB+CyR
0BeJpsfOOa/RAvORcQWi/fHowDhq0JhfV+SIjKuFAoGAdPcvh3KDmpD9YqE8v8FJ
fDvw67nA3fb1VU9ZJ2naxCyy6acfdEiBw7qHUs86Fy+Zbb60rnT8Cp++zgC7mgRD
nKJda/lQuJ/PkBKrotXbaiSGtt98wWOT2L28LOR9HQpjdaThoKcJwGdNqf0oEHoD
x0DNGxeORlFewgyegOa92lI=
-----END PRIVATE KEY-----
	`)

// 小程序微信账号root证书
var SpRootPem = []byte(`
-----BEGIN CERTIFICATE-----
MIIDIDCCAomgAwIBAgIENd70zzANBgkqhkiG9w0BAQUFADBOMQswCQYDVQQGEwJV
UzEQMA4GA1UEChMHRXF1aWZheDEtMCsGA1UECxMkRXF1aWZheCBTZWN1cmUgQ2Vy
dGlmaWNhdGUgQXV0aG9yaXR5MB4XDTk4MDgyMjE2NDE1MVoXDTE4MDgyMjE2NDE1
MVowTjELMAkGA1UEBhMCVVMxEDAOBgNVBAoTB0VxdWlmYXgxLTArBgNVBAsTJEVx
dWlmYXggU2VjdXJlIENlcnRpZmljYXRlIEF1dGhvcml0eTCBnzANBgkqhkiG9w0B
AQEFAAOBjQAwgYkCgYEAwV2xWGcIYu6gmi0fCG2RFGiYCh7+2gRvE4RiIcPRfM6f
BeC4AfBONOziipUEZKzxa1NfBbPLZ4C/QgKO/t0BCezhABRP/PvwDN1Dulsr4R+A
cJkVV5MW8Q+XarfCaCMczE1ZMKxRHjuvK9buY0V7xdlfUNLjUA86iOe/FP3gx7kC
AwEAAaOCAQkwggEFMHAGA1UdHwRpMGcwZaBjoGGkXzBdMQswCQYDVQQGEwJVUzEQ
MA4GA1UEChMHRXF1aWZheDEtMCsGA1UECxMkRXF1aWZheCBTZWN1cmUgQ2VydGlm
aWNhdGUgQXV0aG9yaXR5MQ0wCwYDVQQDEwRDUkwxMBoGA1UdEAQTMBGBDzIwMTgw
ODIyMTY0MTUxWjALBgNVHQ8EBAMCAQYwHwYDVR0jBBgwFoAUSOZo+SvSspXXR9gj
IBBPM5iQn9QwHQYDVR0OBBYEFEjmaPkr0rKV10fYIyAQTzOYkJ/UMAwGA1UdEwQF
MAMBAf8wGgYJKoZIhvZ9B0EABA0wCxsFVjMuMGMDAgbAMA0GCSqGSIb3DQEBBQUA
A4GBAFjOKer89961zgK5F7WF0bnj4JXMJTENAKaSbn+2kmOeUJXRmm/kEd5jhW6Y
7qj/WsjTVbJmcVfewCHrPSqnI0kBBIZCe/zuf6IWUrVnZ9NA2zsmWLIodz2uFHdh
1voqZiegDfqnc1zqcPGUIWVEX/r87yloqaKHee9570+sB3c4
-----END CERTIFICATE-----%
	`)
