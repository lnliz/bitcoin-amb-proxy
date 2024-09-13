package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type AwsSignedTransport struct {
	credsProvider aws.CredentialsProvider
	signer        *v4.Signer
	host          string
	region        string
}

func (t AwsSignedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.Background()

	req.Host = t.host
	req.Header.Del("X-Amz-Date")
	req.Header.Del("X-Forwarded-For")

	payloadHash, err := hexEncodedSha256OfRequest(req)
	if err != nil {
		fmt.Printf("Failed to sign request: %v\n", err)
		return nil, err
	}

	creds, err := t.credsProvider.Retrieve(ctx)
	if err != nil {
		fmt.Printf("Failed to retrieve AWS credentials: %v\n", err)
		return nil, err
	}

	if err = t.signer.SignHTTP(ctx, creds, req, payloadHash, "managedblockchain", t.region, time.Now()); err != nil {
		fmt.Printf("Failed to sign request: %v\n", err)
		return nil, nil
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	return resp, err
}

func NewProxy(credsProvider aws.CredentialsProvider, region string, network string) (*httputil.ReverseProxy, error) {
	host := fmt.Sprintf("%s.bitcoin.managedblockchain.%s.amazonaws.com", network, region)
	u, err := url.Parse("https://" + host)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = AwsSignedTransport{
		credsProvider: credsProvider,
		signer:        v4.NewSigner(),
		host:          host,
		region:        region,
	}

	return proxy, nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	flagListenAddr := flag.String("listen-addr", "0.0.0.0:8787", "listen addr")
	flagRegion := flag.String("region", "us-east-1", "region")
	flagNetwork := flag.String("network", "mainnet", "network")
	flagAwsKey := flag.String("aws-key", "", "aws key")
	flagAwsSecret := flag.String("aws-secret", "", "aws secret")

	flag.Parse()

	var credsProvider aws.CredentialsProvider

	if *flagAwsKey != "" && *flagAwsSecret != "" {
		credsProvider = credentials.NewStaticCredentialsProvider(*flagAwsKey, *flagAwsSecret, "")
	} else {
		/*
			todo: hook up more credential providers
		*/
		credsProvider = ec2rolecreds.New()
	}

	if credsProvider == nil {
		flag.Usage()
		log.Fatalf("Missing AWS creds")
	}

	proxy, err := NewProxy(credsProvider, *flagRegion, *flagNetwork)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v\n", err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))

	log.Printf("Starting server on %s\n", *flagListenAddr)
	if err := http.ListenAndServe(*flagListenAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func hexEncodedSha256OfRequest(r *http.Request) (string, error) {
	if r.Body == nil {
		return `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`, nil
	}

	hasher := sha256.New()

	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	if err := r.Body.Close(); err != nil {
		return "", err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
	hasher.Write(reqBodyBytes)
	digest := hasher.Sum(nil)

	return hex.EncodeToString(digest), nil
}
