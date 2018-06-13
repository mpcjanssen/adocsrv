package adoc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"io"

	"github.com/husobee/vestigo"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
 <link rel="stylesheet" type="text/css" href="/assets/asciidoctor-default.css">
		<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.0.13/css/all.css" integrity="sha384-DNOHZ68U8hZfKXOrtjWvjxusGo9WQnrNx2sqG0tfsghAvtVlRW3tvkXWZh58N9jp" crossorigin="anonymous">
		<meta charset="UTF-8">
	</head>
	<body>
<body class="article">
<div id="content">
<div class="paragraph">
<table class="tableblock frame-all grid-all stretch">
<colgroup>
<col style="width: 5%;">
<col style="width: 85%;">
<col style="width: 10%;">
</colgroup>
<tbody>
		<tr>
		<td class="tableblock halign-left valign-top">
		<a href='/browse/{{ .Parent }}'><i class="fas fa-folder"></i></a>
		</td>
		<td class="tableblock halign-left valign-top">
		..
		</td>
		<td></td>
		</tr>
		{{range .Folders}}
		<tr>
		<td class="tableblock halign-left valign-top">
		<a href='/browse/{{ .Path }}'><i class="fas fa-folder"></i></a>
		</td>
		<td class="tableblock halign-left valign-top">
		{{ .Name }}
		</td>
		<td></td>
		</tr>
		{{end}}


		<tr>
		{{range .ADOCs}}
		<tr>
		<td class="tableblock halign-left valign-top">
		</td>
		<td class="tableblock halign-left valign-top">
		<a href='/view/{{ .Path  }}'>{{ .Name }}</a>
		</td>
		<td class="tableblock halign-left valign-top">
		<a href='/edit/{{ .Path  }}'><i class="fas fa-edit"></i></a> 
		<a href='/reveal/{{ .Path  }}'><i class="fas fa-file-powerpoint"></i></a>
		</td>
		</tr>
		{{end}}
</tbody>
</table>

</div>
</div>
	</body>
</html>`

type fInfo struct {
	Name string
	Path string
}

func BrowseHandler(w http.ResponseWriter, r *http.Request) {

	var filesToShow struct {
		Parent  string
		Folders []fInfo
		ADOCs   []fInfo
	}
	processPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	relPath := vestigo.Param(r, "_name")
	if relPath == "browse/" {
		relPath = "."
	}
	files, err := ioutil.ReadDir(filepath.Join(processPath, relPath))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	if relPath != "." {
		filesToShow.Parent = filepath.Join(relPath, "..")
	}

	for _, f := range files {
		if f.IsDir() {
			filesToShow.Folders = append(filesToShow.Folders, fInfo{Name: f.Name(), Path: filepath.Join(relPath, f.Name())})
		} else if strings.HasSuffix(f.Name(), ".adoc") {
			filesToShow.ADOCs = append(filesToShow.ADOCs, fInfo{Name: f.Name(), Path: filepath.Join(relPath, f.Name())})
		}
	}
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	err = t.Execute(w, filesToShow)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {

	processPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	filePath := filepath.Join(processPath, vestigo.Param(r, "_name"))
	if strings.HasSuffix(filePath, ".adoc") {
		cmd := exec.Command("asciidoctor", "-a", "imagesdir=./", "-r", "asciidoctor-diagram", "-q", "-o", "-", filePath)
		cmd.Stdout = w
		cmd.Stderr = w
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(w, err)
		}
	} else {
		fileReader, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			io.Copy(w, fileReader)
			fileReader.Close()
		}
	}
}

func RevealHandler(w http.ResponseWriter, r *http.Request) {

	processPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	filePath := filepath.Join(processPath, vestigo.Param(r, "_name"))
	cmd := exec.Command("asciidoctor-revealjs", "-a", "revealjsdir=/assets/reveal.js", "-o", "-", filePath)
	cmd.Stdout = w
	cmd.Stderr = w
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	processPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	filePath := filepath.Join(processPath, vestigo.Param(r, "_name"))
	cmd := exec.Command("gvim", "--servername", "ADOC", "--remote-silent", filePath)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Start()
	http.Redirect(w, r, r.Referer(), 301)
}
