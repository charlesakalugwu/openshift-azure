package cluster

//go:generate go get github.com/golang/mock/mockgen
//go:generate mockgen -destination=../util/mocks/mock_$GOPACKAGE/hash.go -package=mock_$GOPACKAGE -source hash.go
//go:generate gofmt -s -l -w ../util/mocks/mock_$GOPACKAGE/hash.go
//go:generate go get golang.org/x/tools/cmd/goimports
//go:generate goimports -local=github.com/openshift/openshift-azure -e -w ../util/mocks/mock_$GOPACKAGE/hash.go

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/arm"
	"github.com/openshift/openshift-azure/pkg/metrics"
	"github.com/openshift/openshift-azure/pkg/startup"
	"github.com/openshift/openshift-azure/pkg/sync"
)

type Hasher interface {
	HashScaleSet(*api.OpenShiftManagedCluster, *api.AgentPoolProfile) ([]byte, error)
	HashSyncPod(cs *api.OpenShiftManagedCluster) ([]byte, error)
}

type Hash struct {
	Log            *logrus.Entry
	TestConfig     api.TestConfig
	StartupFactory func(*logrus.Entry, *api.OpenShiftManagedCluster, api.TestConfig) (startup.Interface, error)
	Arm            arm.Interface
}

var _ Hasher = &Hash{}

// HashScaleSet returns the hash of a scale set
func (h *Hash) HashScaleSet(cs *api.OpenShiftManagedCluster, app *api.AgentPoolProfile) ([]byte, error) {
	hash := sha256.New()

	if armhd, ok := h.Arm.(interface {
		HashData(*api.AgentPoolProfile) ([]byte, error)
	}); ok {
		b, err := armhd.HashData(app) // legacy code path only for v3
		if err != nil {
			return nil, err
		}

		hash.Write(b)

	} else {
		b, err := h.Arm.Hash(app)
		if err != nil {
			return nil, err
		}

		hash.Write(b)
	}

	s, err := h.StartupFactory(h.Log, cs, h.TestConfig)
	if err != nil {
		return nil, err
	}

	if shd, ok := s.(interface {
		HashData(api.AgentPoolProfileRole) ([]byte, error)
	}); ok {
		b, err := shd.HashData(app.Role) // legacy code path only for v3
		if err != nil {
			return nil, err
		}

		hash.Write(b)

	} else {
		b, err := s.Hash(app.Role)
		if err != nil {
			return nil, err
		}

		hash.Write(b)
	}

	if app.Role == api.AgentPoolProfileRoleMaster {
		// add certificates pulled from keyvault by the master to the hash, to
		// ensure the masters update if a cert changes.  We don't add the keys
		// because these are not necessarily stable (sometimes the 'D' value of
		// the RSA key returned by keyvault differs to the one that was sent).
		// I believe that in a given RSA key, there are multiple suitable values
		// of 'D', so this is not a problem, however it doesn't make the value
		// suitable for a hash.  References:
		// https://stackoverflow.com/a/14233140,
		// https://crypto.stackexchange.com/a/46572.
		hash.Write(cs.Config.Certificates.OpenShiftConsole.Certs[0].Raw)
		hash.Write(cs.Config.Certificates.Router.Certs[0].Raw)

		if h.TestConfig.DebugHashFunctions {
			err = ioutil.WriteFile(fmt.Sprintf("cert-console-%d", time.Now().Unix()), cs.Config.Certificates.OpenShiftConsole.Certs[0].Raw, 0666)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(fmt.Sprintf("cert-router-%d", time.Now().Unix()), cs.Config.Certificates.Router.Certs[0].Raw, 0666)
			if err != nil {
				return nil, err
			}
		}
	}

	return hash.Sum(nil), nil
}

// HashSyncPod returns the hash of the sync pod output
func (h *Hash) HashSyncPod(cs *api.OpenShiftManagedCluster) ([]byte, error) {
	s, err := sync.New(h.Log, cs, false, metrics.DefaultCollector())
	if err != nil {
		return nil, err
	}

	return s.Hash()
}
