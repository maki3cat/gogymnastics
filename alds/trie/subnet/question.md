

# IP Geo-Location Service

## Problem Statement

From the network layer, it is a common behavior to find the geographical location of a request from its source ip. Assuming you’re designing a service to resolve the source location of a request, and you can build a map between ip subnets and geo locations like this:

| Subnet      | Location |
| ----------- | -------- |
| 15.0.0.0/24 | fr       |
| 10.0.0.0/16 | cn       |
| 10.0.1.0/25 | sh       |

Note that we should always take the one that has more specific matching (sh over cn if possible).

**API**
```kotlin
class Locator {
    fun setGeoLocation(subnet: String, geo: String)
    fun getGeoLocation(ip: String): String
}
```

## Possible Question Points

1) strings/strconv
2) bit calculation
3) trie tree
   - insert
   - search
   - last match priority
4) testing is hard/time-consuming, need coding accuracy