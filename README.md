# k8s-port-forward

Straightforward and naive implementation of port forwarding of an existing K8S Pod built on top of K8S API.

The package was initially developed for internal use only to provide K8S port forwarding capabilities for a 3rd-party app (not Golang based). The idea was not to get locked on an existence of `kubectl` tool as well as providing of porcelain output.

`cmd/sample_app.go` provides with a sample usage (almost 100% close to a real life) of the package.

# License

The library is released under the MIT license. See LICENSE file.
