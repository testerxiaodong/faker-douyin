package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(wire.Struct(new(UserServiceImpl), "*"), wire.Struct(new(CommentServiceImpl), "*"), wire.Struct(new(VideoServiceImpl), "*"),
	wire.Bind(new(UserService), new(*UserServiceImpl)), wire.Bind(new(VideoService), new(*VideoServiceImpl)), wire.Bind(new(CommentService), new(*CommentServiceImpl)))
