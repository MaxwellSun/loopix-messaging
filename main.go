package main

import (
	"flag"
	"anonymous-messaging/client"
	"anonymous-messaging/server"
	"anonymous-messaging/sphinx"
	"anonymous-messaging/pki"
	"fmt"
	"anonymous-messaging/config"
	"github.com/protobuf/proto"
)


const(
	PKI_DIR = "pki/database.db"
)

func pkiPreSetting(pkiDir string) error {
	db, err := pki.OpenDatabase(pkiDir, "sqlite3")
	if err != nil {
		return err
	}
	defer db.Close()

	params := make(map[string]string)
	params["Id"] = "TEXT"
	params["Typ"] = "TEXT"
	params["Config"] = "BLOB"

	err = pki.CreateTable(db, "Pki", params)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	typ := flag.String("typ", "", "A type of entity we want to run")
	id := flag.String("id", "", "Id of the entity we want to run")
	host := flag.String("host", "", "The host on which the entity is running")
	port := flag.String("port", "", "The port on which the entity is running")
	providerId := flag.String("provider", "", "The port on which the entity is running")
	flag.Parse()

	err := pkiPreSetting(PKI_DIR)
	if err != nil{
		panic(err)
	}


	switch *typ {
	case "client":
		db, err := pki.OpenDatabase(PKI_DIR, "sqlite3")

		if err != nil{
			panic(err)
		}


		row := db.QueryRow("SELECT Config FROM Pki WHERE Id = ? AND Typ = ?", providerId, "Provider")

		var results []byte
		err = row.Scan(&results)
		if err != nil {
			fmt.Println(err)
		}
		var providerInfo config.MixConfig
		err = proto.Unmarshal(results, &providerInfo)

		pubC, privC, err := sphinx.GenerateKeyPair()
		if err != nil{
			panic(err)
		}

		client, err := client.NewClient(*id, *host, *port, pubC, privC, PKI_DIR, providerInfo)
		if err != nil{
			panic(err)
		}

		err = client.Start()
		if err != nil{
			panic(err)
		}

	case "mix":
		pubM, privM, err := sphinx.GenerateKeyPair()
		if err != nil{
			panic(err)
		}

		mixServer, err := server.NewMixServer(*id, *host, *port, pubM, privM, PKI_DIR)
		if err != nil{
			panic(err)
		}

		err = mixServer.Start()
		if err != nil{
			panic(err)
		}
	case "provider":
		pubP, privP, err := sphinx.GenerateKeyPair()
		if err != nil{
			panic(err)
		}

		providerServer, err := server.NewProviderServer(*id, *host, *port, pubP, privP, PKI_DIR)
		if err != nil{
			panic(err)
		}

		err = providerServer.Start()
		if err != nil{
			panic(err)
		}
	}
}