package sample

import "github.com/bandros/framework"

func main() {
	fw := framework.Init{}
	fw.Get()
	r := fw.Begin
	router.Init(r)
	//log.Fatal(fw.RunTls("bandros.tss.my.id"))
	fw.Run()
	//http.Handle("/", r)
	//appengine.Main()
}
