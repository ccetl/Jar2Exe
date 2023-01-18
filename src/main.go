/*
Copyright (C) 2023 ccetl (pseudonym)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Config -----------
var everyJar bool = true           //execute every jar
var jarName string = "yourJar.jar" //only if everyJar is false: set the file name of the file to excetute (lowercase sensitive)
var debug bool = true              //send info to the terminal
// -----------------

var tempFolder string = os.TempDir() + "\\src2exe"
var executables []string

func main() {
	if createFolder() == true {
		sendMessage("[main] Created a temporaly folder.")
	} else {
		sendMessage("[main] The creation of a folder failed.")
		return
	}

	if dumpFiles() == true {
		sendMessage("[main] Moved your files.")
	} else {
		sendMessage("[main] No jar found.")
		return
	}

	if executeJar() == true {
		sendMessage("[main] Succsefull exectued you jar.")
	}
}

func executeJar() bool {
	for _, file := range executables {
		cmd := exec.Command("java", "-jar", file)
		output, err := cmd.CombinedOutput()
		sendMessage(tempFolder + "\\" + file)
		if err != nil {
			sendMessage("[executeJar] " + string(output))
			sendError("executeJar", err)
		}
	}
	return true
}

func dumpFiles() bool {
	var paths []string
	files, err := ioutil.ReadDir("./resources")
	sendError("dumbFiles", err)

	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, file.Name())
			sendMessage("[dumpFiles] " + file.Name())
		}
	}
	sendError("dumpFiles", err)

	for _, fileName := range paths {
		var in string = "resources\\" + fileName
		var out string = tempFolder + "\\" + fileName

		inFile, err := os.Open(in)
		sendError("dumpFiles", err)
		defer inFile.Close()

		outFile, err := os.Create(out)
		sendError("dumpFiles", err)
		defer outFile.Close()

		_, err = io.Copy(outFile, inFile)
		sendError("dumpFiles", err)

		dotIndex := strings.LastIndex(fileName, ".")
		extension := fileName[dotIndex:]

		if extension == ".zip" {
			exctractZip(out)
		}

		if everyJar {
			if extension == ".jar" {
				executables = append(executables, out)
			}
		} else {
			if fileName == jarName {
				executables = append(executables, out)
			}
		}
	}

	for _, s := range executables {
		sendMessage(s)
	}

	if len(executables) == 0 {
		return false
	}

	return true
}

func createFolder() bool {
	tempFolder := os.TempDir() + "\\src2exe"
	sendMessage(tempFolder)
	_, err := os.Stat(tempFolder)
	if err == nil {
		err = os.RemoveAll(tempFolder)
		sendError("createFolder", err)
		sendMessage("[createFolder] Deleted " + tempFolder + ".")
	}
	err2 := os.Mkdir(tempFolder, os.ModePerm)
	sendError("createFolder", err2)
	return true
}

func exctractZip(File string) {
	r, err := zip.OpenReader("File")
	sendError("exctractZip", err)
	defer r.Close()

	for _, file := range r.File {
		rc, err := file.Open()
		sendError("exctractZip", err)
		defer rc.Close()

		path := file.Name
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			sendError("exctractZip", err)
			defer f.Close()

			_, err = io.Copy(f, rc)
			sendError("exctractZip", err)
		}
	}
}

func sendError(Function string, Error error) {
	if Error != nil {
		panic("[" + Function + "] ERROR:" + Error.Error())
	}
	return
}

func sendMessage(Message string) {
	if debug == true {
		fmt.Println(Message)
	}
}
