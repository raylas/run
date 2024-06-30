package kube

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func tty(ctx context.Context, c *rest.Config, cs *kubernetes.Clientset, pod *corev1.Pod) error {
	pc := cs.CoreV1().Pods(pod.Namespace)

	for {
		pod, err := pc.Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
	}

	req := cs.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("attach").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	executor, err := remotecommand.NewSPDYExecutor(c, "POST", req.URL())
	if err != nil {
		return err
	}

	prevTerm, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), prevTerm)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	resizeChan := make(chan resizeEvent, 1)
	sizeQueue := &terminalSizeQueue{
		resize: resizeChan,
	}

	go func() {
		for {
			if w, h, err := term.GetSize(0); err == nil {
				resizeChan <- resizeEvent{
					width:  uint16(w),
					height: uint16(h),
				}
			}
		}
	}()

	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: sizeQueue,
	})
	if err != nil {
		return err
	}

	return nil
}

type resizeEvent struct {
	width  uint16
	height uint16
}

func (r resizeEvent) Next() *remotecommand.TerminalSize {
	return &remotecommand.TerminalSize{
		Width:  r.width,
		Height: r.height,
	}
}

type terminalSizeQueue struct {
	resize chan resizeEvent
}

func (t *terminalSizeQueue) Next() *remotecommand.TerminalSize {
	event := <-t.resize
	return &remotecommand.TerminalSize{
		Width:  event.width,
		Height: event.height,
	}
}
