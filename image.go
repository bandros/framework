package framework

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/satori/go.uuid"
	"image"
	"log"
	"mime/multipart"
	"path"
	"time"
)

type Image struct {
	img image.Image
	ext string
	file *multipart.FileHeader
	Encrypt bool
	Width int
	Height int
}

type Filename struct {
	Filename string
	Dir string
	Fullpath string
}


func(img *Image) Set(file *multipart.FileHeader) error {
	img.Encrypt = false
	f,err := file.Open()
	if err!= nil {
		return err
	}
	defer f.Close()
	src, err := imaging.Decode(f)
	if err!= nil {
		return err
	}
	img.img = src
	img.file = file
	return nil
}

func(img *Image) ResizeSave(location string)  (Filename,error){
	var f Filename
	img.img = imaging.Resize(img.img,img.Width,img.Height, imaging.Lanczos)
	filename := img.file.Filename
	if img.Encrypt {
		filename = string(time.Now().Unix())+uuid.Must(uuid.NewV4()).String() + path.Ext(filename)
	}
	f.Dir = location
	f.Filename = filename
	f.Fullpath = location+filename
	err := imaging.Save(img.img, f.Fullpath)
	if err != nil {
		return Filename{},err
	}
	return f,nil
}

func(img *Image) ResizeUpload(bucket string)  (Filename,error){
	var f Filename
	img.img = imaging.Resize(img.img,img.Width,img.Height, imaging.Lanczos)
	filename := img.file.Filename
	filename = string(time.Now().Unix())+uuid.Must(uuid.NewV4()).String() + ".png"
	f.Dir = "https://storage.googleapis.com/"
	f.Filename = filename
	f.Fullpath = f.Dir+filename


	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	//wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	//wc.ContentType = img.file.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	wc.CacheControl = "public, max-age=86400"

	err1 := imaging.Encode(wc, img.img,imaging.PNG )
	if err1 != nil {
		fmt.Println("to bucket ", err1)
	}

	if err := wc.Close(); err != nil {
		fmt.Println("error wc close ", err)
		return Filename{},err
	}

	return f,nil
}
