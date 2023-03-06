InvadeGo - Insertionbias branch
================


# parameters

## mandatory parameters

* --N population size; diploid individuals; hermaphrodites
* --genome; define the genome; provide size and number of chromosomes; eg. 'MB:2,3,1,5' defines 4 chromosomes, the first has a size of 2Mbp, the second a size of 3Mbp than 1Mbp and the fourth 5Mbp
* --generations the number of generations to simulate (eg 5000)
* --basepop define the starting population the base population; two major options i) specify the exact position of each insertion in a file (eg: "file:/user/rober/sim/basepop.txt"); for the syntax of this file see below ii) provide the number and the bias for insertions with random positions; eg 100(50) will randomly distribute 100 insertions with a bias of 50 in the chromosomes; another example 100(50),100(-50) will randomly distribute 100 with a bias of 50 and another 100 with a bias of -50 in the chromosomes




## optional
* --mu-bias: mutation rate of the insertion bias; probability that the insertion bias of a TE may change; either randomly increase or decrease by 1 (range -100 to +100); the insertion bias of any TE (piRNA cluster, genomic, novel, old, fixed) may change. this is justified by the fact that base substitutions could also change the sequence of any TE.
* --u transposition rate; propability that a TE generates a novel copy
* --cluster define the size of piRNA clusters; one cluster for each chromosome needs to be specified, eg 'kb:1,1,1,1' specifies a cluster of 1kb at the beginning of each chromosome
* --sampleid the ID of the sample; convenience option for integration with R, especially tidyverse; will be a help in R to group samples like with facete_grid()
* --rr the recombination rate for each chromosome in cM/Mb  e.g. '3,4,4,5' the first chromosome has 3cM/Mb, the second 4cM/Mb and so on
* --x negative effect of a TE insertion: fitness w = 1-xn where n is the number of TE insertions
* --t synergistic effect of a TE insertion assuming: fitness w=1-xn^t 
* --no-x-cluins: switch; when activated insertions in piRNA cluster have no negative effect (ie. x=0 for cluster insertions)
* --selective-cluins: experimental feature; not yet sure it will be useful; insertions with a given bias (say 50) will solely be silenced by cluster insertions of the same bias (50); in other word, when activated a cluster insertion will not deactivate all TEs but only the TEs with the same bias; biological justification: different in the bias will be due to mutations in the TE, due to this mutations the TE may escape the piRNAs (pirnas solely tolerate a sequence divergence around 10-20%)
* --clonal: simulate clones; hence the simulated individuals will not mate and recombine anymore; they will make exact copies of themselves, but transpositions and mutations (--mu-bias) may still happen
* --uc residual transposition rate; per default it is assumed that a cluster insertion silences the TE completely (--uc 0), however to simulate some residual activity despite a cluster insertions use this feature
* --steps provide an output at each --steps generation
* --rep number of replicates to simulate
* --replicate-offset usually the very first replicate has number 1; sometimes this may not be wanted. eg when merging several files with simulations. with this option the first simulation may have a higher number; than files from different simulations may be merged and each simulation still has an unique replicate ID (useful for R tidyverse)
* --file-mhp output file given the position of the TE insertions (future options; show insertion bias?)
* --file-debug output solely relevant for debugging; currently used to compute LD decay
* --file-sfs output file with site-frequency-spectrum of TEs; not yet implemented; not sure it will be required
* --file-tally output file of insertion counts per each individual; not yet implemented; likely not needed
* --max-insertions in the absence of an effective host defence the TE copy numbers may increase exponentially, no computer can deal with super-large TE copy numbers; if the number of insertions exceeds --max-insertions the simulation will be terminated
* --min-w minimum fitness of a population; if the fitness gets lower than this value the simulations will be terminated; the population is extinct (default 0.1)
 simulations wi
* --seed for the simulations; if not provided current time in ticks is used
* --threads number of threads; not yet supported; likely never will
* --silent just show relevant output, not the messages

## Base pop file
With the flag `--basepop` a user may specify a base population file, eg. `file:/user/rober/sim/basepop.txt`

This file should look like in the following example:

```
250; 1(-2) 400(-12); 50000(-22)
250;;
500; 2(0) 10000(10) 50000(20);
```

The file may have between 1 and infinite many rows, where each row specifies the genotype for one or more individuals in the base population;
Each row has 3-columns that are delimited by `;`

* col-1: number of indivdiuals having the given genotype; **Note** the total number of individuals needs to sum to the population size `--N`
* col-2: genotype of the haploid genome 1
* col-3: genotype of the haploid genome 2

For each haploid genome the positions and the insertion bias (in bracket) of the TE insertions can be provided.
Some examples
* `500;;` simulate 500 individuals without any TE insertion
* `1;10(0);` simulate one indivdual with a heterozygous TE insertion at position 10, that has no insertion bias
* `5;10(0);10(0)` simulate 5 indivdiuals with a homozygous TE insertion at position 10, the insertion has no insertion bias
* `3;1(-50);5(50)
