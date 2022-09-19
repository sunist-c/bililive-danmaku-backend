package backend

import (
	"os"

	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
)

var (
	exitService *ExitService
)

func init() {
	exitService = &ExitService{
		channelService: GetDefaultChannelService(),
	}

	exitService.Serve()
}

func GetDefaultExitService() *ExitService {
	return exitService
}

type ExitService struct {
	channelService *ChannelService
}

func (s *ExitService) Exit() {
	logging.Info("execute exit operation")
	for i := 0; i < 3; i++ {
		success, total := s.channelService.RemoveAllChannel()
		if total == success {
			return
		}
		logging.Warn("failed to remove all channels, closed %v, total %v, retrying...", success, total)
	}
	logging.Info("closed all channels, have a good day!")
	os.Exit(0)
}

func (s *ExitService) Serve() {

}
