package main

import (
	"github.com/ianwoolf/golib/tar"
)

func main() {
	tar1 := tar.Tar{
		Gz:      true,
		DoPath:  false,
		Dest:    "tar/result.tgz",
		OriPath: "file/",
	}

	tar1.Init()
	tar1.Run(tar.ChanMode)

	tmp := tar.Content{
		Text:     []byte("test tar\ntest tar package!\n"),
		FileName: "test.log",
		Mode:     0755,
	}
	tar1.AddFile(tmp)
	tar1.Done()
}
