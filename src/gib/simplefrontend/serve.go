package simplefrontend
//a simple proof of concept front end
//access via localhost:8000/catalog, localhost:8000/postform

import "gib"
import "net/http"
import "fmt"
import "strconv"

//obviously a real frontend would handle multiple boards so it would use Handler instead of
//HandleFunc and globals
var board *gib.Board

func cataloghandler(w http.ResponseWriter, req *http.Request){
	threads, err := board.FetchThreads(0, 100)
	for _, thread := range threads{
		fmt.Fprintln(w, thread)
		replies, _ := board.FetchThread(thread.Number)
		for _, post := range replies[1:]{
			fmt.Fprintln(w, "reply: ", post)
		}
	}
	fmt.Fprintln(w, err)
}
func posthandler(w http.ResponseWriter, req *http.Request){
	var p gib.Post
	p.Comment = req.FormValue("comment")
	replyto, _ := strconv.Atoi(req.FormValue("replyto"))
	if replyto < 0{
		board.PostThread(p, nil)
	}else{
		board.PostReply(p, nil, replyto)
	}
	http.Redirect(w, req, "/catalog", http.StatusSeeOther)
}
func postformhandler(w http.ResponseWriter, req *http.Request){
	w.Write([]byte(`<html><body><form action="/post" method="post">
<textarea name="comment"/></textarea>
<input type="number" name="replyto"/>
<input type="submit" value="post"/>
</form></body></html>`))
}


func Startserver(b *gib.Board, addr string){
	board = b
	http.HandleFunc("/catalog", cataloghandler)
	http.HandleFunc("/post", posthandler)
	http.HandleFunc("/postform", postformhandler)
	http.ListenAndServe(addr, nil)
}
