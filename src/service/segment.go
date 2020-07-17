package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/huichen/sego"
	"github.com/vidorg/vid_backend/src/config"
	"strings"
)

type SegmentService struct {
	Config *config.Config `di:"~"`

	Segmenter *sego.Segmenter `di:"-"`
}

func NewSegmentService(dic *xdi.DiContainer) *SegmentService {
	srv := &SegmentService{}
	dic.MustInject(srv)

	// TODO
	// var segmenter sego.Segmenter
	// segmenter.LoadDictionary(srv.Config.Search.DictPath)
	// srv.Segmenter = &segmenter

	return srv
}

func (s *SegmentService) Seg(str string) []string {
	// TODO
	// segments := s.Segmenter.Segment([]byte(str))
	// return sego.SegmentsToSlice(segments, true)
	return []string{""}
}

func (s *SegmentService) Cat(tokens []string) string {
	sign := "，。、？！；："
	ret := ""
	for _, token := range tokens {
		if !strings.Contains(sign, token) {
			ret += token + " "
		}
	}
	if len(ret) != 0 {
		ret = ret[0 : len(ret)-1]
	}
	return ret
}
