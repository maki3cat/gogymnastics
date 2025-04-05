

### Simple Patterns of Golang

1. Heartbeating Server
    * building primitives:
      - Ticker v.s Timer (stop, reset, defer stop to prevent leaking)
      - Select, Channel (close, ok, buffered, len(), etc)
      - Context (explicit cancellation)
    * where it is used:
      - Raft Leader
