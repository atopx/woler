# woler
Wake on LAN Library


## Install
```bash
go get github.com/yanmengfei/woler
```

## Use
```go
package main

import (
    "fmt"
    "github.com/yanmengfei/woler"
)


func main() {
    var macaddr := "F0-2F-74-B0-1D-E0"
    if err := woler.Do(macaddr); err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("Wake up device successfully")
    }
}
```
