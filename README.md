# Code Reviewer.

![License: MIT](https://img.shields.io/badge/Language-Golang-blue.svg)


### Requirements:
* Gemini API Key
* Golang
* Add .env file with values to specified keys:
    ```txt
    export GEMINI_API_KEY="your-api-key-goes-here"
    ```
### Run code locally
 ```
 $ git clone https://github.com/edwinnduti/code-reviewer.git 
 $ cd code-reviewer
 $ go mod download
 $ go build main.go -o code-reviewer
 $ ./code-reviewer -file=./mycode.py
 ```

### Build on Windows
```go
GOOS=windows GOARCH=amd64 go build -o code-reviewer.exe main.go
```

### Build on Linux
```go
GOOS=linux GOARCH=amd64 go build -o code-reviewer main.go
```

Have a day full of ❤️