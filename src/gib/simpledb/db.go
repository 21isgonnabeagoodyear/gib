package simpledb
//this is a simple and not very good implementation of DataStore.
//it's meant to get things started.  better implementations will come later.

import "gib"
import "sync"
import "errors"

/*
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
*/
type BoardStore struct{
	threads [][]gib.Post
	//single lock will cause starvation when lots of things are happening
	lock sync.Mutex
}
func (b *BoardStore) AddThread(op gib.Post) error{
	b.lock.Lock()
	b.threads = append([][]gib.Post{[]gib.Post{op}}, b.threads...)
	b.lock.Unlock()
	return nil
}
func (b *BoardStore) AddReply(post gib.Post, Op int) error{
	b.lock.Lock()
	for index := range b.threads{
		if b.threads[index][0].Number == Op{
			b.threads[index] = append(b.threads[index], post)
			b.lock.Unlock()
			return nil
		}
	}
	b.lock.Unlock()
	return errors.New("thread does not exist")
}
func (b *BoardStore) DeletePost(num int) error{
	return errors.New("not implemented")
}
func (b *BoardStore) GetPost(num int) (gib.Post, error){
	b.lock.Lock()
	for _, thread := range b.threads{
		for _, post := range thread{
			if post.Number == num{
				b.lock.Unlock()
				return post, nil
			}
		}
	}
	b.lock.Unlock()
	return gib.Post{}, errors.New("post does not exist")
}
func (b *BoardStore) GetThread(num int) ([]gib.Post, error){
	b.lock.Lock()
	for _, thread := range b.threads{
		if thread[0].Number == num{
			b.lock.Unlock()
			return thread, nil
		}
	}
	b.lock.Unlock()
	return nil, errors.New("thread does not exist")
}
func (b *BoardStore) GetThreads(skip int, maxcount int) ([]gib.Post, error){
	var rv []gib.Post
	b.lock.Lock()
	for _, thread := range b.threads{
		rv = append(rv, thread[0])
	}
	b.lock.Unlock()
	return rv, nil
}
