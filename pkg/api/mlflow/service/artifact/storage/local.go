package storage

import (
	"io"
	"net/url"
	"os"

	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/config"
)

// Local represents local file storage adapter to work with artifacts.
type Local struct {
	config *config.ServiceConfig
}

// NewLocal creates new Local storage instance.
func NewLocal(config *config.ServiceConfig) (*Local, error) {
	return &Local{
		config: config,
	}, nil
}

// List implements Provider interface.
func (s Local) List(runArtifactPath, itemPath string) (string, []ArtifactObject, error) {
	path, err := url.JoinPath(s.config.ArtifactRoot, runArtifactPath, itemPath)
	if err != nil {
		return "", nil, eris.Wrap(err, "error constructing full path")
	}

	// test local storage existence
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path does not (yet) exist -- returning no artifacts
		log.Infof("artifact dir does not exist: %s", path)
		return path, []ArtifactObject{}, nil
	}

	// read objects
	objects, err := os.ReadDir(path)
	if err != nil {
		return "", nil, eris.Wrapf(err, "error reading object from local storage")
	}

	log.Debugf("got %d objects from local storage for path: %s", len(objects), path)
	artifactList := make([]ArtifactObject, len(objects))
	for i, object := range objects {
		info, err := object.Info()
		if err != nil {
			return "", nil, eris.Wrapf(err, "error getting info for object: %s", object.Name())
		}
		artifactList[i] = ArtifactObject{
			Path:  info.Name(),
			Size:  info.Size(),
			IsDir: object.IsDir(),
		}
	}
	return "/artifacts" + runArtifactPath, artifactList, nil
}

// GetArtifact will return actual item URI in the storage location
func (s Local) GetArtifact(runArtifactPath, itemPath string) (io.Reader, error) {
	path, err := url.JoinPath(s.config.ArtifactRoot, runArtifactPath, itemPath)
	if err != nil {
		return nil, eris.Wrap(err, "error constructing full path")
	}
	return os.Open(path)
}
