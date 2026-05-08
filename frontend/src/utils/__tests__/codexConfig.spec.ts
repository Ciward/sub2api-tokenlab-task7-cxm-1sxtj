import { describe, expect, it } from 'vitest'
import { buildCodexAuthJson, buildCodexConfigToml } from '@/utils/codexConfig'

describe('codexConfig utils', () => {
  it('builds a Codex TOML config that keeps the OpenAI provider id stable', () => {
    const config = buildCodexConfigToml('https://tokenlab.cc.cd')

    expect(config).toContain('model_provider = "OpenAI"')
    expect(config).toContain('[model_providers.OpenAI]')
    expect(config).toContain('base_url = "https://tokenlab.cc.cd"')
    expect(config).toContain('requires_openai_auth = true')
    expect(config).not.toContain('[model_providers.newapi]')
    expect(config).not.toContain('responses_websockets_v2 = true')
  })

  it('adds websocket flags only when requested', () => {
    const config = buildCodexConfigToml('https://tokenlab.cc.cd', {
      supportsWebsockets: true
    })

    expect(config).toContain('supports_websockets = true')
    expect(config).toContain('[features]')
    expect(config).toContain('responses_websockets_v2 = true')
  })

  it('renders auth.json with the expected OpenAI API key field', () => {
    expect(buildCodexAuthJson('sk-test')).toBe('{\n  "OPENAI_API_KEY": "sk-test"\n}')
  })
})
