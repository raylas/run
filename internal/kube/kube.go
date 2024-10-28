package kube

import (
	"context"
	"fmt"
	"time"

	"github.com/linecard/run/internal/output"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Run(ctx context.Context, attach, hostNetwork bool, secrets []string, cmd, name string) (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubernetes.config_path"))
	if err != nil {
		return "", err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	pc := c.CoreV1().Pods(viper.GetString("kubernetes.namespace"))

	nameRoot := "run-" + name
	timestamp := fmt.Sprint(time.Now().Unix())

	spec := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: nameRoot + "-" + timestamp[len(timestamp)-6:],
			Labels: map[string]string{
				"run.linecard.io/script": name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            nameRoot,
					Image:           viper.GetString("image"),
					ImagePullPolicy: corev1.PullIfNotPresent,
					EnvFrom:         envFrom(secrets),
					Command:         viper.GetStringSlice("entrypoint"),
					Args:            []string{cmd},
					Stdin:           attach,
					TTY:             attach,
				},
			},
			HostNetwork:   hostNetwork,
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	pod, err := pc.Create(ctx, &spec, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	switch {
	case attach:
		if err := tty(ctx, config, c, pod); err != nil {
			return "", err
		}
	default:
		return output.Colorize(output.Green, "Pod started: %s (%s)", pod.GetObjectMeta().GetName(), Context()), nil
	}

	return output.Colorize(output.Green, "Pod TTY closed"), nil
}

func Context() string {
	config, _ := clientcmd.LoadFromFile(viper.GetString("kubernetes.config_path"))
	return config.CurrentContext
}

func envFrom(secrets []string) []corev1.EnvFromSource {
	envFrom := []corev1.EnvFromSource{}
	for _, secret := range secrets {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secret,
				},
			},
		})
	}
	return envFrom
}
