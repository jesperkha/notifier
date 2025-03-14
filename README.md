# notifier

This package implements a `Notifier`, which wraps `context.Context` and `sync.WaitGroup`, making a bidirectional synchronization tool. Passing the notifier to goroutines lets them register for a shutdown signal, and then communicate back when they have finished cleaning up.

```sh
go get github.com/jesperkha/notifier
```
