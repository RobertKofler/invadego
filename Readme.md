InvadeGo - Chassis branch
================

# Info about branch
The insertionbias branch performs simulations with TEs having an insertion bias towards or against piRNA clusters; It also allows to compete different families with different insertion biases


## Assumptions
* hermaphrodites with fitness depending mating probability; self-mating is not excluded;
* non-overlapping generations of a fixed population size (N)
* genome may be several chromosomes with some recombination rate
* simulated regions: piRNA clusters
* each TE insertion has a position and an insertion bias into piRNA clusters (int64 and byte)
* the insertion bias can evolve at some rate (--im)
* clonal propagation can be simulated
* two modes of TE silencing by an insertion into a piRNA cluster; i) a cluster insertion silences all insertions, irrespective of the insertion bias ii) a cluster insertion with a given insertion bias silences only those with the same insertion bias
* TEs may have negative effects; negative effect of insertions in clusters can be switched off

