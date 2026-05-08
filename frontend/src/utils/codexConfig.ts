export interface BuildCodexConfigTomlOptions {
  supportsWebsockets?: boolean
}

export const buildCodexConfigToml = (
  baseUrl: string,
  options: BuildCodexConfigTomlOptions = {}
): string => {
  const lines = [
    'model_provider = "OpenAI"',
    'model = "gpt-5.4"',
    'review_model = "gpt-5.4"',
    'model_reasoning_effort = "xhigh"',
    'disable_response_storage = true',
    'network_access = "enabled"',
    'windows_wsl_setup_acknowledged = true',
    'model_context_window = 1000000',
    'model_auto_compact_token_limit = 900000',
    '',
    '[model_providers.OpenAI]',
    'name = "OpenAI"',
    `base_url = "${baseUrl}"`,
    'wire_api = "responses"'
  ]

  if (options.supportsWebsockets) {
    lines.push(
      'supports_websockets = true',
      'requires_openai_auth = true',
      '',
      '[features]',
      'responses_websockets_v2 = true'
    )
  } else {
    lines.push('requires_openai_auth = true')
  }

  return lines.join('\n')
}

export const buildCodexAuthJson = (apiKey: string): string => {
  return JSON.stringify(
    {
      OPENAI_API_KEY: apiKey
    },
    null,
    2
  )
}
