"""
AgenticSFU Python AI Voice Agent Worker
Connects to AgenticSFU Media Server via Model Context Protocol (MCP) and streams real-time audio.
"""

import json
import urllib.request
import time

SERVER_MCP_URL = "http://localhost:8080/mcp/v1"

def run_agent_worker():
    print("[AI AGENT WORKER] Initializing Multimodal Voice Agent Worker...")
    
    # 1. Discover MCP Server Tools
    try:
        req = urllib.request.Request(f"{SERVER_MCP_URL}/tools")
        with urllib.request.urlopen(req) as response:
            tools = json.loads(response.read().decode())
            print(f"[AI AGENT WORKER] Discovered {len(tools)} MCP Tools on AgenticSFU Server:")
            for t in tools:
                print(f"  - {t['name']}: {t['description']}")
    except Exception as e:
        print(f"[AI AGENT WORKER] Warning: Could not reach MCP Server ({e})")

    # 2. Dispatch Agent into Room
    try:
        data = json.dumps({"room": "demo-room", "prompt": "Act as a voice assistant"}).encode('utf-8')
        req = urllib.request.Request(f"{SERVER_MCP_URL}/dispatch", data=data, headers={'Content-Type': 'application/json'})
        with urllib.request.urlopen(req) as response:
            res = json.loads(response.read().decode())
            print(f"[AI AGENT WORKER] Successfully dispatched into room! Agent ID: {res.get('agent_id')}")
    except Exception as e:
        print(f"[AI AGENT WORKER] Dispatch failed: {e}")

if __name__ == "__main__":
    run_agent_worker()
