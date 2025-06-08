# Makefile for notiboard GUI application

.PHONY: debug-gui kill-debug help

# Variables
DEBUG_PORT=2345

# Debug GUI application with better settings for graphics apps
debug-gui:
	@echo "Starting notiboard GUI with dlv debugger on port $(DEBUG_PORT)..."
	@echo "This target is optimized for GUI applications"
	@echo "Connect with: dlv connect localhost:$(DEBUG_PORT)"
	DISPLAY=${DISPLAY} dlv debug --headless --listen=:$(DEBUG_PORT) --api-version=2 --check-go-version=false main.go

# Kill any running dlv processes
kill-debug:
	@echo "Killing any running dlv processes..."
	-pkill -f "dlv debug"
	-pkill -f "dlv exec"
	-pkill -f "dlv connect"

# Show help
help:
	@echo "Available targets:"
	@echo "  debug-gui   - Run GUI app with dlv (optimized for graphics)"
	@echo "  kill-debug  - Kill any running dlv processes"
	@echo "  help        - Show this help message"
	@echo ""
	@echo "Debug connection:"
	@echo "  Connect your debugger client to localhost:$(DEBUG_PORT)"
	@echo "  Example: dlv connect localhost:$(DEBUG_PORT)"