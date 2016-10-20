package utils

import (
	"os"
	"io"
	"github.com/pandazhuzi/buns/errors"
	"path/filepath"
)

func FileExists(root string) bool {

	_,err := os.Stat(root)

	if(err == nil){
		return true
	}

	return false
}


func FileCopy(source string, targer string) error {


	i,err := os.Stat(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	if(i.IsDir()){
		return errors.MakeError("can not copy file, %v is a dir", source)
	}

	if(FileExists(targer)){
		return errors.MakeError("can not copy %v to %v because targer is exist", source,targer)
	}

	fsource, err := os.Open(source)

	if(err != nil){
		return errors.MakeError(err)
	}

	defer fsource.Close()

	ftarger, err := os.OpenFile(targer,os.O_WRONLY|os.O_CREATE,0755)

	if(err != nil){
		return errors.MakeError(err)
	}

	defer ftarger.Close()

	_, err = io.Copy(ftarger,fsource)

	if(err != nil){
		return errors.MakeError(err)
	}

	return err

}

func FileCopyToFolder(source string, targer string) error {
	base := filepath.Base(source)
	targer = filepath.Join(targer,base)

	return FileCopy(source,targer)
}

func innerCopy(source string, targer string) error {

	i, err := os.Stat(source)

	i.Mode()

	if(err != nil){
		return err
	}

	if(i.IsDir()){
		return os.MkdirAll(targer,0755)
	}

	return FileCopy(source,targer)

}

func FolderCopy(source string, target string) error {

	i,err := os.Stat(source)

	if(err != nil){
		return err
	}

	if(!i.IsDir()){
		return errors.MakeError("can not copy folder because %v is not dir.", source)
	}

	if(FileExists(target)){
		i,err := os.Stat(target)

		if(err != nil && !os.IsNotExist(err)){
			return err
		}

		if(!i.IsDir()){
			return errors.MakeError("cant not copy folder %v to %v because targer is not dir",source, target)
		}
	}else{
		err = os.MkdirAll(target, 0755)
	}


	err = filepath.Walk(source, func(visit string, info os.FileInfo, err error) error{

		remain := visit[len(source):]

		if(len(visit) == 0){
			return nil
		}

		return innerCopy(visit, filepath.Join(target, remain))
	})

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil
}

