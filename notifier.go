package notifier

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// Notifier wraps context.Context and sync.WaitGroup, making a bidirectional
// synchronization tool. Passing the notifier to goroutines lets them register
// for a shutdown signal, and then communicate back when they have finished
// cleaning up.
type Notifier struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Notifier {
	ctx, cancel := context.WithCancel(context.Background())
	return &Notifier{
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Register a listener. Returns a done channel which will close
// when the shutdown signal is sent, and a finish function to call when
// cleanup is completed.
func (n *Notifier) Register() (done <-chan struct{}, finish func()) {
	n.wg.Add(1)
	return n.ctx.Done(), func() {
		n.wg.Done()
	}
}

// Wait until all goroutines have finished cleaning up. Does not send shutdown
// signal.
func (n *Notifier) Wait() {
	n.wg.Wait()
}

// Notify sends the shutdown signal immediately and does not wait until
// registered goroutines have finished cleaning up.
func (n *Notifier) Notify() {
	n.cancel()
}

// NotifyAndWait sends the shutdown signal to all registered goroutines.
// Blocks until all listeners have called finish.
func (n *Notifier) NotifyAndWait() {
	n.cancel()
	n.wg.Wait()
}

// NotifyOnSignal waits until one of the given os signals are received before
// sending the shutdown signal, then blocks until all registered goroutines
// have finished cleaning up. Also returns if the shutdown signal was sent
// elsewhere.
func (n *Notifier) NotifyOnSignal(signals ...os.Signal) {
	defer func() {
		n.cancel()
		n.wg.Wait()
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, signals...)

	select {
	case <-sigchan:
		return

	case <-n.ctx.Done():
		return
	}
}
