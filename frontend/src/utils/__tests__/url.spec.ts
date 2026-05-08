import { describe, expect, it } from 'vitest'
import { resolveApiBaseUrl, sanitizeUrl } from '../url'

describe('url utils', () => {
  it('uses configured api_base_url when present', () => {
    expect(resolveApiBaseUrl('https://gateway.example.com', 'https://app.example.com')).toBe(
      'https://gateway.example.com'
    )
  })

  it('falls back to the current site origin when api_base_url is empty', () => {
    expect(resolveApiBaseUrl('', 'https://app.example.com')).toBe('https://app.example.com')
    expect(resolveApiBaseUrl('   ', 'https://app.example.com')).toBe('https://app.example.com')
    expect(resolveApiBaseUrl(undefined, 'https://app.example.com')).toBe(
      'https://app.example.com'
    )
  })

  it('sanitizes absolute http urls', () => {
    expect(sanitizeUrl('https://example.com/logo.png')).toBe('https://example.com/logo.png')
  })
})
