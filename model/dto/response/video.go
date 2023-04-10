package response

import "faker-douyin/model/entity"

type PublishVideoRes struct {
	Video entity.TableVideo `json:"video"`
}
