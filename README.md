<div style="float: right; margin-left: 20px;">
  <img src="docs/cycle-detect.png" alt="Cycle Detect Image" width="300" height="300">
</div>

# Cycle-detect

English | [Türkçe](README-tr_tr.md) 
___
## Motivation
Hi there! I'm Eray, a software developer who is new to Golang. I was fascinated by the language, but I often encountered import cycle errors which became a problem for me over time. Especially after writing code for a long time, tracking these errors in a project made it difficult to track the import statements used.

Based on this problem I experienced, I thought: Why not have a project that shows directly in which file the import cycles are? With this idea, I created this project to prevent people from facing similar problems and to offer a more efficient development process. This way, you can quickly identify and resolve import cycle errors, making the development process smoother.
___
## Getting Started
Cycle detect is a process that aims to detect cyclic dependencies in import operations within the source code of a program. In other words, it is an analysis process that determines the situation where a module or package directly or indirectly imports another module and this operation creates a loop. Cycle detect guides software developers to resolve such dependencies by identifying the problem in advance in terms of code organization and maintenance.

## Installation

###  `go install`
Use Go 1.20:
___

```bash
git clone https://github.com/eray-can/cycle-detect.git
```


___
## Usage
The usage is very simple. All you need to do is download the project, define your file path and run it. You can follow the example below.

```
//You just need to give the file path of the related golang project to the 'project_path' variable from the detect.toml file
```
___

![img.png](docs%2Fimg.png)
[Detect.toml](./detect.toml) You can access the detect.toml file here

___

```go
func main() {

    engine := runner.NewEngine()
    engine.Run()
    defer engine.Close()

}
```
You can access the [Main:](./main.go)  file here
___
## Contact
If you have any questions, suggestions or feedback, please do not hesitate to contact me.

- E-posta: ceray6575@gmail.com
- Twitter: [@Eraynac13](https://twitter.com/Eraynac13)
- LinkedIn: [eraycan](https://www.linkedin.com/in/eraycan/)


You can also reach me via GitHub. You can report bugs or contribute to the project by opening an issue.
