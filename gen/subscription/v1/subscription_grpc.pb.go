package subscriptionv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andreysavvinov/petproject/internal/grpcjson"
)

type Subscription struct {
	UsersId     string `json:"users_id"`
	TrackListId string `json:"track_list_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type CreateSubscriptionRequest struct {
	UsersId     string `json:"users_id"`
	TrackListId string `json:"track_list_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type GetSubscriptionRequest struct {
	UsersId     string `json:"users_id"`
	TrackListId string `json:"track_list_id"`
}

type ListSubscriptionsRequest struct{}
type ListSubscriptionsResponse struct {
	Subscriptions []*Subscription `json:"subscriptions"`
}

type UpdateSubscriptionRequest struct {
	UsersId     string `json:"users_id"`
	TrackListId string `json:"track_list_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type DeleteSubscriptionRequest struct {
	UsersId     string `json:"users_id"`
	TrackListId string `json:"track_list_id"`
}

type DeleteResponse struct{ Success bool `json:"success"` }

type SubscriptionServiceServer interface {
	CreateSubscription(context.Context, *CreateSubscriptionRequest) (*Subscription, error)
	GetSubscription(context.Context, *GetSubscriptionRequest) (*Subscription, error)
	ListSubscriptions(context.Context, *ListSubscriptionsRequest) (*ListSubscriptionsResponse, error)
	UpdateSubscription(context.Context, *UpdateSubscriptionRequest) (*Subscription, error)
	DeleteSubscription(context.Context, *DeleteSubscriptionRequest) (*DeleteResponse, error)
}

type UnimplementedSubscriptionServiceServer struct{}

func (UnimplementedSubscriptionServiceServer) CreateSubscription(context.Context, *CreateSubscriptionRequest) (*Subscription, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSubscription not implemented")
}
func (UnimplementedSubscriptionServiceServer) GetSubscription(context.Context, *GetSubscriptionRequest) (*Subscription, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubscription not implemented")
}
func (UnimplementedSubscriptionServiceServer) ListSubscriptions(context.Context, *ListSubscriptionsRequest) (*ListSubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSubscriptions not implemented")
}
func (UnimplementedSubscriptionServiceServer) UpdateSubscription(context.Context, *UpdateSubscriptionRequest) (*Subscription, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSubscription not implemented")
}
func (UnimplementedSubscriptionServiceServer) DeleteSubscription(context.Context, *DeleteSubscriptionRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSubscription not implemented")
}

func RegisterSubscriptionServiceServer(s grpc.ServiceRegistrar, srv SubscriptionServiceServer) {
	s.RegisterService(&SubscriptionService_ServiceDesc, srv)
}

func unaryHandler[Req, Resp any](
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
	fullMethod string,
	fn func(SubscriptionServiceServer, context.Context, *Req) (*Resp, error),
) (any, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return fn(srv.(SubscriptionServiceServer), ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return fn(srv.(SubscriptionServiceServer), ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

var SubscriptionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "subscription.v1.SubscriptionService",
	HandlerType: (*SubscriptionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateSubscription", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[CreateSubscriptionRequest, Subscription](srv, ctx, dec, interceptor, "/subscription.v1.SubscriptionService/CreateSubscription", func(s SubscriptionServiceServer, ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error) {
				return s.CreateSubscription(ctx, req)
			})
		}},
		{MethodName: "GetSubscription", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[GetSubscriptionRequest, Subscription](srv, ctx, dec, interceptor, "/subscription.v1.SubscriptionService/GetSubscription", func(s SubscriptionServiceServer, ctx context.Context, req *GetSubscriptionRequest) (*Subscription, error) {
				return s.GetSubscription(ctx, req)
			})
		}},
		{MethodName: "ListSubscriptions", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[ListSubscriptionsRequest, ListSubscriptionsResponse](srv, ctx, dec, interceptor, "/subscription.v1.SubscriptionService/ListSubscriptions", func(s SubscriptionServiceServer, ctx context.Context, req *ListSubscriptionsRequest) (*ListSubscriptionsResponse, error) {
				return s.ListSubscriptions(ctx, req)
			})
		}},
		{MethodName: "UpdateSubscription", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[UpdateSubscriptionRequest, Subscription](srv, ctx, dec, interceptor, "/subscription.v1.SubscriptionService/UpdateSubscription", func(s SubscriptionServiceServer, ctx context.Context, req *UpdateSubscriptionRequest) (*Subscription, error) {
				return s.UpdateSubscription(ctx, req)
			})
		}},
		{MethodName: "DeleteSubscription", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[DeleteSubscriptionRequest, DeleteResponse](srv, ctx, dec, interceptor, "/subscription.v1.SubscriptionService/DeleteSubscription", func(s SubscriptionServiceServer, ctx context.Context, req *DeleteSubscriptionRequest) (*DeleteResponse, error) {
				return s.DeleteSubscription(ctx, req)
			})
		}},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/subscription/v1/subscription.proto",
}

type SubscriptionServiceClient interface {
	CreateSubscription(ctx context.Context, in *CreateSubscriptionRequest, opts ...grpc.CallOption) (*Subscription, error)
	GetSubscription(ctx context.Context, in *GetSubscriptionRequest, opts ...grpc.CallOption) (*Subscription, error)
	ListSubscriptions(ctx context.Context, in *ListSubscriptionsRequest, opts ...grpc.CallOption) (*ListSubscriptionsResponse, error)
	UpdateSubscription(ctx context.Context, in *UpdateSubscriptionRequest, opts ...grpc.CallOption) (*Subscription, error)
	DeleteSubscription(ctx context.Context, in *DeleteSubscriptionRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type subscriptionServiceClient struct{ cc grpc.ClientConnInterface }

func NewSubscriptionServiceClient(cc grpc.ClientConnInterface) SubscriptionServiceClient {
	return &subscriptionServiceClient{cc}
}

func (c *subscriptionServiceClient) invoke(ctx context.Context, method string, in, out any, opts ...grpc.CallOption) error {
	opts = append(opts, grpc.CallContentSubtype(grpcjson.Name))
	return c.cc.Invoke(ctx, method, in, out, opts...)
}

func (c *subscriptionServiceClient) CreateSubscription(ctx context.Context, in *CreateSubscriptionRequest, opts ...grpc.CallOption) (*Subscription, error) {
	out := new(Subscription)
	return out, c.invoke(ctx, "/subscription.v1.SubscriptionService/CreateSubscription", in, out, opts...)
}
func (c *subscriptionServiceClient) GetSubscription(ctx context.Context, in *GetSubscriptionRequest, opts ...grpc.CallOption) (*Subscription, error) {
	out := new(Subscription)
	return out, c.invoke(ctx, "/subscription.v1.SubscriptionService/GetSubscription", in, out, opts...)
}
func (c *subscriptionServiceClient) ListSubscriptions(ctx context.Context, in *ListSubscriptionsRequest, opts ...grpc.CallOption) (*ListSubscriptionsResponse, error) {
	out := new(ListSubscriptionsResponse)
	return out, c.invoke(ctx, "/subscription.v1.SubscriptionService/ListSubscriptions", in, out, opts...)
}
func (c *subscriptionServiceClient) UpdateSubscription(ctx context.Context, in *UpdateSubscriptionRequest, opts ...grpc.CallOption) (*Subscription, error) {
	out := new(Subscription)
	return out, c.invoke(ctx, "/subscription.v1.SubscriptionService/UpdateSubscription", in, out, opts...)
}
func (c *subscriptionServiceClient) DeleteSubscription(ctx context.Context, in *DeleteSubscriptionRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	return out, c.invoke(ctx, "/subscription.v1.SubscriptionService/DeleteSubscription", in, out, opts...)
}
