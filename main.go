package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	Banner()
	dirGenerate, dirReception := ReadEnviroment()
	if len(dirGenerate) == 0 || len(dirReception) == 0 {
		ExitProgram("\tAlguna variable de entorno no esta especificada.")
		return
	}
	homeDir, _ := os.UserHomeDir()
	desGenerate := homeDir + "\\Documents\\GoFind\\Generacion\\"
	desReception := homeDir + "\\Documents\\GoFind\\Recepcion\\"
	choice := 0

	for choice != 3 {
		choice = ChoiceMenu()
		filename := ""
		filesFind := 0
		dirOpen := ""
		var err error = nil
		switch choice {
		case 1:
			filename = ReadFileName()
			filesFind, err = SearchFiles(dirGenerate, filename, desGenerate)
			if filesFind == 0 && err == nil {
				filesFind, err = SearchFiles(dirGenerate+"RESPALDO"+string(os.PathSeparator), filename, desGenerate)
			}
			dirOpen = desGenerate + filename
		case 2:
			filename = ReadFileName()
			filesFind, err = SearchFiles(dirReception, filename, desReception)
			if filesFind == 0 && err == nil {
				filesFind, err = SearchFiles(dirReception+"RESPALDOA"+string(os.PathSeparator), filename, desReception)
			}
			dirOpen = desReception + filename
		case 3:
			continue
		}
		if err != nil {
			fmt.Println("==================================================================")
			fmt.Printf("\tOcurio un error inesperado\nError: %s\n", err.Error())
			fmt.Println("==================================================================")
			choice = 3
		} else {
			if filesFind > 0 {
				exec.Command("explorer", dirOpen).Start()
			}
			fmt.Println("==================================================================")
			fmt.Printf("\tSe encontraron %d archivos.\n", filesFind)
			fmt.Println("==================================================================")
		}
	}
	ExitProgram("\t\tSaliendo de la aplicacion")
}

func Banner() {
	fmt.Println("\t" + `   ______      _______           __`)
	fmt.Println("\t" + `  / ____/___  / ____(_)___  ____/ /`)
	fmt.Println("\t" + ` / / __/ __ \/ /_  / / __ \/ __  / `)
	fmt.Println("\t" + `/ /_/ / /_/ / __/ / / / / / /_/ /  `)
	fmt.Println("\t" + `\____/\____/_/   /_/_/ /_/\__,_/   `)
	fmt.Println("==================================================================")
}

func ExitProgram(message string) {
	fmt.Println(message)
	fmt.Println("==================================================================")
	fmt.Println("\t\tgracias por usar la aplicacion")
	fmt.Println("\t\tBy - Yosimar Zahid Aquino Sosa")
	fmt.Println("==================================================================")
	fmt.Println("\tPresiona Enter para finalizar el programa...")
	fmt.Scanln()
}

func ChoiceMenu() int {
	var choice int
	var successchoice bool = false
	for !successchoice {
		fmt.Println("Por favor seleccione en donde desea realizar la busqueda:")
		fmt.Println("1 -> Generacion de Archivos")
		fmt.Println("2 -> Recepcion de Archivos")
		fmt.Println("3 -> salir")
		fmt.Print("->")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("==================================================================")
			fmt.Println("\tSe esta proporcionando una valor invalido como opcion.")
			fmt.Println("==================================================================")
			continue
		} else {
			if choice <= 0 || choice > 3 {
				fmt.Println("==================================================================")
				fmt.Println("\tLas opciones permitidas son 1, 2 y 3.")
				fmt.Println("==================================================================")
				continue
			}
		}
		successchoice = true
	}
	return choice
}

func ReadEnviroment() (string, string) {
	return os.Getenv("FileGenerate"), os.Getenv("FileReception")
}

func ReadFileName() string {
	var filename string
	var successfilename bool = false
	for !successfilename {
		fmt.Println("Por favor ingrese el nombre del archivo a buscar")
		fmt.Print("->")
		_, err := fmt.Scanln(&filename)
		if err != nil {
			fmt.Println("==================================================================")
			fmt.Println("\tEl nombre del archivo solo pueden ser caracteres alfanumericos.")
			fmt.Println("==================================================================")
			continue
		}
		successfilename = true
	}
	return filename
}

func SearchFiles(origin string, filename string, destiny string) (int, error) {
	dirRead, err := os.Open(origin)
	if err != nil {
		return 0, err
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return 0, err
	}
	counter := 0
	for fileIndex := range dirFiles {
		currentFile := dirFiles[fileIndex]
		dirDestiny := destiny + strings.Split(currentFile.Name(), ".")[0] + string(os.PathSeparator)
		if strings.Split(currentFile.Name(), ".")[0] == filename {
			os.MkdirAll(dirDestiny, 0666)
			originFile, err := os.Open(origin + currentFile.Name())
			if err != nil {
				return 0, err
			}
			defer originFile.Close()
			destinationFile, err := os.Create(dirDestiny + currentFile.Name())
			if err != nil {
				return 0, err
			}
			defer destinationFile.Close()
			_, err = io.Copy(destinationFile, originFile)
			if err != nil {
				return 0, err
			}
			err = destinationFile.Sync()
			if err != nil {
				return 0, err
			}
			counter++
		}
	}
	defer dirRead.Close()
	return counter, nil
}
