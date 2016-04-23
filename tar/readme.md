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
		Text:     []byte("test\ntest\ntest"),
		FileName: "test.log",
		Mode:     0755,
	}
	tar1.AddFile(tmp)
	tar1.Done()
	/////////////

	bf := new(bytes.Buffer)
	log := bytes.NewBuffer([]byte(""))
	log.WriteString("test\n")
	log.WriteString("congratulation\n")
	log.WriteString("you\n")
	tar2 := tar.Tar{
		Gz: true,
		// Dest: "ng.tgz",
		// OriPath: log.String(),
	}
	tar2.Init()
	if err := tar2.Run(tar.IoMode, bf); err != nil {
		fmt.Println("run tar2 fail:", err.Error())
		return
	}
	tmp2 := tar.Content{
		Text:     log.Bytes(),
		FileName: "test.log",
		Mode:     0755,
	}
	tar2.AddFile(tmp2)
	tar2.Done()
	tar2.Close()
}
