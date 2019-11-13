package specs

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var kubePod = &corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name: "test-pod",
	},
	TypeMeta: metav1.TypeMeta{
		Kind:       "Pod",
		APIVersion: "v1",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:    "container1",
				Image:   "busybox",
				Command: []string{"/bin/sleep"},
				Args:    []string{"100s"},
			},
		},
	},
}
