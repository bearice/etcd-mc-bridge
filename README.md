etcd-mc-bridge
===

A simple bridge translates memcached protocol to etcd clusters

usage:

```
./etcd-mc-bridge -port 11211 -etcd config.json
```

supported opreations:
* `get` : return string value of etcd node, or lists of child nodes if node is a dir.
* `gets` : return json of node.

```
 %echo get / | nc master 22122                                                                                       
VALUE / 1 27
/dnsmasq
/fig
/nginx
/test

END

%echo gets / | nc master 22122                                                                                      
VALUE / 0 302
{"key":"/","dir":true,"nodes":[{"key":"/dnsmasq","dir":true,"modifiedIndex":164,"createdIndex":164},{"key":"/fig","dir":true,"modifiedIndex":319,"createdIndex":319},{"key":"/nginx","dir":true,"modifiedIndex":5,"createdIndex":5},{"key":"/test","value":"test","modifiedIndex":4061,"createdIndex":4061}]}
END
```
