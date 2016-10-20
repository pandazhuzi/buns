package utils

import (
	"html/template"
	"github.com/pandazhuzi/buns/errors"
	"bytes"
	"os"
	"io/ioutil"
	"strings"
	"io"
	"path/filepath"
)

func Render(source string, w io.Writer, data interface{}) error {

	tpl := template.New("template-render")

	tpl, err := tpl.Parse(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	err = tpl.Execute(w, data)

	if err != nil {
		return errors.MakeError(err)
	}

	return nil
}

func RenderString(source string, data interface{}) (string, error){

	buf := bytes.NewBuffer(nil)

	err := Render(source, buf, data)

	if(err != nil){
		return "", err
	}

	return buf.String(),nil
}


func RenderTemplateFile(source string, patten string, data interface{}) error {

	if(len(patten) ==0 ){
		return errors.MakeError("patten is empty.")
	}


	if(!strings.Contains(source,patten)){
		return nil
	}

	tpl, err := ioutil.ReadFile(source)

	if(err!=nil){
		return errors.MakeError(err)
	}

	target := strings.Replace(source,patten,".",-1)


	if(FileExists(target)){
		return errors.MakeError("render targer file %v alerady exist", target)
	}

	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY,0755)

	if(err != nil){
		return errors.MakeError(err)
	}

	defer f.Close()

	err = Render(string(tpl), f, data)

	if(err != nil){
		return errors.MakeError(err)
	}

	err = os.Remove(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil
}


func FolderRenameByTemplate(source string, v interface{}) error {

	err := filepath.Walk(source,func(visit string, info os.FileInfo, err error) error{

		target, err := RenderString(visit, v)

		if err != nil {
			return errors.MakeError(err)
		}

		if(visit != target){
			return os.Rename(visit,target)
		}

		return nil

	})

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil

}

func FolderRenderByTemplate(source string, patten string,data interface{}) error {

	err := filepath.Walk(source,func(visit string, info os.FileInfo, err error) error{

		if (!info.IsDir()){
			err = RenderTemplateFile(visit, patten, data)

			if(err != nil){
				return errors.MakeError(err)
			}
		}
		return nil

	})

	if err != nil {
		return errors.MakeError(err)
	}

	return nil
}