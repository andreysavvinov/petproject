package catalogv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andreysavvinov/petproject/internal/grpcjson"
)

type Album struct {
	Id          string `json:"id"`
	AlbumName   string `json:"album_name"`
	ReleaseYear string `json:"release_year"`
	Genre       string `json:"genre"`
}

type Performer struct {
	Id            string `json:"id"`
	PerformerName string `json:"performer_name"`
}

type Track struct {
	Id           string `json:"id"`
	TrackName    string `json:"track_name"`
	Genre        string `json:"genre"`
	Duration     string `json:"duration"`
	AlbumsId     string `json:"albums_id"`
	PerformersId string `json:"performers_id"`
	MyMusicId    int64  `json:"my_music_id"`
}

type CreateAlbumRequest struct {
	AlbumName   string `json:"album_name"`
	ReleaseYear string `json:"release_year"`
	Genre       string `json:"genre"`
}

type GetAlbumRequest struct{ Id string `json:"id"` }
type ListAlbumsRequest struct{}
type ListAlbumsResponse struct{ Albums []*Album `json:"albums"` }

type UpdateAlbumRequest struct {
	Id          string `json:"id"`
	AlbumName   string `json:"album_name"`
	ReleaseYear string `json:"release_year"`
	Genre       string `json:"genre"`
}

type DeleteAlbumRequest struct{ Id string `json:"id"` }

type CreatePerformerRequest struct{ PerformerName string `json:"performer_name"` }
type GetPerformerRequest struct{ Id string `json:"id"` }
type ListPerformersRequest struct{}
type ListPerformersResponse struct{ Performers []*Performer `json:"performers"` }

type UpdatePerformerRequest struct {
	Id            string `json:"id"`
	PerformerName string `json:"performer_name"`
}

type DeletePerformerRequest struct{ Id string `json:"id"` }

type CreateTrackRequest struct {
	TrackName    string `json:"track_name"`
	Genre        string `json:"genre"`
	Duration     string `json:"duration"`
	AlbumsId     string `json:"albums_id"`
	PerformersId string `json:"performers_id"`
	MyMusicId    int64  `json:"my_music_id"`
}

type GetTrackRequest struct{ Id string `json:"id"` }
type ListTracksRequest struct{}
type ListTracksResponse struct{ Tracks []*Track `json:"tracks"` }

type UpdateTrackRequest struct {
	Id           string `json:"id"`
	TrackName    string `json:"track_name"`
	Genre        string `json:"genre"`
	Duration     string `json:"duration"`
	AlbumsId     string `json:"albums_id"`
	PerformersId string `json:"performers_id"`
	MyMusicId    int64  `json:"my_music_id"`
}

type DeleteTrackRequest struct{ Id string `json:"id"` }
type DeleteResponse struct{ Success bool `json:"success"` }

type CatalogServiceServer interface {
	CreateAlbum(context.Context, *CreateAlbumRequest) (*Album, error)
	GetAlbum(context.Context, *GetAlbumRequest) (*Album, error)
	ListAlbums(context.Context, *ListAlbumsRequest) (*ListAlbumsResponse, error)
	UpdateAlbum(context.Context, *UpdateAlbumRequest) (*Album, error)
	DeleteAlbum(context.Context, *DeleteAlbumRequest) (*DeleteResponse, error)

	CreatePerformer(context.Context, *CreatePerformerRequest) (*Performer, error)
	GetPerformer(context.Context, *GetPerformerRequest) (*Performer, error)
	ListPerformers(context.Context, *ListPerformersRequest) (*ListPerformersResponse, error)
	UpdatePerformer(context.Context, *UpdatePerformerRequest) (*Performer, error)
	DeletePerformer(context.Context, *DeletePerformerRequest) (*DeleteResponse, error)

	CreateTrack(context.Context, *CreateTrackRequest) (*Track, error)
	GetTrack(context.Context, *GetTrackRequest) (*Track, error)
	ListTracks(context.Context, *ListTracksRequest) (*ListTracksResponse, error)
	UpdateTrack(context.Context, *UpdateTrackRequest) (*Track, error)
	DeleteTrack(context.Context, *DeleteTrackRequest) (*DeleteResponse, error)
}

type UnimplementedCatalogServiceServer struct{}

func (UnimplementedCatalogServiceServer) CreateAlbum(context.Context, *CreateAlbumRequest) (*Album, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlbum not implemented")
}
func (UnimplementedCatalogServiceServer) GetAlbum(context.Context, *GetAlbumRequest) (*Album, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbum not implemented")
}
func (UnimplementedCatalogServiceServer) ListAlbums(context.Context, *ListAlbumsRequest) (*ListAlbumsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAlbums not implemented")
}
func (UnimplementedCatalogServiceServer) UpdateAlbum(context.Context, *UpdateAlbumRequest) (*Album, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAlbum not implemented")
}
func (UnimplementedCatalogServiceServer) DeleteAlbum(context.Context, *DeleteAlbumRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAlbum not implemented")
}
func (UnimplementedCatalogServiceServer) CreatePerformer(context.Context, *CreatePerformerRequest) (*Performer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePerformer not implemented")
}
func (UnimplementedCatalogServiceServer) GetPerformer(context.Context, *GetPerformerRequest) (*Performer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPerformer not implemented")
}
func (UnimplementedCatalogServiceServer) ListPerformers(context.Context, *ListPerformersRequest) (*ListPerformersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPerformers not implemented")
}
func (UnimplementedCatalogServiceServer) UpdatePerformer(context.Context, *UpdatePerformerRequest) (*Performer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePerformer not implemented")
}
func (UnimplementedCatalogServiceServer) DeletePerformer(context.Context, *DeletePerformerRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePerformer not implemented")
}
func (UnimplementedCatalogServiceServer) CreateTrack(context.Context, *CreateTrackRequest) (*Track, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTrack not implemented")
}
func (UnimplementedCatalogServiceServer) GetTrack(context.Context, *GetTrackRequest) (*Track, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrack not implemented")
}
func (UnimplementedCatalogServiceServer) ListTracks(context.Context, *ListTracksRequest) (*ListTracksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTracks not implemented")
}
func (UnimplementedCatalogServiceServer) UpdateTrack(context.Context, *UpdateTrackRequest) (*Track, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTrack not implemented")
}
func (UnimplementedCatalogServiceServer) DeleteTrack(context.Context, *DeleteTrackRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTrack not implemented")
}

func RegisterCatalogServiceServer(s grpc.ServiceRegistrar, srv CatalogServiceServer) {
	s.RegisterService(&CatalogService_ServiceDesc, srv)
}

func unaryHandler[Req, Resp any](
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
	fullMethod string,
	fn func(CatalogServiceServer, context.Context, *Req) (*Resp, error),
) (any, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return fn(srv.(CatalogServiceServer), ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return fn(srv.(CatalogServiceServer), ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

var CatalogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "catalog.v1.CatalogService",
	HandlerType: (*CatalogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateAlbum", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[CreateAlbumRequest, Album](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/CreateAlbum", func(s CatalogServiceServer, ctx context.Context, req *CreateAlbumRequest) (*Album, error) {
				return s.CreateAlbum(ctx, req)
			})
		}},
		{MethodName: "GetAlbum", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[GetAlbumRequest, Album](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/GetAlbum", func(s CatalogServiceServer, ctx context.Context, req *GetAlbumRequest) (*Album, error) {
				return s.GetAlbum(ctx, req)
			})
		}},
		{MethodName: "ListAlbums", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[ListAlbumsRequest, ListAlbumsResponse](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/ListAlbums", func(s CatalogServiceServer, ctx context.Context, req *ListAlbumsRequest) (*ListAlbumsResponse, error) {
				return s.ListAlbums(ctx, req)
			})
		}},
		{MethodName: "UpdateAlbum", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[UpdateAlbumRequest, Album](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/UpdateAlbum", func(s CatalogServiceServer, ctx context.Context, req *UpdateAlbumRequest) (*Album, error) {
				return s.UpdateAlbum(ctx, req)
			})
		}},
		{MethodName: "DeleteAlbum", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[DeleteAlbumRequest, DeleteResponse](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/DeleteAlbum", func(s CatalogServiceServer, ctx context.Context, req *DeleteAlbumRequest) (*DeleteResponse, error) {
				return s.DeleteAlbum(ctx, req)
			})
		}},
		{MethodName: "CreatePerformer", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[CreatePerformerRequest, Performer](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/CreatePerformer", func(s CatalogServiceServer, ctx context.Context, req *CreatePerformerRequest) (*Performer, error) {
				return s.CreatePerformer(ctx, req)
			})
		}},
		{MethodName: "GetPerformer", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[GetPerformerRequest, Performer](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/GetPerformer", func(s CatalogServiceServer, ctx context.Context, req *GetPerformerRequest) (*Performer, error) {
				return s.GetPerformer(ctx, req)
			})
		}},
		{MethodName: "ListPerformers", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[ListPerformersRequest, ListPerformersResponse](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/ListPerformers", func(s CatalogServiceServer, ctx context.Context, req *ListPerformersRequest) (*ListPerformersResponse, error) {
				return s.ListPerformers(ctx, req)
			})
		}},
		{MethodName: "UpdatePerformer", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[UpdatePerformerRequest, Performer](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/UpdatePerformer", func(s CatalogServiceServer, ctx context.Context, req *UpdatePerformerRequest) (*Performer, error) {
				return s.UpdatePerformer(ctx, req)
			})
		}},
		{MethodName: "DeletePerformer", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[DeletePerformerRequest, DeleteResponse](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/DeletePerformer", func(s CatalogServiceServer, ctx context.Context, req *DeletePerformerRequest) (*DeleteResponse, error) {
				return s.DeletePerformer(ctx, req)
			})
		}},
		{MethodName: "CreateTrack", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[CreateTrackRequest, Track](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/CreateTrack", func(s CatalogServiceServer, ctx context.Context, req *CreateTrackRequest) (*Track, error) {
				return s.CreateTrack(ctx, req)
			})
		}},
		{MethodName: "GetTrack", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[GetTrackRequest, Track](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/GetTrack", func(s CatalogServiceServer, ctx context.Context, req *GetTrackRequest) (*Track, error) {
				return s.GetTrack(ctx, req)
			})
		}},
		{MethodName: "ListTracks", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[ListTracksRequest, ListTracksResponse](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/ListTracks", func(s CatalogServiceServer, ctx context.Context, req *ListTracksRequest) (*ListTracksResponse, error) {
				return s.ListTracks(ctx, req)
			})
		}},
		{MethodName: "UpdateTrack", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[UpdateTrackRequest, Track](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/UpdateTrack", func(s CatalogServiceServer, ctx context.Context, req *UpdateTrackRequest) (*Track, error) {
				return s.UpdateTrack(ctx, req)
			})
		}},
		{MethodName: "DeleteTrack", Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return unaryHandler[DeleteTrackRequest, DeleteResponse](srv, ctx, dec, interceptor, "/catalog.v1.CatalogService/DeleteTrack", func(s CatalogServiceServer, ctx context.Context, req *DeleteTrackRequest) (*DeleteResponse, error) {
				return s.DeleteTrack(ctx, req)
			})
		}},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/catalog/v1/catalog.proto",
}

type CatalogServiceClient interface {
	CreateAlbum(ctx context.Context, in *CreateAlbumRequest, opts ...grpc.CallOption) (*Album, error)
	GetAlbum(ctx context.Context, in *GetAlbumRequest, opts ...grpc.CallOption) (*Album, error)
	ListAlbums(ctx context.Context, in *ListAlbumsRequest, opts ...grpc.CallOption) (*ListAlbumsResponse, error)
	UpdateAlbum(ctx context.Context, in *UpdateAlbumRequest, opts ...grpc.CallOption) (*Album, error)
	DeleteAlbum(ctx context.Context, in *DeleteAlbumRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	CreatePerformer(ctx context.Context, in *CreatePerformerRequest, opts ...grpc.CallOption) (*Performer, error)
	GetPerformer(ctx context.Context, in *GetPerformerRequest, opts ...grpc.CallOption) (*Performer, error)
	ListPerformers(ctx context.Context, in *ListPerformersRequest, opts ...grpc.CallOption) (*ListPerformersResponse, error)
	UpdatePerformer(ctx context.Context, in *UpdatePerformerRequest, opts ...grpc.CallOption) (*Performer, error)
	DeletePerformer(ctx context.Context, in *DeletePerformerRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	CreateTrack(ctx context.Context, in *CreateTrackRequest, opts ...grpc.CallOption) (*Track, error)
	GetTrack(ctx context.Context, in *GetTrackRequest, opts ...grpc.CallOption) (*Track, error)
	ListTracks(ctx context.Context, in *ListTracksRequest, opts ...grpc.CallOption) (*ListTracksResponse, error)
	UpdateTrack(ctx context.Context, in *UpdateTrackRequest, opts ...grpc.CallOption) (*Track, error)
	DeleteTrack(ctx context.Context, in *DeleteTrackRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type catalogServiceClient struct{ cc grpc.ClientConnInterface }

func NewCatalogServiceClient(cc grpc.ClientConnInterface) CatalogServiceClient {
	return &catalogServiceClient{cc}
}

func (c *catalogServiceClient) invoke(ctx context.Context, method string, in, out any, opts ...grpc.CallOption) error {
	opts = append(opts, grpc.CallContentSubtype(grpcjson.Name))
	return c.cc.Invoke(ctx, method, in, out, opts...)
}

func (c *catalogServiceClient) CreateAlbum(ctx context.Context, in *CreateAlbumRequest, opts ...grpc.CallOption) (*Album, error) {
	out := new(Album)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/CreateAlbum", in, out, opts...)
}
func (c *catalogServiceClient) GetAlbum(ctx context.Context, in *GetAlbumRequest, opts ...grpc.CallOption) (*Album, error) {
	out := new(Album)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/GetAlbum", in, out, opts...)
}
func (c *catalogServiceClient) ListAlbums(ctx context.Context, in *ListAlbumsRequest, opts ...grpc.CallOption) (*ListAlbumsResponse, error) {
	out := new(ListAlbumsResponse)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/ListAlbums", in, out, opts...)
}
func (c *catalogServiceClient) UpdateAlbum(ctx context.Context, in *UpdateAlbumRequest, opts ...grpc.CallOption) (*Album, error) {
	out := new(Album)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/UpdateAlbum", in, out, opts...)
}
func (c *catalogServiceClient) DeleteAlbum(ctx context.Context, in *DeleteAlbumRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/DeleteAlbum", in, out, opts...)
}
func (c *catalogServiceClient) CreatePerformer(ctx context.Context, in *CreatePerformerRequest, opts ...grpc.CallOption) (*Performer, error) {
	out := new(Performer)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/CreatePerformer", in, out, opts...)
}
func (c *catalogServiceClient) GetPerformer(ctx context.Context, in *GetPerformerRequest, opts ...grpc.CallOption) (*Performer, error) {
	out := new(Performer)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/GetPerformer", in, out, opts...)
}
func (c *catalogServiceClient) ListPerformers(ctx context.Context, in *ListPerformersRequest, opts ...grpc.CallOption) (*ListPerformersResponse, error) {
	out := new(ListPerformersResponse)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/ListPerformers", in, out, opts...)
}
func (c *catalogServiceClient) UpdatePerformer(ctx context.Context, in *UpdatePerformerRequest, opts ...grpc.CallOption) (*Performer, error) {
	out := new(Performer)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/UpdatePerformer", in, out, opts...)
}
func (c *catalogServiceClient) DeletePerformer(ctx context.Context, in *DeletePerformerRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/DeletePerformer", in, out, opts...)
}
func (c *catalogServiceClient) CreateTrack(ctx context.Context, in *CreateTrackRequest, opts ...grpc.CallOption) (*Track, error) {
	out := new(Track)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/CreateTrack", in, out, opts...)
}
func (c *catalogServiceClient) GetTrack(ctx context.Context, in *GetTrackRequest, opts ...grpc.CallOption) (*Track, error) {
	out := new(Track)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/GetTrack", in, out, opts...)
}
func (c *catalogServiceClient) ListTracks(ctx context.Context, in *ListTracksRequest, opts ...grpc.CallOption) (*ListTracksResponse, error) {
	out := new(ListTracksResponse)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/ListTracks", in, out, opts...)
}
func (c *catalogServiceClient) UpdateTrack(ctx context.Context, in *UpdateTrackRequest, opts ...grpc.CallOption) (*Track, error) {
	out := new(Track)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/UpdateTrack", in, out, opts...)
}
func (c *catalogServiceClient) DeleteTrack(ctx context.Context, in *DeleteTrackRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	return out, c.invoke(ctx, "/catalog.v1.CatalogService/DeleteTrack", in, out, opts...)
}
