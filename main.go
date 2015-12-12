package main
import "gib/simpledb"
import "gib/simplefrontend"
import "gib"
//import "fmt"


func main(){
	db := simpledb.BoardStore{}

	board := gib.Board{Datastore:&db, BoardName: "test board"}


	//the following stuff would happen in the frontend in response to http requests etc.
/*	fmt.Println("add thread")
	fmt.Println(board.PostThread(gib.Post{Comment:"test post"}, nil))
	fmt.Println("add thread")
	fmt.Println(board.PostThread(gib.Post{Comment:"test post 2"}, nil))
	fmt.Println("add reply")
	fmt.Println(board.PostReply(gib.Post{Comment:"reply to test post"}, nil, 0))


	fmt.Println("show threads")
	fmt.Println(board.FetchRangeThreads(0,100))
	fmt.Println("show thread 0")
	fmt.Println(board.FetchThread(0))
*/
	simplefrontend.Startserver(&board, ":8000")


}
