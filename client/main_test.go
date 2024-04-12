package main

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ezeeyahoo/demoBlogServiceInGrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	createBlogs = `
	{
		"blogs" : [
			{
				"Title": "blog1",
				"Content": "this is blog1",
				"Author": "alice",
				"Publication Date": "01/01/2009",
				"Tags": [
					"blog",
					"grpc"
				]
			},
			{
				"Title": "blog1",
				"Content": "this is blog2",
				"Author": "joe",
				"Publication Date": "01/12/2009",
				"Tags": [
					"blog",
					"grpc"
				]
			}
		]
	}`

	address string = "localhost"
	port    string = ":7878"
)

type blogItems struct {
	BlogTitle []struct {
		Title   string   `json:"Title"`
		Content string   `json:"Content"`
		Author  string   `json:"Author"`
		PubDate string   `json:"Publication Date"`
		Tags    []string `json:"Tags"`
	} `json:"blogs"`
}

func Test_main(t *testing.T) {

	blogs := blogItems{}
	if err := json.Unmarshal([]byte(createBlogs), &blogs); err != nil {
		t.Errorf("unmarshal error")
	}

	conn, err := grpc.Dial(string(address+port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	client := proto.NewBlogServicerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// create blog
	var blogID string
	var entry1 *proto.BlogEntryPost
	for _, blog := range blogs.BlogTitle {
		t.Logf("sending request for new blog- %+v", blog)

		postResp, err := client.CreatePost(ctx, &proto.CreateRequest{
			BlogEntry: &proto.BlogEntryPost{
				Title:   blog.Title,
				Author:  blog.Author,
				Content: blog.Content,
				PubDate: blog.PubDate,
				Tags:    blog.Tags,
			},
		})

		if err != nil {
			t.Errorf("POST blog failed- %v", err)
		}

		blogID, entry1 = postResp.PostID, postResp.BlogEntry
	}

	// get blog
	getResp, err := client.GetPost(ctx,
		&proto.GetRequest{PostID: blogID})
	if err != nil {
		t.Errorf("GET blog failed- %v", err)
	}

	entry2 := getResp.BlogEntry
	t.Logf("read blog output of ID %v : %+v", blogID, entry2)

	t.Logf("%+v \n %+v", entry1, entry2)
	if !reflect.DeepEqual(entry1, entry2) {
		t.Errorf("POST and GET returns different blog")
	}

	// update blog
	entry2.Content = entry2.Content + " updated!"

	putResp, err := client.UpdatePost(ctx,
		&proto.UpdateRequest{BlogEntry: entry2, PostID: blogID})
	if err != nil {
		t.Errorf("PUT blog failed- %v", err)
	}

	// delete blog
	delResp, err := client.DeletePost(ctx,
		&proto.DeleteRequest{PostID: putResp.PostID})
	if err != nil {
		t.Errorf("DELETE blog failed- %v", err)
	}

	if delResp.Msg != "deleted" {
		t.Errorf("delete failed")
	}
}
