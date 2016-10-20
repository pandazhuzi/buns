package utils

import (
	"github.com/pandazhuzi/buns/errors"
	"path/filepath"
	"strings"
	"os"
	"compress/gzip"
	"archive/tar"
	"path"
	"io/ioutil"
	"io"
)

func tarFile(base,relative string, t *tar.Writer) error {

	source := filepath.Join(base, relative)

	info, err := os.Stat(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	hdr, err := tar.FileInfoHeader(info, "")

	if(err != nil){
		return errors.MakeError(err)
	}

	hdr.Name = relative

	err = t.WriteHeader(hdr)

	if(err != nil){
		return errors.MakeError(err)
	}

	file, err := os.OpenFile(source,os.O_RDONLY, 0755)

	if( err != nil){
		return errors.MakeError(err)
	}

	defer file.Close()

	_, err = io.Copy(t,file)

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil
}

func tarDirectory(base,relative string, t *tar.Writer, name string) error{

	source := filepath.Join(base,relative)

	info, err := os.Stat(source)

	if(err != nil){
		return errors.MakeError(err)
	}


	paths, err := ioutil.ReadDir(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	for _, home := range paths {

		if(home.IsDir()){
			err = tarDirectory(base, filepath.Join(relative, home.Name()), t, "")

			if(err != nil){
				return errors.MakeError(err)
			}

		}else{

			err = tarFile(base, filepath.Join(relative, home.Name()), t)

			if(err != nil){
				return errors.MakeError(err)
			}
		}

	}

	hdr, err := tar.FileInfoHeader(info, "")

	if(err != nil){
		return errors.MakeError(err)
	}

	if name != "" {
		hdr.Name = name
	}else{
		hdr.Name = relative
	}


	err = t.WriteHeader(hdr)

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil
}

func Tar(source string, target string, name string) error{

	source, err := filepath.Abs(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	target , err = filepath.Abs(target)

	if(err != nil){
		return errors.MakeError(err)
	}

	if(strings.Index(target, source) == 0){
		return errors.MakeError("target path %v is in source path %v", target, source)
	}

	if(FileExists(target)){
		err = os.Remove(target)

		if(err != nil){
			return errors.MakeError(err)
		}
	}

	fs, err := os.OpenFile(target,os.O_CREATE|os.O_WRONLY, 0755)

	if(err != nil){
		return errors.MakeError(err)
	}

	defer fs.Close()

	zip := gzip.NewWriter(fs)
	tar := tar.NewWriter(zip)


	stat, err := os.Stat(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	base,relative := path.Split(source)

	if(stat.IsDir()){
		err = tarDirectory(base,relative, tar, name)
	}else{
		err = tarFile(base,relative,tar)
	}

	if err != nil {
		return errors.MakeError(err)
	}

	return nil
}
