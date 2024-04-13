package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/ezeeyahoo/demoBlogServiceInGrpc/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type stringConst string

const (
	network stringConst = "tcp"
	address stringConst = ":7878"
)

type BlogServer struct {
	proto.UnimplementedBlogServicerServer
}

func main() {
	listener, err := net.Listen(string(network), string(address))
	if err != nil {
		log.Fatalf("failed to initialize a listener- %v", err)
	}

	server := grpc.NewServer()

	proto.RegisterBlogServicerServer(server, &BlogServer{})

	go func() {
		log.Printf("starting the server at: %+v", listener.Addr())
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to start server- %v", err)
		}
	}()

	chn := make(chan os.Signal, 1)
	signal.Notify(chn, os.Interrupt)

	<-chn
	log.Println("shut down server event")
	server.Stop()
	log.Println("clearing listener")
	listener.Close()
}

type PostID string
type BlogEntry struct {
	Title   string
	Author  string
	Content string
	PubData string
	Tags    []string
}

type BlogStorage struct{ storage map[string]BlogEntry }

var blogStorage *BlogStorage

func init() {
	blogStorage = NewBlogStorage()
}

func getUUID() string {
	return uuid.NewString()
}

func NewBlogStorage() *BlogStorage {
	return &BlogStorage{storage: make(map[string]BlogEntry)}
}

func blogIncomingTransformer(entry *proto.BlogEntryPost) BlogEntry {
	return BlogEntry{
		Title:   entry.Title,
		Author:  entry.Author,
		Content: entry.Content,
		PubData: entry.PubDate,
		Tags:    entry.Tags,
	}
}

// Create is responsible for creating new blog entry in
func (b *BlogStorage) Create(entry *proto.BlogEntryPost) (*string, *BlogEntry) {
	ID := getUUID()

	blogEntry := blogIncomingTransformer(entry)

	b.storage[ID] = blogEntry

	return &ID, &blogEntry
}

// Update is responsible for updating data in the data storage
func (b *BlogStorage) Update(blogID *string, entry *proto.BlogEntryPost) (*string, *BlogEntry, error) {
	s := b.Delete(blogID)

	var (
		ID   *string
		blog *BlogEntry
	)

	if s == "deleted" {
		log.Println(s, *blogID)

		ID, blog = b.Create(entry)
		if ID == nil {
			return nil, nil, fmt.Errorf("deleted, but create failed")
		}

		log.Println("update blog with new ID:", *ID)
	}

	return ID, blog, nil
}

// Get is responsible for fetching blog from data storage
func (b *BlogStorage) Get(blogID *string) (*BlogEntry, error) {
	v, ok := b.storage[*blogID]
	if !ok {
		return nil, fmt.Errorf("blog ID not found")
	}

	return &v, nil

}

// Delete is responsible for deleting blog from data storage
func (b *BlogStorage) Delete(blogID *string) string {
	if _, ok := b.storage[*blogID]; !ok {
		return "blog ID not found"
	}

	delete(b.storage, *blogID)
	return "deleted"
}

// CreatePost is endoint for CREATE blog
func (b *BlogServer) CreatePost(ctx context.Context, in *proto.CreateRequest) (*proto.CreateResponse, error) {
	log.Printf("recvd create request for blog - %+v", in.BlogEntry)

	if in.BlogEntry.Title == "" {
		return nil, fmt.Errorf("no title")

	}

	if in.BlogEntry.Content == "" {
		return nil, fmt.Errorf("no content")

	}

	if in.BlogEntry.Author == "" {
		return nil, fmt.Errorf("no author")

	}

	if in.BlogEntry.PubDate == "" {
		return nil, fmt.Errorf("no publishing date")
	}

	id, blogEntry := blogStorage.Create(in.BlogEntry)

	log.Printf("new blog generated with ID: %v", *id)

	return &proto.CreateResponse{
		PostID: *id,
		BlogEntry: &proto.BlogEntryPost{
			Title:   blogEntry.Title,
			Content: blogEntry.Content,
			Author:  blogEntry.Author,
			PubDate: blogEntry.PubData,
			Tags:    blogEntry.Tags,
		},
	}, nil
}

// GetPost is endoint for GET blog
func (b *BlogServer) GetPost(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	log.Printf("recvd GET request for blog ID: %v", in.PostID)

	blog, err := blogStorage.Get(&in.PostID)
	if err != nil {
		return nil, err
	}

	return &proto.GetResponse{
		PostID: in.PostID,
		BlogEntry: &proto.BlogEntryPost{
			Title:   blog.Title,
			Content: blog.Content,
			Author:  blog.Author,
			PubDate: blog.PubData,
			Tags:    blog.Tags,
		},
	}, nil
}

// DeletePost is endoint for DELETE blog
func (b *BlogServer) DeletePost(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	log.Printf("Recvd delete for blog ID - %v", in.PostID)
	return &proto.DeleteResponse{Msg: blogStorage.Delete(&in.PostID)}, nil
}

// UpdatePost is endoint for PUT blog
func (b *BlogServer) UpdatePost(ctx context.Context, in *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	log.Println("Recvd update for ID:", in.PostID)

	id, blogEntry, err := blogStorage.Update(&in.PostID, in.BlogEntry)
	if err != nil {
		return nil, err
	}

	log.Println("update completed, new ID is:", *id)

	return &proto.UpdateResponse{
		PostID: *id,
		BlogEntry: &proto.BlogEntryPost{
			Title:   blogEntry.Title,
			Content: blogEntry.Content,
			Author:  blogEntry.Author,
			PubDate: blogEntry.PubData,
			Tags:    blogEntry.Tags,
		},
	}, nil
}
