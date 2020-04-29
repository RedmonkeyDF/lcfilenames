package main

import (
	"errors"
	"fmt"
	mfu "github.com/RedmonkeyDF/mfmodfileutil"
	"github.com/akamensky/argparse"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetFilenames(adir string) ([]string, error) {

	d, errd := os.Open(adir)
	if errd != nil {
		log.Fatal(errd)
	}

	files, errread := d.Readdir(-1)
	if errread != nil {

		return []string{}, errread
	}

	errclose := d.Close()
	if errclose != nil {

		return []string{}, errclose
	}

	var outfiles []string

	for _, file := range files {

		outfiles = append(outfiles, filepath.Join(adir, file.Name()))
	}

	return outfiles, nil
}

func lcfnames(adir string, dryrun bool) error {

	if adir == "" {

		return errors.New("Cannot operate on an empty directory.")
	}

	opdir, errabs := filepath.Abs(adir)
	if errabs != nil {

		return errors.New(fmt.Sprintf("Unexpected error retrieving absolute path for \"%s\".  Error: %s", adir, errabs))
	}

	de, errde := mfu.DirectoryExists(opdir)
	if errde != nil {

		return errors.New(fmt.Sprintf("Unexpected error checking if directory \"%s\" exists.  Error: %s", opdir, errde))
	}

	if !de {

		return errors.New(fmt.Sprintf("Directory \"%s\" does not existexists.  Error: %s", opdir))
	}

	fnames, errgetnames := GetFilenames(opdir)
	if errgetnames != nil {

		return errors.New(fmt.Sprintf("GetFilenames on \"%s\" returned error: \"%s\".", opdir, errgetnames))
	}

	errren := renfiles(fnames, dryrun)

	if errren != nil {

		return errors.New(fmt.Sprintf("Renaming files on \"%s\" returned error: \"%s\".", opdir, errren))
	}

	return nil
}

func renfiles(flist []string, notdryrun bool) error {

	for _, file := range flist {

		renfname := filepath.Join(filepath.Dir(file), strings.ToLower(filepath.Base(file)))

		if notdryrun {
			errren := os.Rename(file, renfname)
			if errren != nil {

				return errors.New(fmt.Sprintf("Error renaming file \"%s\".  Error: %s.", file, errren))
			}
		} else {

			fmt.Println(fmt.Sprintf("In drectory %s, %s to be renamed to %s", filepath.Dir(file), filepath.Base(file), filepath.Base(renfname)))
		}
	}

	return nil
}

func main() {
	// Create new parser object
	parser := argparse.NewParser("lcfnames", "Renames all files in a directory to lower case.")
	// Create string flag
	lcdir := parser.String("d", "directory", &argparse.Options{Required: true, Help: "Directory to operate on."})
	notdryrun := parser.Flag("x", "execute", &argparse.Options{Required: false, Help: "Actually perform the renaming, not just a dry run.."})
	// Parse input
	errpars := parser.Parse(os.Args)
	if errpars != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(errpars))
		os.Exit(1)
	}

	errlcf := lcfnames(*lcdir, *notdryrun)
	if errlcf != nil {

		fmt.Println(errlcf)
		os.Exit(1)
	}

}