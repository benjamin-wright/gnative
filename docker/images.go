package docker

import (
	"bytes"
	"text/template"
)

const BASE_IMAGE_NAME = "ko.local/go-base:0.0.1"

type binding struct {
	ImportPath string
	ContainerPath string
}

func getBaseGoBuildDockerfile() string {
	return `
FROM golang:latest
COPY ./ /go/libraries/
WORKDIR /go/src/app
`
}

func getGoDockerfile(libdirs []string, registry string) string {
	libpaths := make([]binding, len(libdirs))
	for i, v := range libdirs {
			libpaths[i] = binding{
				ImportPath: registry + "/" + v,
				ContainerPath: v,
			}
	}

	tmpl, err := template.New("dockerfile").Parse(`
FROM ko.local/go-base:0.0.1 as builder
COPY ./go.* ./
RUN go mod download
COPY ./ ./
{{ range . }}
RUN echo "replace {{ .ImportPath }} => ../../libraries/{{ .ContainerPath }}" >> go.mod
{{ end }}
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app

FROM scratch
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
`)

  	if err != nil {
		panic(err)
	}

	var data bytes.Buffer
	err = tmpl.Execute(&data, libpaths)

	if err != nil {
		panic(err)
	}

	return data.String()
}