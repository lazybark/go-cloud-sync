# go-cloud-sync
The repository holds basic functionality for:
* watching for filesystem events
* processing fs events and comparing to records in database
* exchange events between peers (generally known as client & server, but it's posible to use server-server mode)
* exchange files between peers (client-server, server-server)

This code provides a lib to use in your own projects and production ready client/server code to start syncing files.
<br>
## Projects used here
* fsnotify
* lazyevent
* go-tls-server & go-tls-client
