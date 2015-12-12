package gib
import "time"
import "io"
import "errors"

//A full system requires one Datastore, one Imagestore and some number of Frontends for each board
//each layer should be thread safe internally but make no assumptions about the ordering or timing its methods will be called from

//can a thread just be a slice of posts?
type Post struct{
	Timestamp time.Time
	Name string
	Email string
	Subject string
	//TODO: when do we parse greentext, spoilers, etc?
	//do we do it on the fly or store the htmlified version in Datastore?
	Comment string
	//these include file extensions
	Files []string
	Bump bool
	Sticky bool

	Number int

}

//Each board can have its own datastore object.  
//This makes namespaces simpler, i.e. we can identify uniquely a post by its number and Datastore.
//
//Ideally this should be able to handle loosely structured (flat SQL database) and highly structured (nosql or custom in-memory data structure) storage, but
//this is always a compromise.  If we provide an interface for traversing the board tree logically that means more heavy lifting in the SQL backend, while
//using purely random lookups (getPostById style) means the middle layer will have to reconstruct information that was already present in the backend.
//TODO: should this interface include ordering (bump, sort or something like that) functions?  some backends will need to be kept in order while others
//will sort their results on request.  We could simply rely on the backend to sort itself if necessary either periodically or when posts are added or deleted.
type Datastore interface{
	//these need not take effect immediately if it is more efficient to coalesce them
	AddThread(op Post) error
	AddReply(post Post, Op int) error
	//deletes the thread if it's an OP
	DeletePost(num int) error

	//queries
	GetPost(num int) (Post, error)
	GetThread(num int) ([]Post, error)

	//OPs sorted by last bump
	//note: careless use of skip will lead to race conditions (missing/duplicate threads)
	GetThreads(skip int, maxcount int) ([]Post, error)

}
//gib should write an image and verify its contents before adding the reply or thread
//this will probably just be a folder containing images
//imagestore does not need to be aware of thumbnails, they will just be images with different names
type Imagestore interface{
	WriteImage(name string, data io.Reader) error
	ReadImage(name string) (io.Reader, error)
	DeleteImage(name string) error
}

//Frontends will receive a *Board which they will use to retreive all their data
//this interface can be expanded as features are added.  
//anything that would be useful for multiple frontends is fair game for moving into gib
type Board struct{
	Datastore
	Imagestore
	BoardName string
	lastpost int
	//board settings (allowed file types, supported functionality, etc)
	//...
}
//these 3 are just passthroughs to Datastore.  Is there a better way?
//they could do parsing of comments if we decide not to store html in the Datastore
func (b *Board)FetchPost(num int) (Post, error){return Post{}, errors.New("not imp")}
func (b *Board)FetchThread(num int) ([]Post, error){
	return b.Datastore.GetThread(num)
}
//FIXME:renamed to avoid accidental infinite recursion, but should have a better name
func (b *Board)FetchThreads(skip int, maxcount int) ([]Post, error){
	return b.GetThreads(0, 100)
}
//GetImage and GetThumb will both take the filename in Post, gib must convert that name to a thumbnail name for GetThumb
//my proposal: "thumb_"+name
func (b *Board)FetchImage(name string) (io.Reader, error){return nil, errors.New("not imp")}
func (b *Board)FetchThumb(name string) (io.Reader, error){return nil, errors.New("not imp")}
//these must verify the files, save them, generate thumbs (unless this is on demand) and add the post to the database
//post ordering can be updated here or asynchronously (or on request in the case of rdbms backend)
func (b *Board)PostThread(op Post, files []io.Reader) error{
	op.Timestamp = time.Now()
	op.Number = b.lastpost
	b.lastpost ++//TODO:atomic
	return b.AddThread(op)
}
func (b *Board)PostReply(post Post, files []io.Reader, Op int) error{
	post.Timestamp = time.Now()
	post.Number = b.lastpost
	b.lastpost ++//TODO:atomic
	return b.AddReply(post, Op)
}



