# Fore Template Engine

Fore Template Engine is designed to remove code repetition and planned to make revisions faster and easier by using Fore Components and Global or Local Variables.

## Sample Fore Code
```
<title text="@title" /> //Title for the page

<imports> // CSS, Javascript Imports
    <link rel="stylesheet" href="index.css" />
    <script src="./ezss/js/ezss.js"></script>
</imports>

<styles> // Use can also use styles
	.body {
		width: 100%;
		margin: 0;
		padding: 0;
	}
</styles>

<content>
    <div id="main-container">
        <Header /> //Custom Component Import
        <About /> //Custom Component Import
        <div class="nav-right">
            <a href="#home" class="nav-tab">@home</a>
            <a href="#about-me" class="nav-tab">@about-me</a>
            <a href="#skills" class="nav-tab">@skills</a>
            <a href="#projects" class="nav-tab">@projects</a>
            <a href="#contact" class="nav-tab">@contact</a>
        </div>
    </div>
</content>
```

## To start project, run command on your local repository:
```
fore init
```

## File Structure
- Project
    - html
    - css
    - js
    - assets
- Components
- Fore
- .gitignore
- strings.json
- package.json
- config.json

## Usage

### Commands

**IMPORTANT: Commands only run on root folder!**

```
fore init //Initializes Fore Project
fore compile //Compile Fore files
fore watch //Runs Fore Watcher to auto hot-reload
```

### Files

#### /Components/Header.fore
Header.fore is **NOT** initialized during `fore init`. Fore Components are created manually and **SHOULD** be stored in `/Components` folder. Storing Components to seperated folders will not work and would create a `<ComponentNotFound />` line on file Output.

Sample Header.fore
```
<nav class="ezss-nav nav-clear nav">
    <div class="nav-center">
        <a href="#home" class="nav-tab">HOME</a>
        <a href="#about-me" class="nav-tab">ABOUT ME</a>
        <a href="#skills" class="nav-tab">SKILLS</a>
        <a href="#projects" class="nav-tab">MY PROJECTS</a>
        <a href="#contact" class="nav-tab">CONTACT ME</a>
    </div>
</nav>
<nav class="ezss-nav nav-clear nav-right sidebar">
    <a href="#home" onclick="openNav()" class="nav-tab">&#9776;</a>
</nav>
<div id="mySidenav" class="sidenav">
    <a href="javascript:void(0)" class="closebtn" onclick="closeNav()">&times;</a>
    <a href="#home" class="nav-tab">HOME</a>
    <a href="#about-me" class="nav-tab">ABOUT ME</a>
    <a href="#skills" class="nav-tab">SKILLS</a>
    <a href="#projects" class="nav-tab">MY PROJECTS</a>
    <a href="#contact" class="nav-tab">CONTACT ME</a>
</div>
```

#### /Fore/index.fore
index.fore is **NOT** initialized during `fore init`.
Fore files stored in `/Fore` directory are the structure files that will be used in generating `.html` files.

#### /Project
Project folder will be initialized during `fore init`.
This folder will be used to store fore compiled files and once finished, the only folder needed to be deployed. Using Fore Watcher, saving files will automatically update the `/Project/html` folder files. Fore Watcher does **NOT** allow manual updates of **HTML** files. To compile manually update HTML Files, run `fore compile` in the command line.

#### strings.json
Strings.json will be used to store Global and Local Variables that can be used in Fore files.

To Create Local Variables, create an object using the filename as the key. Creating different key name would not work.

Sample strings.json Code
```
{
    "title": "My Project", //Global Variables
    "home": "HOME",
    "about-me": "ABOUT ME",
    "skills": "SKILLS",
    "project": "MY PROJECTS",
    "contact": "CONTACT US",
    "index": { //Local Variable Group
        "title": "Home Page",
        "con-one": "Lorem Ipsum",
        "con-two": "Lorem Ipsum",
        "con-three": "Lorem Ipsum",
        "css": {
            "nav-class": "ezss-btn nav-right btn-clear"
        }
    },
    "about-us": {
        "title": "About Page",
        "con-one": "Lorem Ipsum",
        "con-two": "Lorem Ipsum",
        "con-three": "Lorem Ipsum"
    }
}
```

 To call variables, simply add `@key-name` to your Fore file. In calling local variables, you won't need to call the filename group. In this case, calling `@title` in `index.fore` will return the value of `Home Page`. For accessing inner groups such as the `css` object, calling `css.nav-class` would work fine. If a variable is not found in a Local Variable Group, Fore will automatically get the values from Global Variables.

 ## Todo
 1. Add Loops.
 2. Make Syntax more shorter!
 3. Add more features.