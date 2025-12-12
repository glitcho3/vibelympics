package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	ctx := context.Background()
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

        reportsPath := os.Getenv("OUTPUT_FOLDER")
        imagesFile := filepath.Join(reportsPath, "images.txt")
	namespace := "default"

	file, err := os.Open(imagesFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	jobs := []*batchv1.Job{}

	for scanner.Scan() {
		img := scanner.Text()
		if img == "" {
			continue
		}

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(img)))

		//sbomFile := filepath.Join(reportsPath, hash+".cdx.json")
		//vulnFile := filepath.Join(reportsPath, hash+".vulns.json")
		//provFile := filepath.Join(reportsPath, hash+".prov.json")
                promChart := os.Getenv("PROM_CHART")
		chartFolder := filepath.Join(reportsPath, promChart) // e.g., kube-prometheus-stack
                os.MkdirAll(reportsPath, 0755)
                
                sbomFile := filepath.Join(chartFolder, hash+".cdx.json")
                vulnFile := filepath.Join(chartFolder, hash+".vulns.json")
                provFile := filepath.Join(chartFolder, hash+".prov.json")

		jobName := "trivy-" + hash[:8]

		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: namespace,
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						InitContainers: []corev1.Container{
							{
								Name:    "trivy-sbom",
								Image:   "aquasec/trivy:0.68.1",
								Command: []string{"trivy", "image"},
								Args: []string{
									"--format", "cyclonedx",
									"--output", sbomFile,
									img,
								},
								VolumeMounts: []corev1.VolumeMount{
									{Name: "reports", MountPath: reportsPath},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:    "trivy",
								Image:   "aquasec/trivy:0.68.1",
								Command: []string{"trivy", "sbom"},
								Args: []string{
									"--format", "json",
									"--output", vulnFile,
									sbomFile,
								},
								VolumeMounts: []corev1.VolumeMount{
									{Name: "reports", MountPath: reportsPath},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "reports",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: "reports-pvc",
									},
								},
							},
						},
					},
				},
			},
		}

		// Provenance job

                // Provenance job
                provJob := &batchv1.Job{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      "prov-" + hash[:8],
                        Namespace: namespace,
                    },
                    Spec: batchv1.JobSpec{
                        Template: corev1.PodTemplateSpec{
                            Spec: corev1.PodSpec{
                                RestartPolicy: corev1.RestartPolicyNever,
                                Containers: []corev1.Container{
                                    {
                                        Name:            "provenor",
                                        Image:           "helm-auditor:latest",
                                        ImagePullPolicy: corev1.PullNever,
                                        Command:         []string{"/auditor-provenor"},
                                        Env: []corev1.EnvVar{
                                            {Name: "PROV_IMAGE", Value: img},
                                            {Name: "OUTPUT_FOLDER", Value: provFile},
                                        },
                                        VolumeMounts: []corev1.VolumeMount{
                                            {Name: "reports", MountPath: reportsPath},
                                        },
                                    },
                                },
                                Volumes: []corev1.Volume{
                                    {
                                        Name: "reports",
                                        VolumeSource: corev1.VolumeSource{
                                            PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                                                ClaimName: "reports-pvc",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                }
		_, err = clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create job for %s: %v\n", img, err)
		} else {
			fmt.Printf("Job %s dispatched for image %s\n", jobName, img)
			jobs = append(jobs, job)
		}

		_, err = clientset.BatchV1().Jobs(namespace).Create(ctx, provJob, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create provenance job for %s: %v\n", img, err)
		} else {
			fmt.Printf("Prov job prov-%s dispatched\n", hash[:8])
			jobs = append(jobs, provJob)
		}
	}

	fmt.Println("Waiting for all jobs to complete...")

	for _, job := range jobs {
		for {
			j, err := clientset.BatchV1().Jobs(namespace).Get(ctx, job.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Printf("Error fetching job %s: %v\n", job.Name, err)
				break
			}
			if j.Status.Succeeded > 0 {
				fmt.Printf("Job %s completed\n", job.Name)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
}

