package ingestion

import "embed"

//go:embed scripts/*
var Scripts embed.FS

// pathConf maps connection types to their ingestion script paths
// Add new connection types here as they are implemented
var pathConf = map[string]string{
	"mqtt":   "scripts/mqtt.py",
	"plc":    "scripts/plc.py",
	"modbus": "scripts/modbus.py",
	"opcua":  "scripts/opcua.py",
}
