# xisdb
A key-value store implemented in Go

### About
XisDB is an experimental low-level key-value database written in Go. Its purpose is to provide an API for a key-value store that is very fast, ACID compliant, thread-safe and easy to use. It is a work in progress, so don't use this in production or anything crazy. Basically I was inspired by LevelDB and BoltDB and was bored so I wrote this.

### Features
- In-memory
- Supports transactions and rollbacks
- Custom Indexes
- Query language
- Buckets of keys
- ACID compliant
- Disk Persistence

### Upcoming features
- PubSub on key changes
- Point-in-time restores
