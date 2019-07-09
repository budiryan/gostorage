package gostorage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	gStorage "cloud.google.com/go/storage"

	"github.com/budiryan/gostorage/storage"
)

func initGoogleClient() (storage.Storage, error) {
	client, err := storage.NewStorage(storage.GCP, storage.GCPStorage(context.Background(), "GCP STORAGE BUCKET", "PATH TO SECRET FILE"))
	if err != nil {
		log.Println("error initializing: ", err)
		return nil, err
	}

	return client, nil
}

func ExampleGoogle_Read() {
	client, err := initGoogleClient()
	if err != nil {
		return
	}

	// Reading with timeout context (completely optional)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	readCloser, err := client.Read("FILEPATH + NAME", storage.OperationCtx(ctx))
	if err != nil {
		log.Println("error reading file: ", err)
		return
	}
	defer readCloser.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(readCloser)
	fmt.Println(buf.String())
}

func ExampleGoogle_Write() {
	client, err := initGoogleClient()
	if err != nil {
		return
	}

	// Writing with timeout context (completely optional)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	writeCloser, err := client.Write("FILEPATH + NAME", storage.OperationCtx(ctx))
	if err != nil {
		log.Println("error initializing writer: ", err)
	}

	defer writeCloser.Close()
	_, err = writeCloser.Write([]byte("jon bodat"))
	if err != nil {
		log.Println("error writing: ", err)
	}
}

func ExampleGoogle_IsExists() {
	client, err := initGoogleClient()
	if err != nil {
		return
	}

	// Checking existence with timeout context (completely optional)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	exists, err := client.IsExists("FILEPATH + NAME", storage.OperationCtx(ctx))
	if err != nil {
		log.Println("error checking file existence: ", err) // e.g.: because of context timeout
	}

	fmt.Println(exists)
}

func ExampleGoogle_GetSignedURL() {
	client, err := initGoogleClient()
	if err != nil {
		return
	}

	url, err := client.GetSignedURL("FILE", &storage.SignedURLOptions{
		HTTPMethod:  http.MethodGet,
		ContentType: "",
		ExpiryTime:  time.Now().Add(time.Minute * 10),
	})
	if err != nil {
		log.Println("error getting signed url: ", err)
	}
	fmt.Println(url)
}


func ExampleGoogle_ListObject() {
	client, err := initGoogleClient()
	if err != nil {
		return
	}

	res, err := client.ListObject(&gStorage.Query{
		Prefix:    "Prefix",
		Delimiter: "Delimiter",
	})
	if err != nil {
		log.Println("error getting list object: ", err)
	}
	fmt.Println(res)
}
