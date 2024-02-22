// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: videoService.proto

package videoPb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for VideoService service

func NewVideoServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for VideoService service

type VideoService interface {
	Feed(ctx context.Context, in *FeedRequest, opts ...client.CallOption) (*FeedResponse, error)
	Publish(ctx context.Context, in *PublishRequest, opts ...client.CallOption) (*PublishResponse, error)
	PublishList(ctx context.Context, in *PublishListRequest, opts ...client.CallOption) (*PublishListResponse, error)
}

type videoService struct {
	c    client.Client
	name string
}

func NewVideoService(name string, c client.Client) VideoService {
	return &videoService{
		c:    c,
		name: name,
	}
}

func (c *videoService) Feed(ctx context.Context, in *FeedRequest, opts ...client.CallOption) (*FeedResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.Feed", in)
	out := new(FeedResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) Publish(ctx context.Context, in *PublishRequest, opts ...client.CallOption) (*PublishResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.Publish", in)
	out := new(PublishResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoService) PublishList(ctx context.Context, in *PublishListRequest, opts ...client.CallOption) (*PublishListResponse, error) {
	req := c.c.NewRequest(c.name, "VideoService.PublishList", in)
	out := new(PublishListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VideoService service

type VideoServiceHandler interface {
	Feed(context.Context, *FeedRequest, *FeedResponse) error
	Publish(context.Context, *PublishRequest, *PublishResponse) error
	PublishList(context.Context, *PublishListRequest, *PublishListResponse) error
}

func RegisterVideoServiceHandler(s server.Server, hdlr VideoServiceHandler, opts ...server.HandlerOption) error {
	type videoService interface {
		Feed(ctx context.Context, in *FeedRequest, out *FeedResponse) error
		Publish(ctx context.Context, in *PublishRequest, out *PublishResponse) error
		PublishList(ctx context.Context, in *PublishListRequest, out *PublishListResponse) error
	}
	type VideoService struct {
		videoService
	}
	h := &videoServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&VideoService{h}, opts...))
}

type videoServiceHandler struct {
	VideoServiceHandler
}

func (h *videoServiceHandler) Feed(ctx context.Context, in *FeedRequest, out *FeedResponse) error {
	return h.VideoServiceHandler.Feed(ctx, in, out)
}

func (h *videoServiceHandler) Publish(ctx context.Context, in *PublishRequest, out *PublishResponse) error {
	return h.VideoServiceHandler.Publish(ctx, in, out)
}

func (h *videoServiceHandler) PublishList(ctx context.Context, in *PublishListRequest, out *PublishListResponse) error {
	return h.VideoServiceHandler.PublishList(ctx, in, out)
}