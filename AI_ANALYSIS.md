# AI-Powered Database Analysis

The MongoDB migration tool now includes AI-powered database analysis capabilities using multiple AI providers.

## Features

### 🤖 Multi-Provider AI Support
- **OpenAI**: GPT-4o, GPT-4o-mini, GPT-3.5-turbo
- **Google Gemini**: Gemini-1.5-flash, Gemini-1.5-pro  
- **Claude**: Coming soon (temporarily disabled)

### 🔍 Analysis Types

#### 1. Comprehensive Database Analysis
```bash
mongo-migrate ai analyze --provider openai
mongo-migrate ai analyze --provider gemini --format json
```

**Provides insights on:**
- Overall database health assessment
- Performance optimization recommendations
- Index optimization suggestions
- Schema design improvements
- Security considerations
- Best practices compliance
- Specific action items

#### 2. Schema Analysis
```bash
mongo-migrate ai schema --provider openai
mongo-migrate ai schema users --detail --provider gemini
```

**Analyzes:**
- Collection schema structures
- Data type optimization
- Normalization/denormalization recommendations
- Field naming and structure improvements
- Validation rule suggestions
- Migration strategies

#### 3. Performance Analysis
```bash
mongo-migrate ai performance --provider openai
```

**Examines:**
- Performance bottleneck identification
- Index utilization analysis
- Query optimization suggestions
- Resource usage patterns
- Scaling recommendations
- Monitoring suggestions

#### 4. Oplog & Replication Analysis
```bash
mongo-migrate ai oplog --provider gemini --detail
```

**Analyzes:**
- Oplog health and sizing
- Replication lag analysis
- Throughput optimization
- Operation pattern insights
- Replica set health assessment
- Replication monitoring recommendations

#### 5. Change Stream Analysis
```bash
mongo-migrate ai changestream --provider openai
mongo-migrate ai changestream --collection users --provider gemini
```

**Reviews:**
- Change stream configuration
- Resume token management
- Performance optimization
- Real-time processing patterns
- Error handling strategies
- Scaling recommendations

## Configuration

### Environment Variables

Add these to your `.env` file:

```bash
# Enable AI Analysis
AI_ENABLED=true
AI_PROVIDER=openai  # openai, gemini, claude

# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-4o-mini

# Google Gemini Configuration  
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-1.5-flash

# Claude (Coming Soon)
CLAUDE_API_KEY=your_claude_api_key_here
CLAUDE_MODEL=claude-3-5-sonnet-20241022
```

### API Key Setup

#### OpenAI
1. Go to https://platform.openai.com/api-keys
2. Create a new API key
3. Set `OPENAI_API_KEY` in your environment

#### Google Gemini
1. Go to https://ai.google.dev/
2. Get your API key from Google AI Studio
3. Set `GEMINI_API_KEY` in your environment

#### Google Docs Integration
1. Create a Google Cloud Project
2. Enable Google Docs and Drive APIs
3. Create a service account and download the JSON credentials
4. Set `GOOGLE_CREDENTIALS_PATH` to the JSON file path
5. Optionally set `GOOGLE_DRIVE_FOLDER_ID` and `GOOGLE_DOCS_SHARE_WITH_EMAIL`

## Usage Examples

### Quick Database Health Check
```bash
# Basic analysis with OpenAI
mongo-migrate ai analyze --provider openai

# Detailed analysis with Gemini
mongo-migrate ai analyze --provider gemini --detail

# Export results as JSON
mongo-migrate ai analyze --provider openai --format json > analysis.json

# Export to Google Docs
mongo-migrate ai analyze --provider openai --google-docs

# Export to Google Docs with custom settings
mongo-migrate ai analyze --provider gemini --google-docs \
  --docs-title "Production DB Analysis $(date)" \
  --docs-share "team@company.com" \
  --docs-folder "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
```

### Collection-Specific Analysis
```bash
# Analyze specific collection schema
mongo-migrate ai schema users --provider gemini --detail

# Performance analysis
mongo-migrate ai performance --provider openai

# Oplog and replication analysis
mongo-migrate ai oplog --provider gemini --detail --google-docs

# Change stream analysis for specific collection
mongo-migrate ai changestream --collection orders --provider openai
```

### Different AI Providers
```bash
# Compare results from different providers
mongo-migrate ai analyze --provider openai > openai-analysis.txt
mongo-migrate ai analyze --provider gemini > gemini-analysis.txt
```

## Output Formats

### Text Format (Default)
Provides human-readable analysis with clear sections:
- Database Health Assessment
- Performance Recommendations  
- Index Optimization
- Schema Improvements
- Security Considerations
- Action Items

### JSON Format
Structured output perfect for automation:
```json
{
  "database_analysis": {
    "database_info": {...},
    "collections": [...],
    "indexes": [...],
    "performance": {...}
  },
  "ai_recommendations": "...",
  "provider": "OpenAI gpt-4o-mini",
  "timestamp": "2025-01-08T10:30:00Z"
}
```

## Command Reference

### Global Flags
- `--provider string`: AI provider (openai, gemini)
- `--detail`: Include detailed analysis
- `--config string`: Config file path

### Analyze Command
```bash
mongo-migrate ai analyze [flags]
```
- `--collection string`: Analyze specific collection
- `--format string`: Output format (text, json)

### Schema Command
```bash
mongo-migrate ai schema [collection] [flags]
```
- `--format string`: Output format (text, json)

### Performance Command
```bash
mongo-migrate ai performance [flags]
```
- `--format string`: Output format (text, json)

## Data Collected

The AI analysis collects the following MongoDB metrics:

### Database Information
- Database size and storage statistics
- Collection and index counts
- Overall database health metrics

### Collection Analysis
- Document counts and average sizes
- Sample document structures
- Index information per collection
- Storage utilization

### Performance Metrics
- Server status and operation counters
- Connection statistics
- Index usage patterns (when available)

### Security Note
**Important**: The tool only collects metadata and schema information. It does NOT send actual document data to AI providers, ensuring your sensitive data remains secure.

## Troubleshooting

### AI Analysis Disabled Error
```
Error: AI analysis is disabled. Set AI_ENABLED=true in your configuration
```
**Solution**: Add `AI_ENABLED=true` to your `.env` file

### API Key Not Configured
```
Error: OpenAI API key not configured. Set OPENAI_API_KEY environment variable
```
**Solution**: Set the appropriate API key for your chosen provider

### Provider Not Found
```
Error: unsupported AI provider: claude. Supported providers: openai, gemini (claude coming soon)
```
**Solution**: Use `openai` or `gemini` as providers. Claude support is coming soon.

## Best Practices

1. **Start with Basic Analysis**: Use the default `analyze` command first
2. **Compare Providers**: Different AI models may provide different insights
3. **Use JSON for Automation**: Export results in JSON format for further processing
4. **Focus on Specific Issues**: Use schema or performance commands for targeted analysis
5. **Regular Health Checks**: Run analysis periodically to catch issues early

## Coming Soon

- ✅ OpenAI GPT-4o support
- ✅ Google Gemini integration
- 🔄 Anthropic Claude support
- 📋 Custom analysis prompts
- 📊 Historical analysis comparison
- 🔄 Integration with monitoring tools
- 📈 Automated recommendations tracking
