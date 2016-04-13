package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	ChanMode   = "channel"
	StringMode = "string"
)

type Tar struct {
	Gz      bool
	DoPath  bool
	Dest    string
	OriPath string
	Prepare sync.WaitGroup
	Job     sync.WaitGroup

	Tw *tar.Writer

	Ch        chan Content
	JobDone   chan bool
	CloseDone chan bool
}

type Content struct {
	Text     []byte
	FileName string
	Mode     int64
	ModTime  time.Time
}

func (t *Tar) Init(nums ...int) {
	var length int = 10
	if len(nums) != 0 {
		length = nums[0]
	}
	t.Prepare.Add(1)
	t.init(length)
	t.Prepare.Wait()
}

func (t *Tar) init(num int) error {
	t.Ch = make(chan Content, num)
	t.JobDone = make(chan bool, num)
	t.CloseDone = make(chan bool, 2)
	t.Prepare.Done()
	return nil
}

func (t *Tar) Run(mode string) error {
	switch mode {
	case ChanMode:
		t.Prepare.Add(1)
		go t.ChanTar()
		t.Prepare.Wait()
	default:
		return fmt.Errorf("invalid mode")
	}
	return nil
}

func (t *Tar) ChanTar() {
	defer sendBoolChan(t.CloseDone, true)
	fw, err := os.Create(t.Dest)
	if err != nil {
		return
	}
	defer fw.Close()
	// if t.Gz {
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	t.Tw = tar.NewWriter(gw)
	// }
	defer t.Tw.Close()

	dir, err := os.Open(t.OriPath)
	if err != nil {
		return
	}
	defer dir.Close()
	t.Prepare.Done()

	for {
		var done bool = false
		select {
		case content := <-t.Ch:
			header := new(tar.Header)
			header.Name = content.FileName
			header.Size = int64(len(content.Text))
			header.Mode = content.Mode
			if year, _, _ := content.ModTime.Date(); year == 1 {
				header.ModTime = time.Now()
			}

			if err := t.Tw.WriteHeader(header); err != nil {
				// todo error
				fmt.Println(err.Error())
				return
				// break
			} else {
				fmt.Println("write header success")
			}
			if _, err := t.Tw.Write(content.Text); err != nil {
				fmt.Println(err.Error())
				return
				// break
			} else {
				fmt.Println("wirte content success")
			}
			t.Job.Done()
		case <-t.JobDone:
			done = true
			break
		}
		if done {
			break
		}
	}
	return
}

func sendBoolChan(c chan<- bool, v bool) {
	c <- v
}

func (t *Tar) AddFile(content Content) {
	t.Job.Add(1)
	t.Ch <- content
}

func (t *Tar) JobWait() {
	t.Job.Wait()
}

// todo: timeout
func (t *Tar) Done() {
	t.JobWait()
	t.JobDone <- true
	<-t.CloseDone
}
