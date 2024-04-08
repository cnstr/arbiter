# Canister Arbiter
> Dead simple orchestration for Canister workers.

Canister Arbiter is a simple tool for orchestrating Canister workers.
These workers (also known as indexers) are responsible for building data that is consumed by the [Canister API](https://github.com/cnstr/api).
Similar to how we maintain transparency by keeping the manifests list public, we also keep a public ledger of the available workers.
This is available at [`cnstr/manifests`](https://github.com/cnstr/manifests) and via an API endpoint which provides information on peering.

## Canister Peering
The current iteration of the Canister indexer runs as a single process that goes through a list of manifests and builds the data.
The issue is that this process can take a very long time to complete, making over 20,000 requests to various sources.
Measures have been implemented to avoid unnecessary work, but a full index can take almost an hour to complete.
On top of that, the indexer tends to trigger rate limits on some sources, leading to temporary IP bans.

To address these issues, we are splitting workloads across multiple volunteered peers.
Individuals have agreed to run a worker on their own infrastructure that gets directed by the orchestrator.
The plan is to eventually make the API, database, and orchestrator decentralized.

## Contributing
We'd appreciate anyone who is willing to host a peer for the Canister indexer.
This not only increases the reliability of the service, but also helps serve other regions faster.
In order to join the peering network, there are a few requirements that need to be met:

- Around `256MB` of RAM available at bursts every hour.
- Barely any CPU, it'll only occasionally burst up to half a thread.
- Relatively stable internet connection capable of opening PostgreSQL connections.
- You don't need to run the worker 24/7, but as much as possible is appreciated.
- Allow WebSocket connections to `peer-t1.canister.me`, `peer-t2.canister.me`, and `peer-t3.canister.me`.

Canister's Peering Network is secured by a custom certificate authority to ensure that only approved peers can participate.
In order to join, you'll need to generate a certificate signing request (CSR) and email it to [support@canister.me](mailto:support@canister.me).
Send that email with the subject `PEERING REQUEST: <Name>` where name is what you want to name the peer.
Make sure to provide adequate information as to who you are, where the peer is being hosted, etc.
If approved, we'll get back to you with a signed certificate which you can use to follow the next steps.

## Running a Peer
TODO
