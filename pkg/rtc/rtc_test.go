package rtc_test

import (
	"context"
	"testing"

	"github.com/agenticsfu/agentic-sfu/pkg/rtc"
	"github.com/pion/webrtc/v3"
)

func TestRTCModule(t *testing.T) {
	ctx := context.Background()
	room := rtc.NewRoom("test-room", "sid-12345")

	if room.Name() != "test-room" || room.SID() != "sid-12345" {
		t.Fatalf("unexpected room name or sid")
	}

	p := rtc.NewParticipant("user-1", "Alice", rtc.DefaultPermissions())
	if err := room.JoinParticipant(p); err != nil {
		t.Fatalf("failed to join participant: %v", err)
	}

	if room.ParticipantCount() != 1 {
		t.Fatalf("expected 1 participant, got %d", room.ParticipantCount())
	}

	// MediaTrack test
	track := rtc.NewMediaTrack("track-v1", webrtc.RTPCodecTypeVideo, "user-1")
	p.AddTrack(track)

	if len(p.PublishedTracks()) != 1 {
		t.Fatalf("expected 1 published track")
	}

	// Leave participant
	if err := room.LeaveParticipant(ctx, "user-1"); err != nil {
		t.Fatalf("failed to leave participant: %v", err)
	}

	if room.ParticipantCount() != 0 {
		t.Fatalf("expected 0 participants after leave")
	}
}
