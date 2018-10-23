package framework

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"os"
	"path"
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
	filename = uuid.Must(uuid.NewV4()).String() + path.Ext(filename)
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	defer wc.Close()
	_,err = io.Copy(wc,f);
	if  err != nil {
		return "",err
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s",bucket,filename)
	return url,nil
}

func StorageUploadFile(file *multipart.FileHeader,bucket,filename string) (string,error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "",err
	}
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	wc.ContentType = file.Header.Get("Content-Type")
	wc.CacheControl = "public, max-age=86400"
	defer wc.Close()
	fh,err := file.Open()
	if err != nil {
		return "",err
	}
	_,err = io.Copy(wc,fh);
	if  err != nil {
		return "",err
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s",bucket,filename)
	return url,nil
}
