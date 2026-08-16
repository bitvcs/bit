# bit

A next-generation centralized Version Control System engineered for binary-heavy projects (e.g., game development, 3D assets, rich media) and large team monorepos.

`bit` combines the strengths of Git (lightweight branching, Merge Requests, client 3-way text merges), Perforce (exclusive file locking, high-throughput binary streaming), and SVN (fine-grained path-based access control, centralized single source of truth) to provide scalable asset tracking with developer autonomy.

- **Language/Tech Stack:** Go (Golang), gRPC, Protobuf, SQLite (`modernc.org/sqlite`), FastCDC, BLAKE3
- **Architecture Model:** Client-Server Monorepo (`bit` / `bitd`)


## License

See [LICENSE](https://github.com/bitvcs/bit/blob/main/LICENSE).

