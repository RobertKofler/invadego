# overview of parameters
## mandatory (for minimal simulation)

*  --basepop string; mandatory; coordinates of the individual(s) that trigger the TE invasion; 1-based coordinates; entries need to be in form y1-y2,x1-x2,count where y1-y2 the a range of the y-coordinates, x1-x2 the range of the x-coordinates; all individuals in this range will have exactly 'count' transposons insertions; the genomic coordinates of the insertions are at random positions in the user specified genome and will thus have an individual population frequency of 1/2N where N is the populatin size
* --epi-inherit string; mannatory; epigenetic inheritance of the silencing of transposons; either dros, ara, none; none= no epigenetic inheritance; dros= uniparental inheritance of the silencing, i.e. mothers transmit the silencing; ara= biparental inheritance, both parents transmit the silencing
* --gen int; mandatory; run the simulations for '--gen' generations (default -1)
* --genome string; mandatory; the genomic landscape; e.g. 'MB:2,3,1,5' specifiies four chromosomes with sizes of 2,3,1,5 Mb
* --grid-x int; mandatory; the spatial grid size on X (default -1)
* --grid-y int; mandatory; the spatial grid size on Y (default -1)
* --mate-radius string; radius in the grid at which potential mates can be found; e.g. given the y-x coordinates 3-3 and a mate radius 1, potential mates may have the coordinates from 2-2 to 4-4 (including 2-3, 3-3, 3-4, 4-3 etc); it is also feasible to provide 'pan' for panmictic which ignores the spatial information for finding
* --rr string; the recombination rate per chromosome in cm/Mb; e.g. '3,4,4,5' 
* --u float; the transposition rate of the transposon

## optional but important

* --seed int seed for the random number generator (default -1); the default -1 translates to the current time in ticks
* --trigger-defense int; the copy number of TEs per individual at which the host defence will be triggerd (default -1); per default the host defence will not be at any TE copy number
* --max-insertions int; the maximum number of insertions per individual (default 10000); simulations will be aborted if the average TE copy number per individual exceeds this threshold; this option can be used to prevent the exponential growth of TE copy numbers from crashing the computitional resources
* --selfing-rate float; the selfing rate as fraction between 0 and 1; for example 0.32 specifies 32% selfing
* --rep int; the number of replicates that should be simulated (default 1)
* --uc float; the transposition rate for TEs silenced by the host defence (the basic idea is that the host defence may be slightly inefficient permitting some residual activity of the TE)
* --steps int; report the output at each '--steps' generations (default 20)
* --p-sil-loss float; probabiltiy that epigenetic silencing can be lost between 0 - 1, for example 0.01 specifies a 1% chance to loose silencing of the TE
* --condinv flag; important convenience parameter; conditional on invasion, i.e. only count replicates where at least one TE copy was present at the very last generation (--gen)
* --console-format string; formating of console output, default summary: summary|gridoverview|gridcount|transectx (default "summary")  
* --file-mhp string; optional output file: position and population frequency of each TE insertion; this output could be used to generate plots resembling Manhattenplots, where the genomic coordinates are on the x-axis and the population frequency of the insertions on the y-axis 
* --file-spatial string; **important optional output file**: detailed info about each individual specimens including coordinates; can be used to visualize invasions; 

## optional for negative effect of TEs
*  --s float; the deleterious effect of a single homozygous TE insertion (w=1-s for homozygot insertions)
*  --h float; the heterozygous effect of a TE insertion (w=1-hs for heterozygous insertions)
*  --min-w float the minimum average fitness (default 0.1); will assume the population died out when the averge fitness is lower than this

## optional for population structure
* --mate-barrier string; a comma separated list of X-coordiantes where mate-barriers will be introduced; for example '50,100' will introduce mate barriers at X-coordinates 50 and 100; individuals on opposite ends of a barrier will have a reduced propability to mate with each other.
* --barrier-strength float; the strength of the mate barriers between 0 - 1; for example 0.8 specifies a 80% reduction in mating proability for individuals on opposite ends of the barrier

## cosmetic parameters

* --sampleid string; specify an ID of the sample; could be convenient for loading the output into R e.g. to group samples or to use with facete_grid()
* --replicate-offset int; numbering of the replicates will start with '--replicate-offset'; could be convenient to merge simulations performed on different computers; eg when using pseudo-parallelization) (default 1)
* --threads int; not supported (number of threads; currently multithreading is not supported)
* --file-debug string; optional output file that may be used for debugging
* --silent suppress unnecesarry output


 










