package jsontool

import (
	"os"
	"testing"
	"time"

	"github.com/digisan/gotk"
	gotkio "github.com/digisan/gotk/io"
)

func TestJSONBreakArrCont(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	bytes, err := os.ReadFile("./data/Activities.json")
	failOnErr("%v", err)
	jsonstr := string(bytes)

	values, ok := BreakArr(jsonstr)
	fPln(ok)
	for _, v := range values {
		fPln(v)
	}
}

func TestJSONBreakBlkContV2(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	if bytes, err := os.ReadFile("./why.json"); err == nil {
		// jsonstr := JSONBlkFmt(string(bytes), "  ")
		jsonstr := string(bytes)
		_, cont := SglEleBlkCont(jsonstr)
		names, values := BreakMulEleBlkV2(cont)
		for i, name := range names {
			fPln(MkSglEleBlk(name, values[i], false))
			fPln(" ------------------------------------------ ")
		}
	}
}

func TestScanArray2Objects(t *testing.T) {

	// file, err := os.OpenFile("/home/qmiao/Desktop/rrd.json", os.O_RDONLY, os.ModePerm)
	file, err := os.OpenFile("./data/data.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		fPln(err)
	}

	mustarray := false
	if chRst, ja := ScanObject(file, mustarray, true, OUT_ORI); !ja && mustarray {

		fPln("NOT JSON array")

	} else {

		I := 1
		for rst := range chRst {

			if rst.Err != nil {
				panic("Not Valid@" + rst.Err.Error())
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
