# CloudPork CLI Agent

Cut the pork from your cloud costs with intelligent codebase analysis.

## 🚀 Quick Start (Free Trial)

1. **Install CloudPork:**
   ```bash
   curl -sSL https://cli.cloudpork.com/install.sh | bash
   ```

2. **Start your free trial:**
   ```bash
   cloudpork auth signup
   ```
   
3. **Analyze your codebase:**
   ```bash
   cloudpork analyze
   ```

**Trial includes:** 1 analysis, 7 days to try, no credit card required.

## 🎯 Trial Strategy

Your trial gives you **1 analysis** - make it count! Choose your most important project to see potential cost savings.

After your analysis, you'll see:
- Projected monthly costs at different scales
- Specific optimization recommendations  
- Architecture bottlenecks and solutions

## 💰 Upgrade Options

- **Starter ($29/mo):** 10 analyses/month + export + history
- **Professional ($149/mo):** 100 analyses/month + team + API  
- **Enterprise ($499/mo):** Unlimited + local AI + security

[View detailed pricing →](https://cloudpork.com/pricing)

## Features

- 🔒 **Privacy-first**: Code analysis happens locally, only summaries are sent
- 🧠 **Claude Code integration**: Leverages advanced AI for deep code understanding  
- 📊 **Cost projections**: Get scaling estimates from 1K to 100M+ users
- 🎯 **Optimization recommendations**: Identify bottlenecks before they become expensive
- 🚀 **Multiple installation methods**: curl, brew, package managers, or direct download
- 💻 **Local LLM support**: Analyze code completely offline with local models

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://cli.cloudpork.com/install.sh | bash
```

### Package Managers

```bash
# macOS
brew install carsor007/tap/cloudpork

# Ubuntu/Debian
curl -sL https://github.com/carsor007/cloudpork-agent/releases/latest/download/cloudpork_linux_amd64.deb
sudo dpkg -i cloudpork_linux_amd64.deb

# CentOS/RHEL
curl -sL https://github.com/carsor007/cloudpork-agent/releases/latest/download/cloudpork_linux_amd64.rpm  
sudo rpm -i cloudpork_linux_amd64.rpm
```

### Manual Download

Download the latest binary from [GitHub Releases](https://github.com/carsor007/cloudpork-agent/releases).

## Prerequisites

CloudPork Agent requires **Claude Code CLI** to function:

```bash
# macOS
brew install claude-ai/tap/claude

# Linux  
curl -fsSL https://claude.ai/install.sh | sh

# Windows
# Download from https://claude.ai/cli/download

# All platforms
# Visit https://claude.ai/cli for installation instructions
```

After installation, authenticate Claude Code:
```bash
claude auth login
```

## Usage

### 1. Start Your Free Trial

```bash
cloudpork auth signup
```

No credit card required - just enter your email and get started!

### 2. Analyze Your Project

```bash
# Analyze current directory
cloudpork analyze

# Analyze specific directory  
cloudpork analyze ./my-project

# Use specific project ID
cloudpork analyze --project-id=proj_abc123

# Output raw JSON
cloudpork analyze --output=json
```

### 3. Check Your Subscription

```bash
cloudpork auth status
```

### 4. Set Up Local AI (Optional)

For complete privacy and offline analysis:

```bash
# Set up local analysis mode
cloudpork setup --mode=local

# Or hybrid mode (local analysis + cloud intelligence)
cloudpork setup --mode=hybrid

# Check system health
cloudpork doctor
```

### 5. View Results

After analysis, view your cost projections and optimization recommendations at:
**https://cloudpork.com/dashboard**

## What Gets Analyzed

CloudPork Agent analyzes your codebase to understand:

- **Infrastructure patterns**: Database queries, API endpoints, background jobs
- **Resource requirements**: Memory, CPU, network, storage needs
- **Scaling bottlenecks**: N+1 queries, connection limits, memory leaks
- **Performance issues**: Large payloads, synchronous operations
- **Cost optimization opportunities**: Caching, connection pooling, architecture changes

## Example Output

```bash
🐷 CloudPork Agent - Cut the pork from your cloud costs
Analyzing your codebase for cost optimizations...

📁 Analyzing: /Users/dev/my-app
🆔 Project ID: proj_abc123

🔍 Running code analysis .... ✅

📊 Analysis Summary
==================================================
🏗️  Framework: Express.js
💾  Language: JavaScript  
📦  Dependencies: 42
🔌  API Endpoints: 15
⚡  Background Jobs: 3

💻 Resource Requirements
  Memory: 512 MB
  CPU: 1.0 cores
  DB Connections: 10
  Network: 50 Mbps

⚠️  Scaling Bottlenecks
  ● database: N+1 query pattern in user endpoints
  ● memory: Potential memory leak in image processing

🎯 Complexity Score: 65/100

📡 Sending results to CloudPork...
✅ Analysis complete!
🌐 View results: https://cloudpork.com/dashboard
```

## Commands

### `cloudpork analyze`
Analyze your codebase for cost optimization opportunities.

**Options:**
- `--project-id, -p`: Specify CloudPork project ID
- `--output, -o`: Output format (dashboard, json, quiet)

### `cloudpork auth`
Manage authentication with CloudPork.

**Subcommands:**
- `signup`: Start your free trial (no credit card required)
- `login`: Authenticate with API key
- `logout`: Remove stored credentials  
- `status`: Check subscription and authentication status

### `cloudpork setup`
Set up local LLM for private code analysis.

**Options:**
- `--mode`: Analysis mode (local, hybrid, cloud)
- `--model`: Specific model to install
- `--validate-hardware`: Check system compatibility
- `--force`: Force reinstall
- `--skip-validation`: Skip validation checks

### `cloudpork doctor`
Diagnose CloudPork setup and configuration issues.

### `cloudpork version`
Show version information.

## Configuration

CloudPork Agent stores configuration in `~/.cloudpork.yaml`:

```yaml
api_key: cp_your_api_key_here
project_id: proj_abc123
```

### Environment Variables

- `CLOUDPORK_API_KEY`: API key for authentication
- `CLOUDPORK_PROJECT_ID`: Default project ID
- `CLOUDPORK_VERBOSE`: Enable verbose output

## Privacy & Security

- **Your code never leaves your machine** - only analysis summaries are sent
- **API keys are stored securely** in your system keychain
- **Configuration files have restricted permissions** (600)
- **All API communication uses HTTPS** with certificate validation
- **Open source** - audit the code yourself

## Development

### Building from Source

```bash
git clone https://github.com/carsor007/cloudpork-agent.git
cd cloudpork-agent
make build
```

### Running Tests

```bash
make test
make test-coverage
```

### Development Mode

```bash
make dev
```

## Troubleshooting

### "Claude Code CLI not found"
Install Claude Code CLI first:
```bash
# macOS
brew install claude-ai/tap/claude

# Other platforms  
curl -fsSL https://claude.ai/install.sh | sh
```

### "No API key found"
Authenticate with CloudPork:
```bash
cloudpork auth login
```

### "Analysis failed"
Ensure you're in a code project directory with recognizable files (package.json, go.mod, requirements.txt, etc.).

### "Request failed"
Check your internet connection and API key validity:
```bash
cloudpork auth status
```

## Support

- 📖 **Documentation**: https://docs.cloudpork.com
- 🌐 **Dashboard**: https://cloudpork.com/dashboard  
- 🐛 **Issues**: https://github.com/carsor007/cloudpork-agent/issues
- 💬 **Support**: support@cloudpork.com

## License

MIT License - see [LICENSE](LICENSE) for details.

---

Made with 🐷 by the CloudPork team
