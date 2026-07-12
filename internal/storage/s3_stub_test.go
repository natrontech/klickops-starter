package storage

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/natrontech/klickops-starter/internal/api"
	"github.com/natrontech/klickops-starter/internal/config"
)

// s3Stub is a minimal in-memory S3 server (path-style) so the store is
// tested at the HTTP boundary without a real bucket.
type s3Stub struct {
	objects map[string]stubObject
}

type stubObject struct {
	body        []byte
	contentType string
}

func (s *s3Stub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/test-bucket")
	key = strings.TrimPrefix(key, "/")

	switch {
	case r.Method == "HEAD" && key == "":
		w.WriteHeader(http.StatusOK)
	case r.Method == "GET" && key == "":
		type contents struct {
			Key          string
			Size         int
			LastModified string
		}
		var list []contents
		for k, o := range s.objects {
			list = append(list, contents{k, len(o.body), time.Now().UTC().Format(time.RFC3339)})
		}
		out, _ := xml.Marshal(struct {
			XMLName  xml.Name `xml:"ListBucketResult"`
			Name     string
			KeyCount int
			Contents []contents
		}{Name: "test-bucket", KeyCount: len(list), Contents: list})
		w.Header().Set("Content-Type", "application/xml")
		w.Write(out)
	case r.Method == "PUT":
		body, _ := io.ReadAll(r.Body)
		s.objects[key] = stubObject{body: body, contentType: r.Header.Get("Content-Type")}
		w.WriteHeader(http.StatusOK)
	case r.Method == "GET":
		obj, ok := s.objects[key]
		if !ok {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `<?xml version="1.0"?><Error><Code>NoSuchKey</Code><Message>no such key</Message></Error>`)
			return
		}
		w.Header().Set("Content-Type", obj.contentType)
		w.Write(obj.body)
	case r.Method == "DELETE":
		delete(s.objects, key)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func TestS3StoreRoundTrip(t *testing.T) {
	stub := &s3Stub{objects: map[string]stubObject{}}
	server := httptest.NewServer(stub)
	defer server.Close()

	store, err := NewS3(context.Background(), config.S3{
		Endpoint:  server.URL,
		Region:    "us-east-1",
		Bucket:    "test-bucket",
		AccessKey: "test",
		SecretKey: "test",
	})
	if err != nil {
		t.Fatalf("NewS3: %v", err)
	}
	ctx := context.Background()

	content := []byte("hello starter")
	if err := store.PutFile(ctx, "hello.txt", "text/plain", bytes.NewReader(content), int64(len(content))); err != nil {
		t.Fatalf("PutFile: %v", err)
	}

	body, contentType, err := store.GetFile(ctx, "hello.txt")
	if err != nil {
		t.Fatalf("GetFile: %v", err)
	}
	got, _ := io.ReadAll(body)
	body.Close()
	if !bytes.Equal(got, content) {
		t.Errorf("GetFile body = %q, want %q", got, content)
	}
	if contentType != "text/plain" {
		t.Errorf("contentType = %q, want text/plain", contentType)
	}

	files, err := store.ListFiles(ctx)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 1 || files[0].Key != "hello.txt" {
		t.Errorf("ListFiles = %+v, want one entry hello.txt", files)
	}

	if _, _, err := store.GetFile(ctx, "missing.txt"); err != api.ErrNotFound {
		t.Errorf("GetFile(missing) err = %v, want api.ErrNotFound", err)
	}

	if err := store.DeleteFile(ctx, "hello.txt"); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}
	if files, _ := store.ListFiles(ctx); len(files) != 0 {
		t.Errorf("after delete: %d files left", len(files))
	}
}
