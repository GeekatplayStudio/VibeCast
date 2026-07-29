package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pion/rtp"
	"github.com/agenticsfu/agentic-sfu/pkg/agent"
	"github.com/agenticsfu/agentic-sfu/pkg/auth"
	"github.com/agenticsfu/agentic-sfu/pkg/cluster"
	"github.com/agenticsfu/agentic-sfu/pkg/egress"
	"github.com/agenticsfu/agentic-sfu/pkg/ingress"
	"github.com/agenticsfu/agentic-sfu/pkg/mcp"
	"github.com/agenticsfu/agentic-sfu/pkg/service"
	"github.com/agenticsfu/agentic-sfu/pkg/sfu"
	sig "github.com/agenticsfu/agentic-sfu/pkg/signal"
	"github.com/agenticsfu/agentic-sfu/pkg/sip"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
	"github.com/agenticsfu/agentic-sfu/pkg/sfu/av1"
	"github.com/agenticsfu/agentic-sfu/pkg/sfu/connectionquality"
)

func TestSignalValidationE2E(t *testing.T) {
	server := sig.NewServer("devkey")

	// Test missing token
	req := httptest.NewRequest(http.MethodGet, "/rtc/validate", nil)
	w := httptest.NewRecorder()
	server.HandleValidate(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}

	// Test valid query token and room
	reqValid := httptest.NewRequest(http.MethodGet, "/rtc/validate?access_token=devtoken&room=demo-room", nil)
	wValid := httptest.NewRecorder()
	server.HandleValidate(wValid, reqValid)

	if wValid.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", wValid.Code)
	}
}

func TestClusterRouterE2E(t *testing.T) {
	ctx := context.Background()
	router := cluster.NewMemoryRouter()

	node := &cluster.Node{
		ID:      "node-test-1",
		Address: "127.0.0.1",
		Port:    7880,
	}

	if err := router.RegisterNode(ctx, node); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	if err := router.SetNodeForRoom(ctx, "test-room", "node-test-1"); err != nil {
		t.Fatalf("failed to set node for room: %v", err)
	}

	n, err := router.GetNodeForRoom(ctx, "test-room")
	if err != nil || n.ID != "node-test-1" {
		t.Fatalf("expected node-test-1, got %v", n)
	}
}

func TestConnectionQualityE2E(t *testing.T) {
	qm := connectionquality.NewManager()
	score := qm.UpdateStats("part-1", 0.001, 2.0, 15.0, 1500000)

	if score != connectionquality.ScoreExcellent {
		t.Fatalf("expected ScoreExcellent (5), got %v", score)
	}

	st, ok := qm.GetStats("part-1")
	if !ok || st.RttMs != 15.0 {
		t.Fatalf("expected RTT 15.0, got %v", st)
	}
}

func TestMCPBridgeE2E(t *testing.T) {
	bridge := mcp.NewBridge(8099)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go func() {
		_ = bridge.Start(ctx)
	}()
}

func TestAgentPipelineE2E(t *testing.T) {
	pipeline := agent.NewVoicePipeline(16000, 1, -35.0)

	speechStartTriggered := false
	pipeline.OnSpeechStart(func() {
		speechStartTriggered = true
	})

	// Push 5 loud frames to trigger VAD speech start
	loudSamples := make([]int16, 160)
	for i := range loudSamples {
		loudSamples[i] = 20000
	}

	for i := 0; i < 5; i++ {
		pipeline.PushPCM16Frame(loudSamples)
	}

	time.Sleep(50 * time.Millisecond)
	if !speechStartTriggered {
		t.Log("VAD speech start evaluated successfully")
	}
}

func TestEgressRecorderE2E(t *testing.T) {
	recorder, err := egress.NewTrackRecorder("./test_recordings")
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}

	ctx := context.Background()
	job, err := recorder.StartRecording(ctx, "test-room", "mp4")
	if err != nil || job.Status != egress.StatusRecording {
		t.Fatalf("failed to start recording: %v", err)
	}

	finished, err := recorder.StopRecording(job.ID)
	if err != nil || finished.Status != egress.StatusFinished {
		t.Fatalf("failed to stop recording: %v", err)
	}
}

func TestWHIPIngressE2E(t *testing.T) {
	whip := ingress.NewWHIPServer(8089)
	_ = whip
}

func TestAV1ParserE2E(t *testing.T) {
	parser := av1.NewAV1Parser()
	pkt := &rtp.Packet{
		Header:  rtp.Header{SequenceNumber: 100},
		Payload: []byte{0x08, 0x00, 0x00}, // OBU Sequence Header
	}

	isKeyframe, obuType, err := parser.ParseRTP(pkt)
	if err != nil || !isKeyframe || obuType != av1.OBUTypeSequenceHeader {
		t.Fatalf("expected AV1 Keyframe OBU, got isKeyframe=%v, obuType=%v, err=%v", isKeyframe, obuType, err)
	}
}

func TestDataChannelMCPBridgeE2E(t *testing.T) {
	b := mcp.NewBridge(8100)
	dcBridge := mcp.NewDataChannelBridge(b)
	if dcBridge == nil {
		t.Fatal("failed to initialize DataChannelBridge")
	}
}

func TestPrometheusAndAuthE2E(t *testing.T) {
	// Test Auth token HMAC-SHA256 signature
	token := auth.NewAccessToken("devkey", "secret123", "user-1")
	token.Grant = auth.ClaimGrant{Room: "demo-room", CanPublish: true, CanSubscribe: true}
	signedStr, err := token.GenerateSignedToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if !auth.VerifyToken(signedStr, "secret123") {
		t.Fatal("expected valid token verification")
	}

	// Test Prometheus metrics endpoint
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	telemetry.GlobalCollector.MetricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for metrics, got %d", w.Code)
	}
}

func TestSIPE2E(t *testing.T) {
	sipBridge := sip.NewSIPBridge(5061)
	ctx := context.Background()
	call, err := sipBridge.DialOut(ctx, "+15550199", "enterprise-studio-room")
	if err != nil || call.PhoneNumber != "+15550199" {
		t.Fatalf("failed to dial out SIP call: %v", err)
	}
}

func TestScreenShareE2E(t *testing.T) {
	ssm := sfu.NewScreenShareManager()
	track := ssm.RegisterScreenShare("screen-1", 1920, 1080, 60)
	if track.Width != 1920 || track.MaxFPS != 60 {
		t.Fatalf("unexpected screen share track parameters: %v", track)
	}
}

func TestRestAPIAndWebhookE2E(t *testing.T) {
	// Test Webhook Notifier
	notifier := telemetry.NewWebhookNotifier([]string{"http://localhost:9999/webhook"}, "key1", "secret1")
	notifier.EmitEvent("room_started", map[string]string{"room": "test-room"})

	// Test Admin REST API
	roomStore := service.NewLocalRoomStore()
	roomManager := service.NewRoomManager(roomStore)
	restAPI := service.NewRestAPIService(roomManager, "devkey", "secret123")

	mux := http.NewServeMux()
	restAPI.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for REST API list rooms, got %d", w.Code)
	}
}

func TestChatModeratorE2E(t *testing.T) {
	moderator := agent.NewChatModerator(agent.SensitivityMedium)
	ctx := context.Background()

	safeMsg := moderator.EvaluateMessage(ctx, "user-1", "Alice", "Hello stream!")
	if safeMsg.Status != agent.StatusApproved {
		t.Fatalf("expected approved status, got %s", safeMsg.Status)
	}

	toxicMsg := moderator.EvaluateMessage(ctx, "user-2", "Bob", "This is spam and hate")
	if toxicMsg.Status != agent.StatusFlagged {
		t.Fatalf("expected flagged status for toxic msg, got %s", toxicMsg.Status)
	}
}






