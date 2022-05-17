package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"time"
)

type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (user User) GetEmail() string {
	return user.Email
}

func (user User) GetRegistration() *registration.Resource {
	return user.Registration
}
func (user User) GetPrivateKey() crypto.PrivateKey {
	return user.key
}

var (
	//let Encrypt 正式环境 API
	letEncryptDirectoryProdUrl = "https://acme-v02.api.letsencrypt.org/directory"
	//let Encrypt 测试环境 API
	letEncryptDirectoryTestUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

type MyDNS struct {
}

func (myDNS MyDNS) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)

	fmt.Printf("fqdn: %s,  value: %s\n", fqdn, value)
	return nil
}

func (myDNS MyDNS) CleanUp(domain, token, keyAuth string) error {
	return nil
}

func (myDNS MyDNS) Timeout() (timeout, interval time.Duration) {
	return 120 * time.Second, 5 * time.Second
}

func main() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("private key err: ", err.Error())
		return
	}

	user := User{
		Email: "qx@startops.com.cn",
		key:   privateKey,
	}

	config := lego.NewConfig(&user)

	config.CADirURL = letEncryptDirectoryTestUrl
	config.Certificate.KeyType = certcrypto.RSA4096

	client, err := lego.NewClient(config)

	if err != nil {
		fmt.Println("config key err: ", err.Error())
		return
	}

	var myDNS MyDNS
	err = client.Challenge.SetDNS01Provider(myDNS)
	if err != nil {
		fmt.Println("set dns err: ", err.Error())
		return
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		fmt.Println("reg err: ", err.Error())
		return
	}

	user.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{"hook.startops.com.cn"},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		fmt.Println("cert err: ", err.Error())
		return
	}

	fmt.Println("Domain: ", certificates.Domain)
	fmt.Println("\n")
	fmt.Println("IssuerCertificate: ", string(certificates.IssuerCertificate))
	fmt.Println("\n")
	fmt.Println("PrivateKey: ", string(certificates.PrivateKey))
	fmt.Println("\n")
	fmt.Println("Certificate: ", string(certificates.Certificate))
}
