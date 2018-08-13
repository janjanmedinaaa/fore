package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"regexp"
	"log"
	"time"
	
	"github.com/radovskyb/watcher"
	"github.com/tidwall/gjson"
)

/****** fore home ******/

func home() {
	fmt.Println("Fore v1.0.0")
}

/****** fore compile ******/

func compile() {
	//Get Files of current directory
	files, _ := ioutil.ReadDir(pwd() + "/Fore")
			
	//Loop to each file in current directory
	for _, f := range files {
		if filetypefunc(f.Name()) == "fore" {
			title, imports, styles, scripts, content, genhtmlfile := readfile(pwd(), f.Name())

			htmlcompile(title, imports, styles, scripts, content, genhtmlfile)
		}
	}
}

/****** fore init ******/

func initrun() {
	//Generate folders
	gendir("Project")
	gendir("Components")
	gendir("Fore")

	//Generate Subfolders
	gendir("Project/html")
	gendir("Project/css")
	gendir("Project/js")
	gendir("Project/assets")

	//Generate strings.json for global 
	//and reusable variable strings
	genfile("strings.json", "{\n\t\"title\": \"My Project\"\n}")
	genfile("package.json", "")
	genfile("config.json", "")
	genfile(".gitignore", "/Components\n/Fore\n*.json")

	fmt.Println("Fore Project Initialized.")
}

func watch() {
	w := watcher.New()
	w.SetMaxEvents(1)

	fmt.Println("Fore Watcher Started")
	
	go func() {
		for {
			select {
			case event := <-w.Event:	
				//Kasi nagcoconflict yung compile(), nababasa as new event
				//So bawal na imanual update yung html.
				//To change html files, manual compile needed
				if filetypefunc(event.Path) != "html" {
					compile()
					fmt.Println(event.Path + " saved. -", time.Now().Format(time.RFC850))
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				fmt.Println("Fore Watcher Stopped")
			}
		}
	}()
	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	go func() {
		w.Wait()
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

/****** Filename and Filetype ******/

func filenamefunc(filename string) string {
	output := ""
	split := strings.Split(filename, ".")

	for a := 0; a < len(split)-1; a++ {
		if output == "" {
			output = output + split[a] 
		} else {
			output = output + "." + split[a] 
		}
	}

	return output
}

func filetypefunc(filename string) string {
	split := strings.Split(filename, ".")

	return split[len(split)-1]
}

/****** File Manipulation ******/

func gettabs(line string) string {
	//Gets the tabs to add before a string
	match, _ := regexp.Compile("[\\s\t]*")
    return match.FindString(line)
}

func checktitle(line string) string {
	//Check if file line is a Title Component
	match, _ := regexp.Compile("<title[\\s]*text=\"([A-Za-z0-9\\s]*)\"[\\s]*/>")
    
	if(len(match.FindStringSubmatch(line)) > 0) {
		return match.FindStringSubmatch(line)[1]
	} else {
		return match.FindString(line)
	}
}

func checkcomponent(line string) string {
	//Check if file line has a Custom Component
	match, _ := regexp.Compile("<([A-Z][A-Za-z0-9\\-\\_]*)[\\s]*/>")

	if(len(match.FindStringSubmatch(line)) > 0) {
		return match.FindStringSubmatch(line)[1]
	} else {
		return match.FindString(line)
	}
}

func checkimports(line string) string {
	//Check if file line is a start or finish of an import statement
	match, _ := regexp.Compile("<[/]*imports>")
    return match.FindString(line)
}

func checkstyles(line string) string {
	//Check if file line is a start or finish of styles
	match, _ := regexp.Compile("<[/]*style>")
    return match.FindString(line)
}

func checkscripts(line string) string {
	//Check if file line is a start or finish of an import statement
	match, _ := regexp.Compile("<[/]*script>")
    return match.FindString(line)
}

func checkcontent(line string) string {
	//Check if file line is a start or finish of styles
	match, _ := regexp.Compile("<[/]*content>")
    return match.FindString(line)
}

func checkvariables(line string) []string {
	//Check if file line has a variable
	match, _ := regexp.Compile("@([A-Za-z0-9\\-\\_]*)")

	if(len(match.FindAllString(line, -1)) > 0) {
		return match.FindAllString(line, -1)
	} else {
		return make([]string, 0)
	}
}

/****** Get JSON Strings ******/

func getjson(filename string) string {
	dat, _ := ioutil.ReadFile(filename)

	r, _ := regexp.Compile("[\n\t]")
	return r.ReplaceAllString(string(dat), "")
}

/****** Get Current Working Directory ******/

func pwd() string{
	//Gets current working directory
	dir, _ := os.Getwd()
	
	//returns working directory
	return dir
}

/****** Generate Directory ******/

func gendir(dirname string) {
	//choose your permissions well
	pathErr := os.Mkdir(dirname, 0777)
 
	//check if you need to panic, fallback or report
	if pathErr != nil {
		fmt.Println(pathErr)
	}
}

/****** Generate File ******/

func genfile(filename string, content string) {
	file, _ := os.Create(filename)
    defer file.Close()
	fmt.Fprintf(file, content)
}

/****** Read Fore File ******/

func readfile(pwd string, filename string) (string, string, string, string, string, string){
	title, imports, styles, scripts, content := "", "", "", "", ""	
	importcheck, stylecheck, scriptcheck, contentcheck := 0, 0, 0, 0

	genhtmlfile := "Project/html/" + filenamefunc(filename) + ".html"

	dat, _ := ioutil.ReadFile(pwd + "/Fore/" + filename)

	splitfile := strings.Split(string(dat), "\n")

	stringsfile := getjson("strings.json")
	
	for _, line := range splitfile {
		vars := checkvariables(line)
	
		for _, tempvars := range vars {
			//This removes the @ symbol of the variable
			newvars := strings.Replace(tempvars, "@", "", -1)

			// Checks if variable name (strings) has an equivalent
			// Key in the JSON File (stringsfile)
			valueGlobal := gjson.Get(stringsfile, newvars)
			valueLocal := gjson.Get(stringsfile, filenamefunc(filename) + "." + newvars)

			//So if value is not null (meaning it has an equivalent)
			if valueLocal.Str != "" {
				//It switches the variable (strings) to the key value (value)
				line = strings.Replace(line, tempvars, valueLocal.Str, -1)
			} else if valueGlobal.Str != "" {
				line = strings.Replace(line, tempvars, valueGlobal.Str, -1)
			}
		}

		if importcheck == 0 && stylecheck == 0 && scriptcheck == 0 && contentcheck == 0 {
			switch {
				case checkimports(line) != "":
					importcheck = 1
				case checkstyles(line) != "":
					stylecheck = 1
				case checkscripts(line) != "":
					scriptcheck = 1
				case checkcontent(line) != "":
					contentcheck = 1
				case checktitle(line) != "":
					title = checktitle(line)
			}
		} else if importcheck == 1 {
			switch {
				case checkimports(line) != "":
					importcheck = 0
				default: 
					imports = imports + line + "\n"
			}
		} else if stylecheck == 1 {
			switch {
				case checkstyles(line) != "":
					stylecheck = 0
				default: 
					styles = styles + line + "\n"
			}
		} else if scriptcheck == 1 {
			switch {
				case checkscripts(line) != "":
					scriptcheck = 0
				default: 
					scripts = scripts + line + "\n"
			}
		} else if contentcheck == 1 {
			switch {
				case checkcontent(line) != "":
					contentcheck = 0
				case checkcomponent(line) != "":
					dat, err := ioutil.ReadFile(pwd + "/Components/" + checkcomponent(line) + ".fore")
					if err != nil {
						content = content + gettabs(line) + "<ComponentNotFound />"
					}
		
					splitcomp := strings.Split(string(dat), "\n")
		
					for _, compline := range splitcomp {
						content = content + gettabs(line) + compline + "\n"
					}
				default: 
					content = content + line + "\n"
			}
		}
	}

	return title, imports, styles, scripts, content, genhtmlfile
}

func htmlcompile(title string, imports string, styles string, scripts string, content string, genhtmlfile string) {
	i, s, sc, c := "", "", "", ""

	importsplit := strings.Split(imports, "\n")
	for _, line := range importsplit {
		if i == "" {
			i = i + "\t" + line
		} else {
			i = i + "\n\t" + line
		}
	}

	stylessplit := strings.Split(styles, "\n")
	for _, line := range stylessplit {
		if s == "" {
			s = s + "\t" + line
		} else {
			s = s + "\n\t\t" + line
		}
	}

	scriptsplit := strings.Split(scripts, "\n")
	for _, line := range scriptsplit {
		if sc == "" {
			sc = sc + "\t" + line
		} else {
			sc = sc + "\n\t\t" + line
		}
	}

	contentsplit := strings.Split(content, "\n")
	for _, line := range contentsplit {
		if c == "" {
			c = c + "\t" + line
		} else {
			c = c + "\n\t" + line
		}
	}

	html := "<html>\n\t<head>\n"
	html = html + "\t\t<title>" + title + "</title>\n"

	if imports != "" { html = html + i }
	if styles != "" { html = html + "\t\t<style>\n\t" + s + "</style>\n" }

	html = html + "\t</head>\n\t<body>\n"

	if scripts != "" { html = html + "\t\t<script>\n\t" + sc + "</script>\n"}
	if content != "" { html = html + c }

	html = html + "</body>\n</html>"
	
	genfile(genhtmlfile, html)
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
			case "init":
				initrun()
			case "compile", "comp":
				compile()
			case "watch", "watcher":
				watch()
			case "version":
				fmt.Println("Fore v1.0.0")
			default: 
				fmt.Println("Fore: Command not found.")
		}
	} else {
		home()
	}
}