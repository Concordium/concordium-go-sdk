# Changelog

## 0.4.0

- Added support for protocol version 8.

## 0.3.0

- Added support for protocol version 7.

## 0.2.0

- Added support for GRPC V2 `GetWinningBakersEpoch` for getting a list of bakers that won the lottery in a particular historical epoch. Only available when querying a node with version at least 6.1.
- Added support for GRPC V2 `GetFirstBlockEpoch` for getting the block hash of the first finalized block in a specified epoch. Only available when querying a node with version at least 6.1.
- Added support for GRPC V2 `GetBakerEarliestWinTime` for getting the projected earliest time at which a particular baker will be required to bake a block. Only available when querying a node with version at least 6.1.
- Added support for GRPC V2 `GetBakerRewardPeriodInfo` for getting all the bakers in the reward period of a block. Only available when querying a node with version at least 6.1.
- Added support for GRPC V2 `GetBlockCertificates` for retrieving certificates for a block supporting ConcordiumBF, i.e. a node with at least version 6.1.
