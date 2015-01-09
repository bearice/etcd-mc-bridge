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
