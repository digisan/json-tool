package jsontool

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	gotkio "github.com/digisan/gotk/io"
)

func TestScanObjectInArray(t *testing.T) {
	// set up a context to manage ingest pipeline
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// file, err := os.OpenFile("/home/qmiao/Desktop/rrd.json", os.O_RDONLY, os.ModePerm)
	file, err := os.Open("./data/mixed.json")
	if err != nil {
		fPln(err)
	}

	cOut, cErr, err := ScanObjectInArray(ctx, file, true)
	if err != nil {
		fmt.Println("Not Valid Array")
		return
	}

	go func() {
		I := 1
		for out := range cOut {
			// if I == 3 {
			// 	cancelFunc()
			// }
			if I > 0 {
				fPln(I, "cOut")
				gotkio.MustWriteFile(fSf("dump%02d.json", I), []byte(out))
			}
			I++
		}
	}()

	go func() {
		I := 1
		for e := range cErr {
			if e != nil {
				panic("channel error")
			}
			fPln(I, "cErr")
			I++
		}
	}()

	time.Sleep(1 * time.Second)
}

func TestScanObject(t *testing.T) {

	// set up a context to manage ingest pipeline
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// file, err := os.OpenFile("/home/qmiao/Desktop/rrd.json", os.O_RDONLY, os.ModePerm)
	file, err := os.Open("./data/mixed.json")
	if err != nil {
		fPln(err)
	}

	mustarray := false
	if cOut, ja := ScanObject(ctx, file, mustarray, true, OUT_FMT); !ja && mustarray {

		fPln("NOT JSON array")

	} else {

		I := 1
		for rst := range cOut {

			if rst.Err != nil {
				panic("Not Valid@" + rst.Err.Error())
			}

			if I == 3 {
				cancelFunc()
			}

			if I > 0 {
				fPln(I)
				// fPln(rst.Obj)
				gotkio.MustWriteFile(fSf("dump%02d.json", I), []byte(rst.Obj))
			}

			I++
		}

		// for {
		// 	if rst, more := <-chRst; more {
		// 		fPln(I)
		// 		fPln(rst.Obj)
		// 		fPln(rst.Err)
		// 		I++
		// 	} else {
		// 		break
		// 	}
		// }
	}
}

func Test_analyse(t *testing.T) {

	l1 := `[  {`
	l2 := `"Id": 1,`
	l3 := ` "Name": "Ahmad,Ahmad",`
	l4 := `"Age": "21"`
	l5 := `  },  {"Id": 2,    "Name": "","Age": "50"},{"Id": 3,"Name": "Test1","Age": ""},  {`
	l6 := `"Id": 4 } ]`

	fPln(analyseJL(l1, 0))
	fPln(analyseJL(l2, 1))
	fPln(analyseJL(l3, 1))
	fPln(analyseJL(l4, 1))
	fPln(analyseJL(l5, 1))
	fPln(analyseJL(l6, 1))
}
