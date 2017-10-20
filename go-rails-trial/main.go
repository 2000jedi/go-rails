package main

import (
    "fmt"
    "go-rails/go-rails-trial/template"
)

func main() {
    fmt.Println(Tables["user"].All())
    fmt.Println(template.Construct_example(map[string]interface{}{
        "A": "",
        "B": "2",
        "content": "v",
        "xs":[]map[string]interface{}{{
            "id": map[string]interface{}{
                "first": "1", "second": "2"},
                "prim": map[string]interface{}{
                    "first":"prim",
                },
            },
        },
    }))
}
