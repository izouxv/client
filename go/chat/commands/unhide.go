package commands

import (
	"context"
	"strings"

	"github.com/keybase/client/go/chat/globals"
	"github.com/keybase/client/go/protocol/chat1"
	"github.com/keybase/client/go/protocol/gregor1"
)

type Unhide struct {
	*baseCommand
}

func NewUnhide(g *globals.Context) *Unhide {
	return &Unhide{
		baseCommand: newBaseCommand(g, "unhide", "<conversation to unhide>"),
	}
}

func (h *Unhide) Execute(ctx context.Context, uid gregor1.UID, _ chat1.ConversationID,
	tlfName, text string) (err error) {
	defer h.Trace(ctx, func() error { return err }, "Execute")()
	if !h.Match(ctx, text) {
		return ErrInvalidCommand
	}
	toks := strings.Split(text, " ")
	if len(toks) < 2 {
		return ErrInvalidArguments
	}
	conv, err := h.getConvByName(ctx, uid, toks[1])
	if err != nil {
		return err
	}
	if err = h.G().InboxSource.RemoteSetConversationStatus(ctx, uid, conv.GetConvID(),
		chat1.ConversationStatus_UNFILED); err != nil {
		return err
	}
	return nil
}
