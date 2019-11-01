package services

import (
	"bytes"
	"text/template"
)

type serviceFile struct {
	file string
	data string
}

func getMongoOutputs(namespace string) []serviceFile {
	outputs := []serviceFile{}

	outputs = append(outputs, getMongoDeployment(namespace))
	outputs = append(outputs, getMongoService(namespace))

	return outputs
}

func getMongoDeployment(namespace string) serviceFile {
	tmpl, err := template.New("mongoDeployment").Parse(`
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: mongo
  namespace: {{ . }}
spec:
  serviceName: "mongo"
  replicas: 1
  template:
    metadata:
      labels:
        name: mongo
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: mongo
        image: mongo
        ports:
        - containerPort: 27017
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
		file: "mongo-deployment.yaml",
		data: data.String(),
	}
}

func getMongoService(namespace string) serviceFile {
  tmpl, err := template.New("mongoDeployment").Parse(`
apiVersion: v1
kind: Service
metadata:
  name: mongo
  namespace: {{ . }}
  labels:
    name: mongo
spec:
  ports:
  - port: 27017
    targetPort: 27017
  selector:
    name: mongo
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
		file: "mongo-service.yaml",
		data: data.String(),
	}
}
