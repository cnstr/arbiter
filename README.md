# Arbiter

A very simple peering control plane used to power Canister's peering capabilities.
Using websockets, the Arbiter facilitates communication between peers that do work for Canister's indexer.
Here's how it works:
- A peer attempts to upgrade to a Websocket to the arbiter (`workpool-connect.canister.me`)
- Peers message each other through the arbiter in order to participate in the quorum
  - Peers without certificates signed by the Canister Root CA cannot participate
  - Invalid peers get automatically disconnected by the arbiter
- The peers are able to split work evenly by themselves and keep Canister's indices up-to-date
  - This is a part of the closed-source indexer that Canister uses to build its data
