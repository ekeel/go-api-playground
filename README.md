# go API Playground

*A playground for testing go API features.*  

The current API allows for creating routes based on "route plugins". These "route plugins" are made up of a JSON descriptor and a code file. The coding languages that can currently be used are listed below.  

All plugins are currently required to return a JSON string. It is possible that other configuration languages will be implemented in the future.  

- [go API Playground](#go-api-playground)
  - [Plugin Languages](#plugin-languages)
  - [Plugin Examples](#plugin-examples)
    - [Python](#python)
      - [Descriptor JSON (python_example_route.json)](#descriptor-json-python_example_routejson)
      - [Code File (python_example_route.py)](#code-file-python_example_routepy)
    - [Go](#go)
      - [Descriptor JSON (go_example_route.json)](#descriptor-json-go_example_routejson)
      - [Code File (go_example_route.go)](#code-file-go_example_routego)
    - [PowerShell](#powershell)
      - [Descriptor JSON (pwsh_example_route.json)](#descriptor-json-pwsh_example_routejson)
      - [Code File (pwsh_example_route.ps1)](#code-file-pwsh_example_routeps1)
  - [Tasks](#tasks)

----
## Plugin Languages
*The following languages can be used to create plugins. The languages supported can be expanded by adding a route handler for the language to the main.go file.*

| Name | Status | Notes |
| ---- | ------ | ----- |
| Python | MVP | This type of "route plugin" will execute `python /file/path.py`. There is currently no way to execute `python3 /file/path.py`, meaning that if you need to use python3 it must be aliased to python. |
| Go | MVP | This type of "route plugin" will execute `go run /file/path.go`. In the future, JIT compilation of the go files will be possible. |
| PowerShell | MVP | This type of "route plugin" will execute `pwsh -File /file/path.ps1`. Currently only PowerShell core is supported. In the future, the version of PowerShell (standard vs. core) will be selected based on the OS. |

----
## Plugin Examples

### Python

#### Descriptor JSON (python_example_route.json)

>```json
>{
>  "path": "/testpy/{arg}",
>  "method": "GET",
>  "type": "python",
>  "main_file": "python_example_route.py",
>  "parameters": [
>    "arg"
>  ]
>}
>```

#### Code File (python_example_route.py)

>```python
>import sys
>
>print("{\"test\":\"" + sys.argv[1] + "\"}")
>```

<br>

### Go

#### Descriptor JSON (go_example_route.json)

>```json
>{
>  "path": "/test",
>  
>  "method": "GET",
>  "type": "go",
>  "main_file": "go_example_route.go",
>  "parameters": []
>}
>```

#### Code File (go_example_route.go)

>```go
>package main
>
>import (
>	"encoding/json"
>	"fmt"
>)
>
>func main() {
>	ret, _ := json.Marshal([]string{"1111", "2222", "3333", "4444"})
>
>	fmt.Print(string(ret))
>}
>```

<br>

### PowerShell

#### Descriptor JSON (pwsh_example_route.json)

>```json
>{
>  "path": "/testps/{var1}/{var2}",
>  "method": "GET",
>  "type": "powershell",
>  "main_file": "testps.ps1",
>  "parameters": [
>    "var1",
>    "var2"
>  ]
>}
>```

#### Code File (pwsh_example_route.ps1)

>```powershell
>$Arg1 = $args[0]
>$Arg2 = $args[1]
>
>"[{`"arg1`": `"$Arg1`"}, {`"arg2`": `"$Arg2`"}]"
>```

<br>

----
## Tasks

- [X] Route Plugins - Allow for custom routes in the form of JSON/code files.
- [X] Concurrency Safe Logging - Allow for logging in a concurrently safe way.
- [ ] Testing - Add automated testing (currently uses some basic functional pester testing).
- [ ] Enable Okta Authentication - Require auth. via Okta to process calls.
- [ ] Simple DSL - Allow for the creation of a plugin driven API via a simple DSL.
