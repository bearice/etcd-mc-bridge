package main

import (
    "io"
    "os"
    "fmt"
    "net"
    "log"
    "flag"
    "bytes"
    "bufio"
    "encoding/json"
    "github.com/zobo/mrproxy/protocol"
    "github.com/coreos/go-etcd/etcd"
)

var version = "MCBridge 0.0.1"
var cfile   = flag.String("etcd","config.json","path of etcd config file")
var mc_port = flag.Int("port", 22122, "port of memcached protocol")
var debug   = flag.Bool("debug",false, "debug enabled?")

var etcdClient *etcd.Client

func main() {
    var err error

    flag.Parse()

    if *debug {
        etcd.SetLogger(log.New(os.Stdout,"go-etcd", log.LstdFlags))
    }

    etcdClient,err = etcd.NewClientFromFile(*cfile)
    if err != nil {
        log.Fatal(err)
    }
    etcdClient.SyncCluster()

    log.Printf("Listening on mc_port:%v\n", *mc_port)
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *mc_port))
    if err != nil {
        log.Fatal(err)
    }

    sessId := 1
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go serve(conn, sessId)
        sessId = sessId + 1
    }
}

func serve(conn net.Conn, sessId int) {

    defer conn.Close()

    br := bufio.NewReader(conn)
    bw := bufio.NewWriter(conn)

    for {

        req, err := protocol.ReadRequest(br)
        res := protocol.McResponse{}

        if err != nil {
            if err == io.EOF {
                return
            }
            log.Print(err)
            res = protocol.McResponse{Response: "ERROR"}
        }else{
            switch req.Command {
            case "quit":
                return
            case "version":
                res = protocol.McResponse{Response: fmt.Sprintf("version %s", version)}
                //case "stats":
                //    res = getStats(idPool)
            case "get":
                fallthrough
            case "gets":
                for i, _ := range req.Keys {
                    result,err := etcdClient.Get(req.Keys[i], false, false)
                    if err != nil {
                        log.Print(err)
                    }else{
                        //log.Printf("result= %v",result.Node)
                        if req.Command == "get" {
                            if result.Node.Dir {
                                var b bytes.Buffer
                                for _,n := range result.Node.Nodes {
                                    b.WriteString(n.Key)
                                    b.WriteString("\n")
                                }
                                res.Values = append(res.Values, protocol.McValue{Key: req.Keys[i], Flags: "1", Data: b.Bytes() })
                            }else{
                                res.Values = append(res.Values, protocol.McValue{Key: req.Keys[i], Flags: "0", Data: []byte(result.Node.Value) })
                            }
                        }else{
                            val, err := json.Marshal(result.Node)
                            if err != nil {
                                log.Print(err)
                            } else {
                                res.Values = append(res.Values, protocol.McValue{Key: req.Keys[i], Flags: "0 0", Data: val})
                            }
                        }
                    }
                }
                res.Response = "END"
            default:
                res = protocol.McResponse{Response: "ERROR"}
            }
        }

        bw.WriteString(res.Protocol())
        bw.Flush()
    }
}

