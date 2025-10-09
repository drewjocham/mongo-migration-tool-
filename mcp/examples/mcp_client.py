"""
 MCP client for mongo-essential migration tool.

Usage:
    python mcp_client.py [command] [args...]

Commands:
    init            - Initialize the MCP server
    tools           - List available tools
    status          - Get migration status
    up [version]    - Apply migrations (optionally up to version)
    down [version]  - Roll back migrations (optionally to version)
    create <name> <description> - Create a new migration
    list            - List all registered migrations
    interactive     - Interactive mode with command prompt
"""

import json
import subprocess
import sys
import argparse
import os
from typing import Optional, Dict, Any

class MCPClient:
    def __init__(self, mcp_binary_path: str = "./build/mongo-essential"):
        self.mcp_binary_path = mcp_binary_path
        self.request_id = 0

    def _next_id(self) -> int:
        """Get next request ID."""
        self.request_id += 1
        return self.request_id

    def _run_mcp_request(self, method: str, params: Optional[Dict[str, Any]] = None) -> Optional[Dict[str, Any]]:
        """Run an MCP request and return the response."""
        request = {
            "jsonrpc": "2.0",
            "id": self._next_id(),
            "method": method,
            "params": params or {}
        }
        
        try:
            process = subprocess.Popen(
                [self.mcp_binary_path, "mcp", "--with-examples"],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            # Send request and get response
            stdout, stderr = process.communicate(json.dumps(request) + '\n')
            
            if stderr.strip():
                print(f"MCP Server Error: {stderr.strip()}", file=sys.stderr)
            
            if not stdout.strip():
                print("No response from MCP server", file=sys.stderr)
                return None
            
            try:
                return json.loads(stdout.strip())
            except json.JSONDecodeError as e:
                print(f"Failed to parse JSON response: {e}", file=sys.stderr)
                print(f"Raw response: {stdout}", file=sys.stderr)
                return None
                
        except FileNotFoundError:
            print(f"MCP binary not found at {self.mcp_binary_path}", file=sys.stderr)
            print("Build the binary with: make build", file=sys.stderr)
            return None
        except Exception as e:
            print(f"Error running MCP request: {e}", file=sys.stderr)
            return None

    def initialize(self) -> bool:
        """Initialize the MCP server."""
        print("Initializing MCP server...")
        response = self._run_mcp_request("initialize")
        
        if response and 'result' in response:
            server_info = response['result'].get('serverInfo', {})
            print(f"✅ Connected to {server_info.get('name', 'mongo-essential')} v{server_info.get('version', 'unknown')}")
            return True
        else:
            print("❌ Failed to initialize MCP server")
            return False

    def list_tools(self) -> bool:
        """List available tools."""
        print("Listing available tools...")
        response = self._run_mcp_request("tools/list")
        
        if response and 'result' in response:
            tools = response['result'].get('tools', [])
            print(f"\n📋 Available Tools ({len(tools)}):")
            print("=" * 50)
            for tool in tools:
                print(f"🔧 {tool['name']}")
                print(f"   {tool['description']}")
                print()
            return True
        else:
            print("❌ Failed to list tools")
            return False

    def migration_status(self) -> bool:
        """Get migration status."""
        print("Getting migration status...")
        response = self._run_mcp_request("tools/call", {
            "name": "migration_status",
            "arguments": {}
        })
        
        if response and 'result' in response:
            content = response['result'].get('content', [])
            if content and len(content) > 0:
                print("\n" + content[0].get('text', 'No status available'))
                return True
        
        print("❌ Failed to get migration status")
        return False

    def migration_up(self, version: Optional[str] = None) -> bool:
        """Apply migrations."""
        args = {"version": version} if version else {}
        action = f"up to version {version}" if version else "up (all pending)"
        
        print(f"Running migrations {action}...")
        response = self._run_mcp_request("tools/call", {
            "name": "migration_up",
            "arguments": args
        })
        
        if response and 'result' in response:
            content = response['result'].get('content', [])
            if content and len(content) > 0:
                print("\n" + content[0].get('text', 'Migration completed'))
                return True
        
        print("❌ Failed to run migrations")
        return False

    def migration_down(self, version: Optional[str] = None) -> bool:
        """Roll back migrations."""
        args = {"version": version} if version else {}
        action = f"down to version {version}" if version else "down (last migration)"
        
        print(f"Rolling back migrations {action}...")
        response = self._run_mcp_request("tools/call", {
            "name": "migration_down",
            "arguments": args
        })
        
        if response and 'result' in response:
            content = response['result'].get('content', [])
            if content and len(content) > 0:
                print("\n" + content[0].get('text', 'Migration rolled back'))
                return True
        
        print("❌ Failed to roll back migrations")
        return False

    def create_migration(self, name: str, description: str) -> bool:
        """Create a new migration."""
        print(f"Creating migration: {name}")
        response = self._run_mcp_request("tools/call", {
            "name": "migration_create",
            "arguments": {
                "name": name,
                "description": description
            }
        })
        
        if response and 'result' in response:
            content = response['result'].get('content', [])
            if content and len(content) > 0:
                print("\n" + content[0].get('text', 'Migration created'))
                return True
        
        print("❌ Failed to create migration")
        return False

    def list_migrations(self) -> bool:
        """List all registered migrations."""
        print("Listing registered migrations...")
        response = self._run_mcp_request("tools/call", {
            "name": "migration_list",
            "arguments": {}
        })
        
        if response and 'result' in response:
            content = response['result'].get('content', [])
            if content and len(content) > 0:
                print("\n" + content[0].get('text', 'No migrations found'))
                return True
        
        print("❌ Failed to list migrations")
        return False

    def interactive_mode(self):
        """Run in interactive mode."""
        print("🚀 MongoDB Migration MCP Client - Interactive Mode")
        print("Type 'help' for available commands, 'quit' to exit")
        print()
        
        # Initialize server once
        if not self.initialize():
            return
        
        while True:
            try:
                command = input("mcp> ").strip()
                
                if command in ['quit', 'exit', 'q']:
                    print("Goodbye! 👋")
                    break
                elif command in ['help', 'h']:
                    self._print_help()
                elif command == 'tools':
                    self.list_tools()
                elif command == 'status':
                    self.migration_status()
                elif command == 'up':
                    self.migration_up()
                elif command.startswith('up '):
                    version = command.split(' ', 1)[1]
                    self.migration_up(version)
                elif command == 'down':
                    self.migration_down()
                elif command.startswith('down '):
                    version = command.split(' ', 1)[1]
                    self.migration_down(version)
                elif command.startswith('create '):
                    parts = command.split(' ', 2)
                    if len(parts) < 3:
                        print("Usage: create <name> <description>")
                    else:
                        name_desc = parts[2].split(' ', 1)
                        if len(name_desc) < 2:
                            print("Usage: create <name> <description>")
                        else:
                            self.create_migration(name_desc[0], name_desc[1])
                elif command == 'list':
                    self.list_migrations()
                elif command == '':
                    continue
                else:
                    print(f"Unknown command: {command}. Type 'help' for available commands.")
                    
                print()  # Add spacing between commands
                
            except KeyboardInterrupt:
                print("\nGoodbye! 👋")
                break
            except EOFError:
                print("\nGoodbye! 👋")
                break

    def _print_help(self):
        """Print help information."""
        print("""
Available commands:
  help                     - Show this help
  tools                    - List available MCP tools
  status                   - Get migration status
  up [version]             - Apply migrations (optionally up to version)
  down [version]           - Roll back migrations (optionally to version)
  create <name> <desc>     - Create a new migration
  list                     - List all registered migrations
  quit/exit/q              - Exit interactive mode

Examples:
  status                   - Check what migrations are applied
  up                       - Apply all pending migrations
  up 20240101_001          - Apply migrations up to specific version
  down                     - Roll back the last applied migration
  create add_index "Add user email index" - Create a new migration
""")

def main():
    parser = argparse.ArgumentParser(description="MCP client for mongo-essential")
    parser.add_argument('--binary', default='./build/mongo-essential',
                        help='Path to mongo-essential binary')
    parser.add_argument('command', nargs='?',
                        choices=['init', 'tools', 'status', 'up', 'down', 'create', 'list', 'interactive'],
                        help='Command to run')
    parser.add_argument('args', nargs='*', help='Command arguments')

    args = parser.parse_args()

    # Check if binary exists
    if not os.path.exists(args.binary):
        print(f"Error: Binary not found at {args.binary}")
        print("Build it with: make build")
        sys.exit(1)

    client = MCPClient(args.binary)

    if not args.command or args.command == 'interactive':
        client.interactive_mode()
        return

    success = False

    needs_init = args.command not in ['init']
    if needs_init and not client.initialize():
        sys.exit(1)

    if args.command == 'init':
        success = client.initialize()
    elif args.command == 'tools':
        success = client.list_tools()
    elif args.command == 'status':
        success = client.migration_status()
    elif args.command == 'up':
        version = args.args[0] if args.args else None
        success = client.migration_up(version)
    elif args.command == 'down':
        version = args.args[0] if args.args else None
        success = client.migration_down(version)
    elif args.command == 'create':
        if len(args.args) < 2:
            print("Usage: create <name> <description>")
            sys.exit(1)
        success = client.create_migration(args.args[0], ' '.join(args.args[1:]))
    elif args.command == 'list':
        success = client.list_migrations()

    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()

if __name__ == "__main__":
    main()
