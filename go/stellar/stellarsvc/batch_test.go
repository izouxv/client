package stellarsvc

import (
	"context"
	"testing"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/stellar1"
	"github.com/keybase/client/go/stellar"
	"github.com/keybase/stellarnet"
	"github.com/stretchr/testify/require"
)

// TestPrepareBatchRelays checks that a PrepareBatchPayments
// with a destination username that is a valid user but someone who
// doesn't have a wallet will succeed and create a relay payment.
func TestPrepareBatchRelays(t *testing.T) {
	tc, cleanup := setupDesktopTest(t)
	defer cleanup()
	require.NotNil(t, tc.Srv.walletState)

	tcw, cleanupw := setupDesktopTest(t)
	defer cleanupw()

	mctx := libkb.NewMetaContext(context.Background(), tc.G)

	acceptDisclaimer(tc)
	acceptDisclaimer(tcw)
	payments := []stellar1.BatchPaymentArg{
		{Recipient: "t_rebecca", Amount: "1"},
		{Recipient: tcw.Fu.Username, Amount: "2"},
	}

	_, senderAccountBundle, err := stellar.LookupSenderPrimary(mctx)
	require.NoError(t, err)
	senderSeed, err := stellarnet.NewSeedStr(senderAccountBundle.Signers[0].SecureNoLogString())
	require.NoError(t, err)

	prepared, err := stellar.PrepareBatchPayments(mctx, tc.Srv.walletState, senderSeed, payments)
	require.NoError(t, err)
	require.Len(t, prepared, 2)
	for i, p := range prepared {
		t.Logf("result %d: %+v", i, p)

		switch p.Username.String() {
		case "t_rebecca":
			require.Nil(t, p.Direct)
			require.NotNil(t, p.Relay)
			require.True(t, p.Relay.QuickReturn)
			require.Nil(t, p.Error)
			require.NotEmpty(t, p.Seqno)
			require.NotEmpty(t, p.TxID)
		case tcw.Fu.Username:
			require.NotNil(t, p.Direct)
			require.Nil(t, p.Relay)
			require.True(t, p.Direct.QuickReturn)
			require.Nil(t, p.Error)
			require.NotEmpty(t, p.Seqno)
			require.NotEmpty(t, p.TxID)
		default:
			t.Fatalf("unknown username in result: %s", p.Username)
		}
	}
	if prepared[0].Seqno > prepared[1].Seqno {
		t.Errorf("prepared sort failed (seqnos out of order)")
	}
}
