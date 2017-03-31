/*
Copyright 2017 Jean-Christophe Sirot.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	jobv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/rest"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Spawn)
	glog.Fatalf("Unable to start HTTP server: %v", http.ListenAndServe(":80", router))
	glog.Info("Job spawner is ready...")
}

/* Spawn a new job */
func Spawn(w http.ResponseWriter, r *http.Request) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to load client config: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}
	opts := v1.ListOptions{LabelSelector: "app=model-builder"}
	jobs, err := clientset.BatchV1().Jobs("").List(opts)
	if err != nil {
		glog.Fatalf("Failed to get kubernetes Jobs API: %v", err)
	}
	glog.Infof("There are %d jobs with label app=model-builder in the cluster", len(jobs.Items))

	if len(jobs.Items) >= 5 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Max number of concurrent running jobs reached. Please submit later.")
		return
	}

	jobCount := r.URL.Query().Get("count")
	if len(jobCount) == 0 {
		jobCount = "10"
	}
	jobWait := r.URL.Query().Get("sleep")
	if len(jobWait) == 0 {
		jobWait = "10"
	}

	glog.Info("Spawning a new job")

	jobName := fmt.Sprintf("model-builder-%d", rand.Int31n(100000))
	comp := int32(1)
	jobConf := &jobv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"app": "model-builder"}},
		Spec: jobv1.JobSpec{
			Completions: &comp,
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Name: jobName,
					Labels: map[string]string{
						"app": "model-builder"}},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{Name: "job",
							Image: "jcsirot/simple-job",
							Env: []v1.EnvVar{
								{Name: "JOB_COUNT", Value: jobCount},
								{Name: "JOB_WAIT", Value: jobWait}}}},
					RestartPolicy: v1.RestartPolicyNever}}}}
	_, err2 := clientset.BatchV1().Jobs("default").Create(jobConf)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An error occurred: %s\n", err2)
		glog.Errorf("Failed to create new job: %v", err2)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Spawning job %s\n", jobName)
}
