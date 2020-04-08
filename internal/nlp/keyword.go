package nlp

import (
	"context"
	"time"

	"github.com/vivym/ncovis-api/internal/api/protobuf/nlp"
	"github.com/vivym/ncovis-api/internal/model"
)

func (n *NLPToolkit) ExtractKeywords(sentence string, topK int64) ([]model.Keyword, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	request := nlp.KeywordRequest{
		Method:   nlp.KeywordRequest_TextRank,
		Sentence: sentence,
		TopK:     topK,
		AllowPOS: []string{"n", "nr", "nz", "ns", "v", "s", "nt", "nw", "vn"},
	}

	rsp, err := n.client.ExtractKeywords(ctx, &request)
	if err != nil {
		return nil, err
	}

	keywords := make([]model.Keyword, 0, len(rsp.GetKeywords()))
	for _, k := range rsp.GetKeywords() {
		keywords = append(keywords, model.Keyword{
			Name:   k.Word,
			Weight: k.Weight,
			POS:    k.Pos,
		})
	}

	return keywords, nil
}
