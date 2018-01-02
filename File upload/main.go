package main

import "net/http"
import "log"
import "io/ioutil"
import "fmt"

func main() {
	http.HandleFunc("/", fileUpload)
	http.ListenAndServe(":8000", nil)
}

func fileUpload(w http.ResponseWriter, req *http.Request) {
	var s string

	if req.Method == http.MethodPost {
		f, _, err := req.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading the file", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading the file", http.StatusInternalServerError)
			return
		}
		s = string(bs)
	}
	w.Header().Set("CONTENT-TYPE", "text/html; enc=UTF-8")
	fmt.Fprintf(w, `<form action "/" method="post" enctype="multipart/form-data">
                   Choose file to upload:<br>
                   <input type="file" name="file"><br>
                   <input type="submit"> 
					</form>
					<h2>%v</h2>`, s)

}