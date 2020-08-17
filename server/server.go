package server

import (
	"fmt"
	// "net"
	"crypto/rand"
	"crypto/tls"
	"time"

	"github.com/c-delta/torii/utils"
	reuse "github.com/libp2p/go-reuseport"
	tlsreuse "github.com/c-delta/go-tlsreuse"
)

var toriilogo = `
      _/_/_/_/_/_/_/_/_/
       _/          _/
    _/_/_/_/_/_/_/_/_/  _/_/_/_/_/                    _/  _/  
     _/          _/        _/      _/_/    _/  _/_/            
    _/          _/        _/    _/    _/  _/_/      _/  _/     
   _/          _/        _/    _/    _/  _/        _/  _/      
  _/          _/        _/      _/_/    _/        _/  _/       
`
var shintoTorii = "⛩️"
var horizonBorder = `─────────────────────────────────────────────────────────────`
var topLeftBorder = `┌`
var topRightBorder = `┐`
var bottomLeftBorder = `└`
var bottomRightBorder = `┘`
var verticalBorder = `│`

// Server
type Server struct {
	Name      string
	IPVersion int
	Protocol  string
	Host      string
	Port      int32
	Config    *utils.Settings
	Handler   *HandlerManager
	MaxHosts  int32
}

// NewServer 載入設定檔
func NewServer(number int32) *Server {
	Config := utils.Config()
	server := &Server{
		Name:     Config.Name,
		Protocol: "tcp4",
		Config:   Config,
		Host:     Config.Host,
		Port:     Config.Port,
		Handler:  NewHandlerManager(),
		MaxHosts: number,
	}
	return server
}

// Start 啟動Server
func (s *Server) Start() {
	for i := int32(0); i < s.MaxHosts; i++ {
		go s.newServer()
	}

	select {}
}

func (s *Server) Stop() {

}

// NewTask 建立新處理任務
func (s *Server) NewTask(h *Handler) error {
	err := s.Handler.Add(h)
	return err
}

func (s *Server) newServer() {
	if !s.Config.SSL.Enable {

		// fmt.Println("[Status]", color.Green("Start"), " Server")
		fmt.Printf("Name: %s | Host: %s | Port: %d | Protocol: %s\n", s.Name, s.Host, s.Port, s.Protocol)

		listener, err := reuse.Listen(fmt.Sprintf("%s", s.Protocol), fmt.Sprintf("%s:%d", s.Host, s.Port))
		if err != nil {
			panic(err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			} else {
				defer conn.Close()
			}
			go s.Handler.AcceptTasks(conn)
		}
	} else {
		// fmt.Println("[Status]", color.Green("Start"), " Server")
		fmt.Printf("Name: %s | Host: %s | Port: %d | Protocol: %s\n", s.Name, s.Host, s.Port, s.Protocol)

		crt, err := tls.LoadX509KeyPair(s.Config.SSL.Cert, s.Config.SSL.Key)
		if err != nil {
			panic(err)
		}
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{crt}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader

		listener, err := tlsreuse.Listen(fmt.Sprintf("%s", s.Protocol), fmt.Sprintf("%s:%d", s.Host, s.Port), tlsConfig)
		if err != nil {
			panic(err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			} else {
				defer conn.Close()
			}
			go s.Handler.AcceptTasks(conn)
		}
	}

	select {}
}

func init() {
	fmt.Println(toriilogo)

	fmt.Printf("%s%s%s\n", topLeftBorder,
		horizonBorder,
		topRightBorder)

	fmt.Printf("%s %s %-58v%s\n",
		verticalBorder,
		shintoTorii,
		" [Github] https://github.com/c-delta/torii",
		verticalBorder)

	fmt.Printf("%s %s %-58v%s\n",
		verticalBorder,
		shintoTorii,
		" [Version] v0.1 Alpha",
		verticalBorder)

	fmt.Printf("%s%s%s\n", bottomLeftBorder,
		horizonBorder,
		bottomRightBorder)

	fmt.Printf("Server Version: %s, MaxPacketSize: %d\n",
		utils.Config().Version,
		utils.Config().MaxPacketSize)
}
