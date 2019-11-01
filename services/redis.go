package services

import (
	"bytes"
	"text/template"
)

func getRedisOutputs(namespace string) []serviceFile {
	outputs := []serviceFile{}

	outputs = append(outputs, getRedisDeployment(namespace))
	outputs = append(outputs, getRedisService(namespace))

	return outputs
}

func getRedisDeployment(namespace string) serviceFile {
	tmpl, err := template.New("redisDeployment").Parse(`
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: redis
  namespace: {{ . }}
spec:
  serviceName: "redis"
  replicas: 1
  template:
    metadata:
      labels:
        name: redis
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: redis
        image: redis
        ports:
        - containerPort: 6379
`)

  if err != nil {
		panic(err)
	}

	var data bytes.Buffer
	err = tmpl.Execute(&data, namespace)

	if err != nil {
		panic(err)
	}

	return serviceFile{
		file: "redis-deployment.yaml",
		data: data.String(),
	}
}

func getRedisService(namespace string) serviceFile {
  tmpl, err := template.New("redisDeployment").Parse(`
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: {{ . }}
  labels:
    name: redis
spec:
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    name: redis
  `)
  
    if err != nil {
      panic(err)
    }
  
    var data bytes.Buffer
    err = tmpl.Execute(&data, namespace)
  
    if err != nil {
      panic(err)
    }

	return serviceFile{
		file: "redis-service.yaml",
		data: data.String(),
	}
}
