# Changelog

## [1.0.1]
- Allow the five Korea-only response fields to be omitted outside Korea.
- Change the default Indonesia REST host to `https://api.upbit.co.id/`; retain explicit URL overrides.
- Add Indonesia API documentation links at `docs.upbit.co.id`; keep `global-docs.upbit.com` for Singapore and Thailand.
- Stop sending a request when body preparation or JWT signing fails.
- Preserve the position of `*bytes.Reader` and `*strings.Reader` request bodies and use independent copies for signing, sending, and retries; `*bytes.Buffer` bodies are also replayed. Other `io.Reader` bodies are sent once without replay.
- Avoid panics with custom HTTP transports.
- Redact sensitive headers regardless of case, duplicate keys, or multiple values.
- Refresh examples.

## [1.0.0]
- Add Pocket API (`pocket.go`): 7 new REST endpoints
- [BREAKING] Add missing `unit` field to minute candle response type (changes the min candle response type)

## [0.9.0]
- Initial release
