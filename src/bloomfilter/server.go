package bloomfilter

import(
    "bufio"
    "io"
    "net"
    "strings"
	"fmt"
	"log"
	"io/ioutil"
)

type Server struct {
    Addr string
    handles map[string] HandlerFn
    bf *BloomFilter
}

type conn struct {
    // server *Server
    conn net.Conn
    rw *bufio.ReadWriter
}

func (s *Server) newConn(originalConn net.Conn) (*conn) {
	c := new(conn)
    // c.server = s
    c.conn = originalConn
    c.rw = bufio.NewReadWriter(bufio.NewReader(originalConn), bufio.NewWriter(originalConn))
    return c
}

func NewServer(address string) *Server {
    // return &Server{address}
    srv := &Server{
        Addr:address,
        handles:make(map[string] HandlerFn),
    }
    srv.initHandles()
    return srv
}

func (s *Server) ListenAndServer() error {
    // db, _ := OpenFile("/var/level/testdb", nil)
    // hashDb, _ := OpenFile("/var/level/testhashdb", nil)

    var network string
    var (
        addr net.Addr
        err error
    )
    if strings.Contains(s.Addr, "/") {
        network = "unix"
        addr, err = net.ResolveUnixAddr("unix", s.Addr)
    } else {
        network = "tcp"
        addr, err = net.ResolveTCPAddr("tcp", s.Addr)
    }
    if err != nil {
        return err
    }

    l, e := net.Listen(network, addr.String())
    if e != nil {
        return e
    }
    defer l.Close()
    for {
        connect, e := l.Accept() //type rw = net.Conn
        if e != nil {
            return e
        }
        c := s.newConn(connect)
        // go c.serve()
        go s.serveClient(c)
    }
}

func (s *Server) serveClient(c *conn) {
    defer func() {
        c.conn.Close()
    }()

    for {
        req, err := c.parseRequest()
        if err != nil {
            if err == io.EOF {
                return
            }
            c.rw.WriteString(err.Error())
        }
        rep := s.Apply(req)
        c.sendResponse(rep)
    }
}

func (s *Server) Apply(req *Request) *Reply {
    fn,exists := s.handles[strings.ToLower(req.method)]
    if !exists {
        rep,err := fn(req.data)
        if err != nil {
            log.Fatal(err)
        }
        return rep
    }
    return nil 
}

func (c *conn) sendResponse(r *Reply) {
    if r.responseType == "status" {
        c.rw.WriteString("+OK\r\n")
    }
    // c.rw.WriteString("+OK\r\n")
    if err := c.rw.Flush();err != nil {
        panic(err)
    }
    fmt.Println("call response")
}


type Request struct {
    method string
    data []string
}

func (c *conn) ReadLine() string {
    line, isPrefix, err := c.rw.ReadLine()
    if isPrefix || err != nil {
        panic(err)
    }
    return string(line)
}

func readArgument(r *bufio.Reader) (string, error) {
    line, err := r.ReadString('\n')
    var argSize int
    if _,err := fmt.Sscanf(line, "$%d\r", &argSize); err != nil {
        return "",err
    }
    data, err := ioutil.ReadAll(io.LimitReader(r, int64(argSize)))
    if err != nil {
        return "",err
    }
    if len(data) != argSize {
        return "", nil
    }
    if b, err := r.ReadByte(); err != nil || b != '\r' {
        panic(err)
        return "", nil
    }
    if b, err := r.ReadByte(); err != nil || b != '\n' {
        panic(err)
        return "", nil
    }
    return string(data), nil
}

func (c *conn) parseRequest() (*Request, error) {
    line, err := c.rw.ReadString('\n')
    if err != nil {
        if err == io.EOF {
            fmt.Println("no byte, the error is EOF")
            return nil, io.EOF
        } else {
            panic(err)
        }
    }

    var argsCount int
    if line[0] == '*' {
        if _,err := fmt.Sscanf(line, "*%d\r", &argsCount); err != nil {
            panic(err)
        }
		args := new(Request)
        args.method,err = readArgument(c.rw.Reader)
        for i := 0; i < argsCount-1; i += 1 {
            if args.data[i],err = readArgument(c.rw.Reader); err != nil {
                //return nil,err
            }
        }
        return args, nil
    }
    // fmt.Println(args)
    return nil, nil //error
}

type Reply struct {
    responseType string
    message string
}

//handle.go
type HandlerFn func([]string) (*Reply, error) //

// var add = func (r *Request) (*Reply, error) {
// func (s *Server) Add (req *Request) (*Reply, error) {
func (s *Server) Add (reqdata []string) (*Reply, error) {
    s.bf.Add([]byte(reqdata[0]))
    rep := &Reply{
        responseType:"status",
        message:"OK",
    }
    return rep,nil
}

func (s *Server) initHandles(){
    s.handles["add"] = s.Add
}

// h["hset"] = func(r *Request) (Reply, error) {
//  db.Hset(r.data[0],r.data[1],r.data[2])
//  // func Hset(key,field,value string)
// }
