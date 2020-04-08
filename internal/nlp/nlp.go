package nlp

import (
	"github.com/vivym/ncovis-api/internal/api/protobuf/nlp"
	"google.golang.org/grpc"
)

type NLPToolkit struct {
	conn   *grpc.ClientConn
	client nlp.NLPClient
}

var toolkit *NLPToolkit

func Init(config Config) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial(config.Address, opts...)
	if err != nil {
		return err
	}

	client := nlp.NewNLPClient(conn)

	toolkit = &NLPToolkit{
		conn:   conn,
		client: client,
	}
	return nil
}

func Get() *NLPToolkit {
	return toolkit
}

func (n *NLPToolkit) Release() {
	if n.conn != nil {
		_ = n.conn.Close()
	}
}
