package main

import (
	"fmt"
	// "git.apache.org/thrift.git/lib/go/thrift"
	"github.com/fatih/color"
	"line"
)

var cyanBold = color.New(color.FgCyan).Add(color.Bold).SprintFunc()

func printLoginResult(result *line.LoginResult_) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("token: [%v]\n", cyanBold(result.GetAuthToken()))
	fmt.Printf("certificate: [%v]\n", cyanBold(result.GetCertificate()))
	fmt.Printf("pincode: [%v]\n", cyanBold(result.GetPinCode()))
	fmt.Printf("loginResult: [%v]\n", cyanBold(result.GetTypeA1().String()))
	fmt.Printf("verifier: [%v]\n", cyanBold(result.GetVerifier()))
}
