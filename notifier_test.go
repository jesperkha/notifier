package notifier_test

import (
	"testing"
	"time"

	"github.com/jesperkha/notifier"
)

func TestRegisterAndNotify(t *testing.T) {
	n := notifier.New()
	done, finish := n.Register()

	pass := false

	go func() {
		<-done
		pass = true
		finish()
	}()

	n.NotifyAndWait()
	if !pass {
		t.Fail()
	}
}

func TestNoListeners(t *testing.T) {
	n := notifier.New()
	n.NotifyAndWait()
}

func TestCancelTwice(t *testing.T) {
	n := notifier.New()
	n.Notify()
	n.Notify()
}

func TestNotifierWithDelayedShutdown(t *testing.T) {
	n := notifier.New()
	done, finish := n.Register()
	pass := false

	go func() {
		<-done
		time.Sleep(500 * time.Millisecond)
		pass = true
		finish()
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		n.Notify()
	}()

	n.Wait()
	if !pass {
		t.Fail()
	}
}

func TestNotifier_EarlyContextCancellation(t *testing.T) {
	n := notifier.New()
	done, finish := n.Register()

	n.Notify()

	select {
	case <-done:
		finish()
	case <-time.After(time.Second):
		t.Fatal("expected done channel to close")
	}
}
