package framework

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
)

func StorageUpload(file string,bucket,filename string) (string,error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "",err
	}

	f,err := os.Open(file)
	if err != nil {
		return "",err
	}
	defer f.Close()
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	defer wc.Close()
	_,err = io.Copy(wc,f);
	if  err != nil {
		return "",err
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s",bucket,filename)
	return url,nil
}
