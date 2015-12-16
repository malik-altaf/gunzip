//This tools is an exercise to develop a utility to download a zip file and extract its contents to disk
package main

import (
	"fmt"
	"flag"
	 "./gunzip"
)

var srcUrl = flag.String("u","https://jdbc.postgresql.org/download/postgresql-jdbc-9.4-1206.src.tar.gz","URL pointing to tar.gz file")
var destinationDir = flag.String("d","D:\\test","destination directory")
var help = flag.Bool("h",false,"Print usage information")

func main() {	

	defer recoverPanic()
	flag.Parse();
	
	if *help {
		flag.PrintDefaults();
		return;
	}
	
	fmt.Printf("Downloading \"%s\" to \"%s\"\n", *srcUrl, *destinationDir)
	
	tz := gunzip.Tarzip{Url:*srcUrl,Dest:*destinationDir}
	error := tz.Extract();
	checkError(error)
	
	fmt.Println("Done!")
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func recoverPanic() {
	if err := recover(); err != nil {
		fmt.Printf("Recovering from panic:%s\n", err)
	}
}
