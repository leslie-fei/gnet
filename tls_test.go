package gnet

import (
	"bufio"
	tls2 "crypto/tls"
	"io"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/panjf2000/gnet/v2/pkg/logging"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/panjf2000/gnet/v2/pkg/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var serverCRT = `-----BEGIN CERTIFICATE-----
MIIDczCCAlugAwIBAgIJAM0/vF9KXTQbMA0GCSqGSIb3DQEBCwUAMFAxCzAJBgNV
BAYTAkNOMQwwCgYDVQQIDANGRUkxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoG
A1UECgwTRGVmYXVsdCBDb21wYW55IEx0ZDAeFw0yNDAzMjIwNzIxMTNaFw0yNTAz
MjIwNzIxMTNaMFAxCzAJBgNVBAYTAkNOMQwwCgYDVQQIDANGRUkxFTATBgNVBAcM
DERlZmF1bHQgQ2l0eTEcMBoGA1UECgwTRGVmYXVsdCBDb21wYW55IEx0ZDCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPXZ3mcbFEMmmDaQWyKOMOfPmxQG
lWGdQg9STQard+K/RY36WTYshTHNAQzZX1TqdCdReCSDgGkE5pcb7Pvw1UBa6c2b
BfpbmPwAG7oEQvyqIRgxW65w9ZUHl0A2+Ynvfz6mRJdZ1VRXbUKBZAjdHTwGbOZb
Tdw6niih42i8TUPrgCEq1F/NomiVsZqJTfXXwxsRCo7oZ1Gepb/ljOJ8Q3Iu78Ng
AQHCj+a2oxuRPMydHAhBbdPgQWSvO/C4864cXt6gtetsj4wR99syFX0k6vYeMdmZ
G9OE7jCmvU+F/hIDa60PXYXG4C+rKe01paUMIhrMyhV6GoGvrr97sX9Bp4UCAwEA
AaNQME4wHQYDVR0OBBYEFIYICdXG8UZbNja0OJP4E0UdhzQXMB8GA1UdIwQYMBaA
FIYICdXG8UZbNja0OJP4E0UdhzQXMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEL
BQADggEBABBsRLAtQP+92tXcTNkiZ86R6PSpRAJQB88OCvqVF6JZPFvOQjDHhjvW
mwhEEgbUZACTetQpcCk/SbOebeZUnmm1dx0XG7qaLVVSedUy9gjHnxkhS5Ws8cOW
vZVot9ucr3p7cGBYDGtUZ7oN/f2HpwoiQg2DqEHgx/1cNKza3x1+WANIZNoKNjP5
ZDhDuHqt2/AHxyrlkVbpVS+g2Mn9K68rwJV6cHe3usXrlc4VT2fDTkFbPgVbyhBx
l8lG16b/ko5jSa4NN8MKZKZZYZnZr/452qc9E3/nk6MTimkCB6zWfDx2RsDR7+nE
jqzuMt+wx35XZSA5Jfms71+cU8KvAEw=
-----END CERTIFICATE-----
`

var serverKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA9dneZxsUQyaYNpBbIo4w58+bFAaVYZ1CD1JNBqt34r9FjfpZ
NiyFMc0BDNlfVOp0J1F4JIOAaQTmlxvs+/DVQFrpzZsF+luY/AAbugRC/KohGDFb
rnD1lQeXQDb5ie9/PqZEl1nVVFdtQoFkCN0dPAZs5ltN3DqeKKHjaLxNQ+uAISrU
X82iaJWxmolN9dfDGxEKjuhnUZ6lv+WM4nxDci7vw2ABAcKP5rajG5E8zJ0cCEFt
0+BBZK878Ljzrhxe3qC162yPjBH32zIVfSTq9h4x2Zkb04TuMKa9T4X+EgNrrQ9d
hcbgL6sp7TWlpQwiGszKFXoaga+uv3uxf0GnhQIDAQABAoIBAGEjEfKfCmR1zVJR
uxBuKfiNWHKTH1aQMden5vMpsJaQiETOk40KDh5dJ3sSL853TCl00QlsQ8I6w3kU
+Y+Z2kgxI4uvIq+GnW2EBXB72pPGDIlohxbnXsDhXOtsTTAoxpfUTemW/ujyGiNa
yFGQO8rJpyFlFeBb8L2RKS7qp4Sbs6cC0qAzMIW+OBoSKu07bTphWKxyiextVB7d
z9JEEnpLb6VMRA2SEQAUNJ0uekM1oXfYh5hFjUtd+UKCF+KFTqXh1N7exFEeY5dO
thateML4BaFH1bvhE3TzR3i/el2NGWJ/BUCkm5kjIoLqw/We2TCnz5RWmtjArpQS
7iAIdgECgYEA+8WB7jktLEQm5Sr1EVVdNfkYl4UYaKQyOOWSQi4XN4cL+4Te62HD
mZqrsD/x84Nv36gzcUzfyQvPzz72Cav3Z8SZHP/YhHSRD+U7bxtyo6z+Kv0WDtQD
ZeLdjZW0X/Vl9kwZU3udM0MHOsPRjcxQVhWNwEknfUwwgvkoJ0buxEUCgYEA+frn
+m2DuWubpSTp8UxIOYPSLag9q+sQKMriJPUXvzIpZpBdcLhvINPAjluzTU2HHekK
nF0ciDxJ0QR1SB8UhbtVGYlASxP/K0BEYhydOdDg053qiD7ZtNobZedEcp72pRhF
lpG1IaqNUUY/MC7qbuoqzum0FTeY7ATzt+fcqkECgYAFu1Vd6wdQPqz0AqpZnqvk
kqswp+MKXtKV4QQXZZ05smVfcILysOMXGCJYfNwxQwjE+rtRFQLzPMWaXkCTr8qO
HjeuJVT0jwLObdWQjftqJoIS0Hg0FjeNlkCc/I+Cr7+Gf+0DrYjCpgAAIn0e0Uvy
6cK0YfX4MEfac3mqiz7cUQKBgFS2a2CZNjdiuQcY4FKOODeNbW9dgTC+YuaAAmzl
xdtcI12HdLLbPYWRXMKhXrzGlhBTxe3BPAcK7zTu3tPG/IxpNGiH84q9xw4gJExf
SASsvVPZQv5S1jIMa5Z37eEt74V5GTPyQsCa7HLrbERm0bkuj704w3chNShecT0U
TNpBAoGBAI4b2Unjae+qIFncPGdZhJBIj4UlEC5ePCtUo16Bys/IvFcIveonP3iq
lHi05PF9PNyPGFY0/h3sHfarq5K4cBZ+dsHBV7/KzJltZ3K0mZodaI3E/hB0SZgn
oHQA9bqSDRzT61MFUomrk/56XIXYFjSPBWQI8AwTy8+8Xm+NczxN
-----END RSA PRIVATE KEY-----
`

var clientCRT = `-----BEGIN CERTIFICATE-----
MIIDHDCCAgQCCQChed2z264DpDANBgkqhkiG9w0BAQsFADBQMQswCQYDVQQGEwJD
TjEMMAoGA1UECAwDRkVJMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAaBgNVBAoM
E0RlZmF1bHQgQ29tcGFueSBMdGQwHhcNMjQwMzIyMDcyMTU4WhcNMjUwMzIyMDcy
MTU4WjBQMQswCQYDVQQGEwJDTjEMMAoGA1UECAwDRkVJMRUwEwYDVQQHDAxEZWZh
dWx0IENpdHkxHDAaBgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQClSoCJXax5xhGM+c1ux7nElmwbjKj3fvS8
opMj9nNp1Q5XuauNjgZSvkMAYBQYsrBn5RG+7EpwzPwQMSowzv8RkviDYz9i0/SO
YZLyX/Y4f7DTL21ryKXoxirxB+LHvZSvs7iunhqB3AIXeIlu6WsvGMXhhHOP0diS
n9Uq1598m8OWB4x9k573Lr3UoyOAOQPPra3wZxsVRn+7W6K5xnmPgzYm0p3ac5jF
+E9ZhUlVbuMxKF6xl+X3pGyu53BqcXdWDWdRI+TJyMRzB0kXdOQx536wW6i/1aUv
crgIo2zEoVBH3vylOXfp0K7l3+jc9fQRqEb6xiZru95jZDpPCuYTAgMBAAEwDQYJ
KoZIhvcNAQELBQADggEBAN58zr83coygDbo61HDSdARQ3XWyeTazsmbuy6WIwXrl
zPwoFcksjCOwraFccnowwQfteHMdL5EB0Lk2ihwtIxUwzJVOQTovb2B4TqcpV0nd
/xC5aWJr3JiX+FKiDk450FXQXx5xMR1zA6YsoDmXd7wmK0Z6pv49V+7WjdIDeZlm
FnZ9l3GOT8WdbLmzZ1K1Qsa1L723ENTJhIQJR/UsiJJl0ATykd6y/wX3aUUUtBmq
9P6jcQoRp/x+9eM5iIa+YTKBaUEwTRUEZmklb1s82FC8jhly7B9bmMYEDs8QfNdX
rngcyY2rD5ATuODBRks1LO2G7JDx4vJt1Wh04RInEg0=
-----END CERTIFICATE-----
`
var clientKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApUqAiV2secYRjPnNbse5xJZsG4yo9370vKKTI/ZzadUOV7mr
jY4GUr5DAGAUGLKwZ+URvuxKcMz8EDEqMM7/EZL4g2M/YtP0jmGS8l/2OH+w0y9t
a8il6MYq8Qfix72Ur7O4rp4agdwCF3iJbulrLxjF4YRzj9HYkp/VKteffJvDlgeM
fZOe9y691KMjgDkDz62t8GcbFUZ/u1uiucZ5j4M2JtKd2nOYxfhPWYVJVW7jMShe
sZfl96RsrudwanF3Vg1nUSPkycjEcwdJF3TkMed+sFuov9WlL3K4CKNsxKFQR978
pTl36dCu5d/o3PX0EahG+sYma7veY2Q6TwrmEwIDAQABAoIBABjl8BbBFYoCXAJw
XpGhZEpInmBH/+vII/WGR8PTC++VOGPehAj+q9Uq88zKrDxMLmZ4CwRHc68qMY3+
2ewgQivTh7vC53lrlBsed5QrAFviYDBvpHvn7ZXyw+Uq2u5bkFonKFqlVjjZezct
4QXkA1FE2gnuebngUk9D5Ehf0pZiYki3t5+UbZ/2XIlU6qcwYXwTD8N3NDmhKDcn
2q/JatBp/l7gE+RJPYyBPgM9ffw273GraQhMwmeBkP0wfmxN/rntrO1ViuHaxwk9
qS6lBNghBcM8uIpihrm7hf080VqDUte4OG8Z6gP3kLeLg7vBubFUN/weGqRGKcDN
Ifoa9oECgYEA0Sfu/54O2fNSI3UiHrzExmeLpCDTYJAWMNViporSML6V0DnM0Hso
D+flogxns8t4QW1pJ5AlANmm13piavl44eKO5UyE5bIzCzbOi3v75aTgtIOd4M77
7HId0TR62xAeK23j7rBNoIpwNu9vYIYhh8kD3RJ+8DC1akhy2+398X8CgYEAyk+M
TTPZEajkk74Q46vu1ByAVL2OIqQq3dO292vW9z0YKl0WvehKPgjXEmxZtuUdvwAK
3dQgDNaF3kMAx0mluwCoPHFDgMR0tw51q3Jxb3sIgQo9RKUMuhTRzAwZQQW2qVYq
w+moY6rZgEiwjSEqKxBEviAiHAeOLjg5E4HobW0CgYEAlWMZYDerUd5px2uNIDkh
rrgcU/EiKLEhOGXgnSU/9ZLafCwcbrfYKGwwZqrHkvJ8BPc5iQl6dH4AnZ0LGG7w
+R2RhJM/xLOfo53UrmmyxBYcbLIpHH2aEAp944PBUB/V/JzzlPMVrtnSRlqDyYRk
BJD9knP4eGPB35S5fP4iR88CgYAcFp1JHri1jfhmVs4rffSDuI1MtqCkSJy2/az6
fxtyaiYPw6rjpPuhz0JlCWwPS5IssYnQsKHehm5BpJXYFW6QHxA8jBEuMJhvze8i
ZZpAOxm8KLrG7A9vSnKyXL+o2m0454C6gJdjdxYJjI4OmhsrrW+N4hI6SAZCfoXW
IoQtpQKBgF/JXD9cy4CJjJC2RQ4SSEz+yljYdQECeBtQjiodDq7x/l2ELr5gIO1A
VHDSuE+GXiRlWz1OmRpVh6DiXR6wIYdHJ8xw4rBPGtxEeBUoMWL+xpa6OJ1OX/2s
nXNUa+27KshlpdWNlwJOtkeeP7yMQGf7irAelQ/P0a3qmoop8YvC
-----END RSA PRIVATE KEY-----
`

func getServerConfig() *tls.Config {
	cert, err := tls.X509KeyPair([]byte(serverCRT), []byte(serverKey))
	if nil != err {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func getClientTLSConfig() *tls.Config {
	crt, err := tls.X509KeyPair([]byte(clientCRT), []byte(clientKey))
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{crt}, InsecureSkipVerify: true}
}

func getGoClientTLSConfig() *tls2.Config {
	crt, err := tls2.X509KeyPair([]byte(clientCRT), []byte(clientKey))
	if err != nil {
		panic(err)
	}
	return &tls2.Config{Certificates: []tls2.Certificate{crt}, InsecureSkipVerify: true}
}

type testTLSServer struct {
	*testServer
}

func (s *testTLSServer) OnClose(c Conn, err error) (action Action) {
	if err != nil {
		logging.Debugf("error occurred on closed, %v\n", err)
	}
	/*if c.LocalAddr().Network() != "udp" {
		require.IsType(s.tester, (*tlsConn)(nil), c.Context(), "want *tlsConn")
	}*/

	atomic.AddInt32(&s.disconnected, 1)
	return
}

func (s *testTLSServer) OnTick() (delay time.Duration, action Action) {
	delay = 100 * time.Millisecond
	if atomic.CompareAndSwapInt32(&s.started, 0, 1) {
		for _, protoAddr := range s.addrs {
			proto, addr, err := parseProtoAddr(protoAddr)
			assert.NoError(s.tester, err)
			for i := 0; i < s.nclients; i++ {
				atomic.AddInt32(&s.clientActive, 1)
				go func() {
					startTLSClient(s.tester, proto, addr, s.multicore, s.async)
					atomic.AddInt32(&s.clientActive, -1)
				}()
			}
		}
	}
	if atomic.LoadInt32(&s.clientActive) == 0 {
		var streamAddrs int
		for _, addr := range s.addrs {
			if !strings.HasPrefix(addr, "udp") {
				streamAddrs++
			}
		}
		streamConns := s.nclients * streamAddrs
		disconnected := atomic.LoadInt32(&s.disconnected)
		if int(disconnected) == streamConns && disconnected == atomic.LoadInt32(&s.connected) {
			action = Shutdown
			s.workerPool.Release()
			require.EqualValues(s.tester, 0, s.eng.CountConnections())
		}
	}
	return
}

func TestTLSServer(t *testing.T) {
	go func() {
		err := http.ListenAndServe(":7070", nil)
		if nil != err {
			t.Logf("listenAndServe http error: %v\n", err)
		}
	}()

	t.Run("poll-LT", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, false, false, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, false, false, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, false, false, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, false, false, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, false, false, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, false, false, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-ET", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, true, false, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, true, false, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, true, false, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, true, false, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, true, false, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, true, false, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-reuseport-LT", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, false, true, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, false, true, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, false, true, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, false, true, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, false, true, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, false, true, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-reuseport-ET", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, true, true, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, true, true, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, true, true, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, true, true, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("tcp-async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991"}, true, true, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9992"}, true, true, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-multi-addrs-LT", func(t *testing.T) {
		t.Run("sync", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, false, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, false, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("sync-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, false, false, false, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, false, true, false, true, 10, LeastConnections)
			})
		})
		t.Run("async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, false, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, false, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, false, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, false, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-multi-addrs-reuseport-LT", func(t *testing.T) {
		t.Run("sync", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, true, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, true, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("sync-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, true, false, false, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, true, true, false, true, 10, LeastConnections)
			})
		})
		t.Run("async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, true, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, true, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, false, true, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, false, true, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-multi-addrs-ET", func(t *testing.T) {
		t.Run("sync", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, false, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, false, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("sync-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, false, false, false, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, false, true, false, true, 10, LeastConnections)
			})
		})
		t.Run("async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, false, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, false, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, false, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, false, true, true, true, 10, LeastConnections)
			})
		})
	})

	t.Run("poll-multi-addrs-reuseport-ET", func(t *testing.T) {
		t.Run("sync", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, true, false, false, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, true, true, false, false, 10, LeastConnections)
			})
		})
		t.Run("sync-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, true, false, false, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, true, true, false, true, 10, LeastConnections)
			})
		})
		t.Run("async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, true, false, true, false, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, true, true, true, false, 10, LeastConnections)
			})
		})
		t.Run("async-writev", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9991", "tcp://:9992"}, true, true, false, true, true, 10, RoundRobin)
			})
			t.Run("N-loop", func(t *testing.T) {
				runTLSServer(t, []string{"tcp://:9995", "tcp://:9996"}, true, true, true, true, true, 10, LeastConnections)
			})
		})
	})
}

func runTLSServer(t *testing.T, addrs []string, et, reuseport, multicore, async, writev bool, nclients int, lb LoadBalancing) {
	internal := &testServer{
		tester:     t,
		addrs:      addrs,
		multicore:  multicore,
		async:      async,
		writev:     writev,
		nclients:   nclients,
		workerPool: goPool.Default(),
	}

	ts := &testTLSServer{
		testServer: internal,
	}

	var err error
	if len(addrs) > 1 {
		err = Rotate(ts,
			addrs,
			WithEdgeTriggeredIO(et),
			WithLockOSThread(async),
			WithMulticore(multicore),
			WithReusePort(reuseport),
			WithTicker(true),
			WithTCPKeepAlive(time.Minute),
			WithTCPNoDelay(TCPDelay),
			WithLoadBalancing(lb),
			WithTLSConfig(getServerConfig()))
	} else {
		err = Run(ts,
			addrs[0],
			WithEdgeTriggeredIO(et),
			WithLockOSThread(async),
			WithMulticore(multicore),
			WithReusePort(reuseport),
			WithTicker(true),
			WithTCPKeepAlive(time.Minute),
			WithTCPNoDelay(TCPDelay),
			WithLoadBalancing(lb),
			WithTLSConfig(getServerConfig()))
	}
	assert.NoError(t, err)
}

func startTLSClient(t *testing.T, network, addr string, multicore, async bool) {
	rand.Seed(time.Now().UnixNano())
	c, err := tls2.Dial(network, addr, getGoClientTLSConfig())
	require.NoError(t, err)
	defer c.Close()
	rd := bufio.NewReader(c)
	if network != "udp" {
		msg, err := rd.ReadBytes('\n')
		require.NoError(t, err)
		require.Equal(t, string(msg), "sweetness\r\n", "bad header")
	}
	duration := time.Duration((rand.Float64()*2+1)*float64(time.Second)) / 2
	logging.Debugf("test duration: %v", duration)
	start := time.Now()
	for time.Since(start) < duration {
		reqData := make([]byte, streamLen)
		if network == "udp" {
			reqData = reqData[:datagramLen]
		}
		_, err = rand.Read(reqData)
		require.NoError(t, err)
		_, err = c.Write(reqData)
		require.NoError(t, err)
		respData := make([]byte, len(reqData))
		_, err = io.ReadFull(rd, respData)
		require.NoError(t, err)
		if !async {
			// require.Equalf(t, reqData, respData, "response mismatch with protocol:%s, multi-core:%t, content of bytes: %d vs %d", network, multicore, string(reqData), string(respData))
			require.Equalf(
				t,
				reqData,
				respData,
				"response mismatch with protocol:%s, multi-core:%t, length of bytes: %d vs %d",
				network,
				multicore,
				len(reqData),
				len(respData),
			)
		}
	}
}
