package commands

import (
	"context"
	"strings"

	"github.com/keybase/client/go/chat/globals"
	"github.com/keybase/client/go/protocol/chat1"
	"github.com/keybase/client/go/protocol/gregor1"
)

type Hide struct {
	*baseCommand
}

func NewHide(g *globals.Context) *Hide {
	return &Hide{
		baseCommand: newBaseCommand(g, "hide", "[ conversation to hide ]"),
	}
}

func (h *Hide) Execute(ctx context.Context, uid gregor1.UID, inConvID chat1.ConversationID,
	tlfName, text string) (err error) {
	defer h.Trace(ctx, func() error { return err }, "Execute")()
	if !h.Match(ctx, text) {
		return ErrInvalidCommand
	}
	var convID chat1.ConversationID
	toks := strings.Split(text, " ")
	if len(toks) == 1 {
		convID = inConvID

	} else if len(toks) == 2 {
		conv, err := h.getConvByName(ctx, uid, toks[1])
		if err != nil {
			return err
		}
		convID = conv.GetConvID()
	} else {
		return ErrInvalidArguments
	}
	if err = h.G().InboxSource.RemoteSetConversationStatus(ctx, uid, convID,
		chat1.ConversationStatus_IGNORED); err != nil {
		return err
	}
	return nil
}
