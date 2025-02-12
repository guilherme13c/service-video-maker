package videomaker

import (
	"context"
	"fmt"
	"os"
)

func (r *serviceVideoMaker) CleanUp(ctx context.Context, requestId string) error {
	return os.RemoveAll(fmt.Sprintf("tmp/%s", requestId))
}
