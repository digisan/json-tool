package jsontool

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
)

// Fmt : use Nodejs 'JSON.stringify' to format JSON with original fields order
func Fmt(bytes []byte, indent string) []byte {
	var m any
	err := json.Unmarshal(bytes, &m)
	failOnErr("%v", err)
	bytes, err = json.MarshalIndent(&m, "", indent)
	failOnErr("%v", err)
	return bytes
}

// FmtStr : use Nodejs 'JSON.stringify' to format JSON with original fields order
func FmtStr(str, indent string) string {
	return string(Fmt([]byte(str), indent))
}

// TryFmt : use Nodejs 'JSON.stringify' to format JSON with original fields order
func TryFmt(bytes []byte, indent string) []byte {
	var m any
	if err := json.Unmarshal(bytes, &m); err != nil {
		return bytes
	}
	bytes, err := json.MarshalIndent(&m, "", indent)
	failOnErr("%v", err)
	return bytes
}

// TryFmtStr : use Nodejs 'JSON.stringify' to format JSON with original fields order
func TryFmtStr(str, indent string) string {
	return string(TryFmt([]byte(str), indent))
}

//////////////////////////////////////////////////////////////////////
// *** NodeJS is required for below field order reserved Format *** //
//////////////////////////////////////////////////////////////////////

const srcFmtJS = `
const strJSON = process.argv[2];
console.log(JSON.stringify(JSON.parse(strJSON), null, 2));
`

const srcFmtFileJS = `
const fs = require('fs');
const fPath = process.argv[2];
try {
    const data = fs.readFileSync(fPath, 'utf8');
    const str = JSON.stringify(JSON.parse(data), null, 2)
    console.log(str)
} catch (e) {
    console.error(e);
}
`

// func _FmtJS(str string) (string, error) {
// 	jsFmt := "fmt.js"
// 	lk.FailOnErrWhen(!fd.FileExists(jsFmt), "%v", fmt.Errorf("%v is not found", jsFmt))
// 	cmd := exec.Command("node", jsFmt, str)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return "", err
// 	}
// 	return string(output), nil
// }

func FmtFileJS(fPath string) (string, error) {
	jsFmtSrc := "fmt.js"
	defer os.RemoveAll(jsFmtSrc)
	fd.MustWriteFile(jsFmtSrc, []byte(srcFmtFileJS))
	cmd := exec.Command("node", jsFmtSrc, fPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return ConstBytesToStr(output), nil
}

func FmtJS(str string) (string, error) {
	temp := "nodejs-fmt-temp.json"
	defer os.RemoveAll(temp)
	fd.MustWriteFile(temp, StrToConstBytes(str))
	return FmtFileJS(temp)
}
