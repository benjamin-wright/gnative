package services

type serviceFile struct {
	file string
	data string
}

func getMongoOutputs() []serviceFile {
	outputs := []serviceFile{}

	outputs = append(outputs, getMongoDeployment())
	outputs = append(outputs, getMongoService())

	return outputs
}

func getMongoDeployment() serviceFile {
	return serviceFile{
		file: "mongo-deployment.yaml",
		data: `
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: mongo
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
        - containerPort: 27017`,
	}
}

func getMongoService() serviceFile {
	return serviceFile{
		file: "mongo-service.yaml",
		data: `
apiVersion: v1
kind: Service
metadata:
  name: mongo
  labels:
    name: mongo
spec:
  ports:
  - port: 27017
    targetPort: 27017
  selector:
    role: mongo`,
	}
}
