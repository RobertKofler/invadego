InvadeGo - Insertionbias branch
================


# parameters

## mandatory
* hermaphrodites with fitness depending mating probability; self-mating is not excluded;
* non-overlapping generations of a fixed population size (N)
* genome may be several chromosomes with some recombination rate


## optional
* --mu-bias: mutation rate of the insertion bias; probability that the insertion bias of a TE may change; either randomly increase or decrease by 1 (range -100 to +100); the insertion bias of any TE (piRNA cluster, genomic, novel, old, fixed) may change. this is justified by the fact that base substitutions could also change the sequence of any TE.
