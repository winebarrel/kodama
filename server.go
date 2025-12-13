package kodama

type Server struct{}

func NewServer(options *Options) *Server {
	server := &Server{}
	return server
}

func (svr *Server) Start() error {
	return nil
}
