package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/tidwall/gjson"
)

/*
 * Fore Functions
 * Fore Template Engine is designed to remove code repetition
 * and planned to make revisions faster and easier by using Fore
 * Components and Global or Local Variables.
 */

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
			title, imports, styles, scripts, content, genhtmlfile := forefile(pwd(), f.Name())

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
				/*
				* Kasi nagcoconflict yung compile(), nababasa as new event
				* So bawal na imanual update yung html.
				* To change html files, manual compile needed
				 */

				if filetypefunc(event.Path) != "html" {
					compile()
					fmt.Println(event.Path+" saved. -", time.Now().Format(time.RFC850))
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

	if len(match.FindStringSubmatch(line)) > 0 {
		return match.FindStringSubmatch(line)[1]
	}

	return match.FindString(line)
}

func checkcomponent(line string) string {
	//Check if file line has a Custom Component
	match, _ := regexp.Compile("<([A-Z][A-Za-z0-9\\-\\_]*)[\\s]*/>")

	if len(match.FindStringSubmatch(line)) > 0 {
		return match.FindStringSubmatch(line)[1]
	}

	return match.FindString(line)
}

func checkcollection(line string) string {
	//Check if file line is a start or finish of a collection
	match, _ := regexp.Compile("<[/]*(content|logic|style|imports)[\\s/]*>")

	if len(match.FindStringSubmatch(line)) > 0 {
		return match.FindStringSubmatch(line)[1]
	}

	return match.FindString(line)

}

func checkendtag(line string) string {
	match, _ := regexp.Compile("</(content|logic|style|imports)>")
	return match.FindString(line)
}

func checkvariables(line string) []string {
	//Check if file line has a variable
	match, _ := regexp.Compile("@([A-Za-z0-9\\-\\_]*)")

	if len(match.FindAllString(line, -1)) > 0 {
		return match.FindAllString(line, -1)
	}

	return make([]string, 0)
}

func escape(line string) string {
	line = strings.Replace(line, "$", "&#36;", -1)
	line = strings.Replace(line, "%", "&#37;", -1)
	line = strings.Replace(line, "^", "&#94;", -1)
	line = strings.Replace(line, "`", "&#96;", -1)
	line = strings.Replace(line, "|", "&#124;", -1)
	line = strings.Replace(line, "~", "&#126;", -1)

	return line
}

/****** Get JSON Strings ******/

func getjson(filename string) string {
	dat, _ := ioutil.ReadFile(filename)

	r, _ := regexp.Compile("[\n\t]")
	return r.ReplaceAllString(string(dat), "")
}

/****** Get Variable Values ******/

func getvariables(line string, keyname string) string {
	stringsfile := getjson("strings.json")
	vars := checkvariables(line)

	for _, tempvars := range vars {
		//This removes the @ symbol of the variable
		newvars := strings.Replace(tempvars, "@", "", -1)

		// Checks if variable name (strings) has an equivalent
		// Key in the JSON File (stringsfile)
		valueGlobal := gjson.Get(stringsfile, newvars)
		valueLocal := gjson.Get(stringsfile, keyname+"."+newvars)

		//So if value is not null (meaning it has an equivalent)
		if valueLocal.Str != "" {
			//It switches the variable (strings) to the key value (value)
			line = strings.Replace(line, tempvars, valueLocal.Str, -1)
		} else if valueGlobal.Str != "" {
			line = strings.Replace(line, tempvars, valueGlobal.Str, -1)
		}
	}

	return escape(line)
}

/****** Get Current Working Directory ******/

func pwd() string {
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

/****** Read File ******/

func readfile(file string) string {
	dat, _ := ioutil.ReadFile(file)
	return string(dat)
}

/****** Read Fore File ******/

func forefile(pwd string, filename string) (string, string, string, string, string, string) {
	title, imports, styles, scripts, content := "", "", "", "", ""
	importcheck, stylecheck, scriptcheck, contentcheck := 0, 0, 0, 0

	genhtmlfile := "Project/html/" + filenamefunc(filename) + ".html"

	filecontents := readfile(pwd + "/Fore/" + filename)

	splitfile := strings.Split(filecontents, "\n")

	for _, line := range splitfile {
		varline := getvariables(line, filenamefunc(filename))
		line = getvariables(varline, filenamefunc(filename))

		// This runs on start when all checkers still have a default value (0),
		// Checks if a string is a starting tag of a collection, else does nothing
		if importcheck == 0 && stylecheck == 0 && scriptcheck == 0 && contentcheck == 0 {
			switch checkcollection(line) {
			case "imports":
				importcheck = 1
			case "style":
				stylecheck = 1
			case "logic":
				scriptcheck = 1
			case "content":
				contentcheck = 1
			default:
				if checktitle(line) != "" {
					title = checktitle(line)
				}
			}
		} else {
			/*
			 * So if one checker has a value of 1, this means it has
			 * already found a starting tag for the collection.
			 * So first we need to check if the line has a starting
			 * tag or ending tag. If the line is an ending tag, return the
			 * checker value to 0, else it is a starting tag, changing
			 * the other checkers to 0 and changing the starting tag checker to 1
			 * If the line doesn't have a starting or ending tag, it
			 * checks the line if not empty and adds it to the string it is connected to.
			 */

			if checkcollection(line) != "" {
				switch checkcollection(line) {
				case "imports":
					if checkendtag(line) != "" {
						importcheck = 0
					} else {
						importcheck, stylecheck, scriptcheck, contentcheck = 1, 0, 0, 0
					}
				case "style":
					if checkendtag(line) != "" {
						stylecheck = 0
					} else {
						importcheck, stylecheck, scriptcheck, contentcheck = 0, 1, 0, 0
					}
				case "logic":
					if checkendtag(line) != "" {
						scriptcheck = 0
					} else {
						importcheck, stylecheck, scriptcheck, contentcheck = 0, 0, 1, 0
					}
				case "content":
					if checkendtag(line) != "" {
						contentcheck = 0
					} else {
						importcheck, stylecheck, scriptcheck, contentcheck = 0, 0, 0, 1
					}
				default:
					if checktitle(line) != "" {
						title = checktitle(line)
					}
				}
			} else if importcheck == 1 {
				if line != "" {
					imports = imports + line + "\n"
				}
			} else if stylecheck == 1 {
				if line != "" {
					styles = styles + line + "\n"
				}
			} else if scriptcheck == 1 {
				if line != "" {
					scripts = scripts + line + "\n"
				}
			} else if contentcheck == 1 {
				switch {
				case checkcomponent(line) != "":
					dat, err := ioutil.ReadFile(pwd + "/Components/" + checkcomponent(line) + ".fore")
					if err != nil {
						content = content + gettabs(line) + "<ComponentNotFound />"
					}

					splitcomp := strings.Split(string(dat), "\n")

					for _, compline := range splitcomp {
						varcompline := getvariables(compline, checkcomponent(line))
						compline = getvariables(varcompline, checkcomponent(line))

						content = content + gettabs(line) + compline + "\n"
					}
				case line != "":
					content = content + line + "\n"
				}
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

	if imports != "" {
		html = html + i
	}
	if styles != "" {
		html = html + "\t\t<style>\n\t" + s + "</style>\n"
	}

	html = html + "\t</head>\n\t<body>\n"

	if scripts != "" {
		html = html + "\t\t<script>\n\t" + sc + "</script>\n"
	}
	if content != "" {
		html = html + c
	}

	html = html + "</body>\n</html>"

	genfile(genhtmlfile, html)
}

/* End of Fore Functions */

/*
 * Fore EzSS Functions
 * EzSS is a CLI tool for creating and compressing CSS Files.
 * EzSS allows auto-generation of CSS Files based on HTML Classes
 * and Ids using CSS Best Practices. EzSS also allows CSS Compression.
 * EzSS Functions also uses some of Fore's Functions like
 * genfile(), pwd(), and escape().
 */

/* Creates and returns comments based on Input type and names */

func makecomments(category string, name string) string {
	comment := ""

	switch category {
	case "info":
		comment += "/*\n"
		comment += "\t!!!DON'T FORGET TO UPDATE THESE INFORMATION!!!\n"
		comment += "\tFilename: index.html\n"
		comment += "\tAuthor: Janjan Medina\n"
		comment += "\tAuthor URI: https://github.com/medinajuanantonio95\n"
		comment += "\tProject: " + pwd() + "\n"
		comment += "*/\n\n"
		break
	case "general":
		comment = "/***** General Styles *****/\n\n"
		break
	case "header":
		comment = "/***** Header Style *****/\n\n"
		break
	case "footer":
		comment = "/***** Footer Style *****/\n\n"
		break
	case "nav":
		comment = "/***** Navigation Style *****/\n\n"
		break
	case "basic":
		comment = "/***** " + name + " *****/\n\n"
		break
	default:
		comment = "/* Comment type unknown */\n\n"
	}

	return comment
}

/* Add Media Queries */

func mediaqueries() string {
	mediaqueries := ""

	mediaqueries += "/* Small devices (tablets, 768px and up) */\n\n"
	mediaqueries += "@media (min-width: 768px) {\n\n}\n\n"
	mediaqueries += "/* Small devices (desktops, 992px and up) */\n\n"
	mediaqueries += "@media (min-width: 992px) {\n\n}\n\n"
	mediaqueries += "/* Small devices (large desktops, 1200px and up) */\n\n"
	mediaqueries += "@media (min-width: 1200px) {\n\n}\n\n"

	return mediaqueries
}

/* Add CSS for Body */

func bodycss() string {
	return "body {\n\tmargin: 0;\n\tpadding: 0;\n\twidth: 100%;\n}\n\n"
}

/* Add CSS Function */

func cssfunc(name string, category string) string {
	classifier := ""

	if category == "id" {
		classifier = "#"
	} else {
		classifier = "."
	}

	defString := " {\n\n}\n\n"

	output := classifier + name + defString

	return output
}

/* Add Import to HTML File */

func addimport(filename string) string {
	return "<link rel=\"stylesheet\" href=\"" + filename + "\" />"
}

/* Getting Classes for one line with Regex */

func getclasses(line string) string {
	//Check if file line has a class
	match, _ := regexp.Compile("class=\"([a-zA-Z0-9\\s-\\_]*)\"")

	if len(match.FindStringSubmatch(line)) > 0 {
		return match.FindStringSubmatch(line)[1]
	}

	return match.FindString(line)
}

/* Getting Ids for one line with Regex */

func getids(line string) string {
	//Check if file line has a class
	match, _ := regexp.Compile("id=\"([a-zA-Z0-9\\s-]*)\"")

	if len(match.FindStringSubmatch(line)) > 0 {
		return match.FindStringSubmatch(line)[1]
	}

	return match.FindString(line)
}

/* Get CSS Ids and Classes */

func csscontent(filename string) ([]string, []string) {
	classes := make([]string, 1) //Define Slices to put Classes and Ids
	ids := make([]string, 1)

	filecontents := readfile(pwd() + "\\" + filename) //Read HTML File
	splitfile := strings.Split(filecontents, "\n")    //Split HTML File to lines

	for _, line := range splitfile {
		lclasses := getclasses(line) //Get Classes and IDs of a line
		lids := getids(line)

		splitclasses := strings.Split(lclasses, " ")
		splitids := strings.Split(lids, " ")

		/*
		 * All the results got from the classes and ids
		 * function, it will get split by "\n" and will get pushed
		 * to the overall slices (classes and ids)
		 * and what will be return as output
		 */

		for _, class := range splitclasses {
			if class != "" {
				classes = append(classes, escape(class))
			}
		}

		for _, id := range splitids {
			if id != "" {
				ids = append(ids, escape(id))
			}
		}
	}

	return classes, ids
}

func makecontent(classes []string, ids []string) string {
	compiled := ""

	compiled += makecomments("info", "")
	compiled += makecomments("general", "") + bodycss()

	for a := 0; a < len(ids); a++ {
		if escape(ids[a]) != "" {
			switch escape(ids[a]) {
			case "header":
				compiled += makecomments("header", escape(ids[a]))
				break
			case "footer":
				compiled += makecomments("footer", escape(ids[a]))
				break
			case "nav":
				compiled += makecomments("nav", escape(ids[a]))
				break
			default:
				compiled += makecomments("basic", escape(ids[a]))
			}

			compiled += cssfunc(escape(ids[a]), "id")
		}
	}

	for b := 0; b < len(classes); b++ {
		if escape(classes[b]) != "" {
			compiled += cssfunc(escape(classes[b]), "classes")
		}
	}

	compiled += mediaqueries()

	return compiled
}

func makecss(filename string) {
	a, b := csscontent(filename)
	content := makecontent(a, b)

	genfile(filenamefunc(filename)+".css", content)
}

/* End of EzSS Functions */

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
			home()
		case "ezss": //EzSS Integration
			switch os.Args[2] {
			case "create":
				makecss(os.Args[3])
			case "generate":
				makecss(os.Args[3])
			case "gen":
				makecss(os.Args[3])
			case "compress":
				fmt.Println("EzSS Compress still not available.")
			case "comp":
				fmt.Println("EzSS Compress still not available.")
			case "read":
				fmt.Println(readfile(os.Args[3]))
			default:
				fmt.Println("Fore EzSS: Command not found.")
			}
		default:
			fmt.Println("Fore: Command not found.")
		}
	} else {
		home()
	}
}
