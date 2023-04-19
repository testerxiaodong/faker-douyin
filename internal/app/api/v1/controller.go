package v1

import "github.com/google/wire"

var ProviderSet = wire.NewSet(wire.Struct(new(CommentController), "*"), wire.Struct(new(UserController), "*"), wire.Struct(new(VideoController), "*"))
