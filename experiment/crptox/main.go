package main

import (
	"encoding/pem"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"gitlab.com/yawning/obfs4.git/common/ntor"
)

var PrivatePem = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VuBCIEIHD6B6kjsP3/15f1msT03+K1FtX27Msz7im8/JSYny5n
-----END PRIVATE KEY-----`
var PublicPem = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VuAyEAPZKBXxihzfmRHfsLpJwlhpJ+YZYulPsZo6YwzBC6BSg=\n-----END PUBLIC KEY-----"

func main() {
	prBlock, _ := pem.Decode([]byte(PrivatePem))
	//if err != nil {
	//	log.Fatal(err)
	//}
	key := prBlock.Bytes[len(prBlock.Bytes)-32:]
	pp.Println(key)
	puBlock, _ := pem.Decode([]byte(PublicPem))
	//pp.Println(puBlock)
	privKey, _ := ntor.NewPublicKey(key)
	fmt.Print("private key: ")
	fmt.Println(privKey.Hex())
	pubKey := puBlock.Bytes[len(puBlock.Bytes)-32:]
	//pp.Println(pubKey)
	pub, _ := ntor.NewPublicKey(pubKey)
	fmt.Print("public key: ")
	fmt.Println(pub.Hex())
	pp.Println(pub)
	//privKeyHex := "70FA07A923B0FDFFD797F59AC4F4DFE2B516D5F6ECCB33EE29BCFC94989F2E67"
	//fmt.Println(privKeyHex)
}
