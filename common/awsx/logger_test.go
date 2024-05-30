package awsx

import (
	"testing"

	"github.com/aws/smithy-go/logging"
	"github.com/monorepo/common/logging/loggingtest"
	"github.com/stretchr/testify/mock"
)

func Test_Logf(t *testing.T) {

	t.Run("testwarn", func(t *testing.T) {
		mocker := loggingtest.NewMock(t)
		logger := AwsLoggerfromLogger(mocker)

		//given a warn classification
		mocker.On("Warningf", mock.Anything, mock.Anything).Once()

		//when logging
		logger.Logf(logging.Warn, "a format : %v", "value")

		mocker.AssertExpectations(t)
	})

	t.Run("testdebug", func(t *testing.T) {
		mocker := loggingtest.NewMock(t)
		logger := AwsLoggerfromLogger(mocker)

		//given a info classification
		mocker.On("Debugf", mock.Anything, mock.Anything).Once()

		//when logging
		logger.Logf(logging.Debug, "a format : %v", "value")

		mocker.AssertExpectations(t)
	})

	t.Run("testother", func(t *testing.T) {
		mocker := loggingtest.NewMock(t)
		logger := AwsLoggerfromLogger(mocker)

		//given an uncharted classification
		mocker.On("Errorf", mock.Anything, mock.Anything).Once()

		//when logging
		logger.Logf("something else", "a format : %v", "value")
		//Calls error
		mocker.AssertExpectations(t)
	})
}
