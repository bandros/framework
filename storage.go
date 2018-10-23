package framework

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
)

func StorageUpload(file *os.File,bucket,filename string) (string,error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "",err
	}

	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	defer wc.Close()
	_,err = io.Copy(wc,file);
	if  err != nil {
		return "",err
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s",bucket,filename)
	return url,nil
}
