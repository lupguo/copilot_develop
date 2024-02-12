package intershim

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func LogAndWrapf(err error, format string, args ...interface{}) error {
	err = errors.Wrapf(err, format, args)
	log.Error(err)
	return err
}
