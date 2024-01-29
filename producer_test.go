package logkafka

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"
)

const (
	certFile = "fixtures/client.cer.pem"
	keyFile  = "fixtures/client.key.pem"
	caFile   = "fixtures/ca.pem"
)

func makeTLSConfiguration(t *testing.T) *tls.Config {
	t.Helper()

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		t.Fatalf("makeTLSConfiguration failed %v", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		t.Fatalf("makeTLSConfiguration failed %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}
}

func TestSimpleProducer(t *testing.T) {
	_, err := SimpleProducer([]string{"127.0.0.1:9092"}, 0, 0, nil)
	// without a running kafka this should result in an error
	if err == nil {
		t.Error("SimpleProducer() expected error got nil")
	} else {
		t.Logf("SimpleProducer() error = %v", err)
	}
}

func TestSimpleProducer_WithTLS(t *testing.T) {
	_, err := SimpleProducer([]string{"127.0.0.1:9092"}, 0, 0, makeTLSConfiguration(t))
	// without a running kafka this should result in an error
	if err == nil {
		t.Error("SimpleProducer() expected error got nil")
	} else {
		t.Logf("SimpleProducer() error = %v", err)
	}
}
