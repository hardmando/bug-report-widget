package main
import (
    "fmt"
    "log"
    "net/http"
)
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Bug service running")
    })
    log.Println("Bug service starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
