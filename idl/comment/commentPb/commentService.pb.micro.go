// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: commentService.proto

package commentPb

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

// Api Endpoints for CommentService service

func NewCommentServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for CommentService service

type CommentService interface {
	CommentAction(ctx context.Context, in *CommentActionRequest, opts ...client.CallOption) (*CommentActionResponse, error)
	CommentList(ctx context.Context, in *CommentListRequest, opts ...client.CallOption) (*CommentListResponse, error)
}

type commentService struct {
	c    client.Client
	name string
}

func NewCommentService(name string, c client.Client) CommentService {
	return &commentService{
		c:    c,
		name: name,
	}
}

func (c *commentService) CommentAction(ctx context.Context, in *CommentActionRequest, opts ...client.CallOption) (*CommentActionResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.CommentAction", in)
	out := new(CommentActionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentService) CommentList(ctx context.Context, in *CommentListRequest, opts ...client.CallOption) (*CommentListResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.CommentList", in)
	out := new(CommentListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CommentService service

type CommentServiceHandler interface {
	CommentAction(context.Context, *CommentActionRequest, *CommentActionResponse) error
	CommentList(context.Context, *CommentListRequest, *CommentListResponse) error
}

func RegisterCommentServiceHandler(s server.Server, hdlr CommentServiceHandler, opts ...server.HandlerOption) error {
	type commentService interface {
		CommentAction(ctx context.Context, in *CommentActionRequest, out *CommentActionResponse) error
		CommentList(ctx context.Context, in *CommentListRequest, out *CommentListResponse) error
	}
	type CommentService struct {
		commentService
	}
	h := &commentServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&CommentService{h}, opts...))
}

type commentServiceHandler struct {
	CommentServiceHandler
}

func (h *commentServiceHandler) CommentAction(ctx context.Context, in *CommentActionRequest, out *CommentActionResponse) error {
	return h.CommentServiceHandler.CommentAction(ctx, in, out)
}

func (h *commentServiceHandler) CommentList(ctx context.Context, in *CommentListRequest, out *CommentListResponse) error {
	return h.CommentServiceHandler.CommentList(ctx, in, out)
}
