
## Gymnastics of Golang

<div style="display: flex; align-items: center;">
  <img src="pic/gymnastics-1.jpg" alt="Gymnastics" width="500">
  <div style="margin-left: 10px;">
  maki is accumulating Golang patterns or just patterns in Golang I've encountered along the way here
  </div>
</div>

### Simple Patterns of Golang

1. Heartbeating Server
    * building primitives:
      - Ticker v.s Timer (stop, reset, defer stop to prevent leaking)
      - Select, Channel (close, ok, buffered, len(), etc)
      - Context (explicit cancellation)
    * where it is used:
      - Raft Leader
