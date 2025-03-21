# overview of parameters
## mandatory (for minimal simulation)

*  --basepop string; mandatory; coordinates of the individual(s) with the TE insertions in the first generation ; 1-based coordinates!; entries need to be in form Y,X,N:Y,X,N; where N' is feasible and means N silenced insertions; can also be in form  Ystart-Yend,Xstart-Xend,N
*   --epi-inherit string; mannatory; epigenetic inheritance of the silencing; either dros, ara, none
* --gen int; mandatory; run the simulations for '--gen' generations (default -1)
* --genome string; mandatory; the genomic landscape; e.g. 'MB:2,3,1,5' specifiies four chromosomes with sizes of 2,3,1,5 Mb
* --grid-x int; mandatory; the spatial grid size on X (default -1)
* --grid-y int; mandatory; the spatial grid size on Y (default -1)
* --mate-radius at which spatial distance in grid should mates be found; either integer or 'pan' for panmictici (default "2")
* --rr string; the recombination rate per chromosome in cm/Mb; e.g. '3,4,4,5' 
* --u float; the transposition rate of the transposon

## optional but important

* --seed int seed for the random number generator (default -1); default means current time in ticks
* -trigger-defense int at which threshold value should the host defense been triggerd (default -1)
* --max-insertions int the maximum number of insertions (default 10000)
* --selfing-rate float the selfing rate as a fraction; 0= no selfing 1= everything mating
* --rep int; the number of replicates to simulate (default 1)
* --uc float the transposition rate silenced TEs (assuming the host defence is not perfect)
* --steps int; report the output at each '--steps' generations (default 20)
* --p-sil-loss float; probabiltiy that epigenetic silencing can be lost
* --condinv flag; important convenience parameter; conditional on invasion, i.e. only show replicates where at least one TE copy was present until the final generation
* --console-format string; formating of console output, default summary: summary|gridoverview|gridcount|transectx (default "summary")
     
* --file-mhp string optional output file: position and population frequency of each TE insertion, resembles a manhatten plot
* --file-spatial string; **important optional output file**: detailed info about each individual specimens including coordinates

## optional for negative effect of TEs
*  --s float; the deleterious effect of a single homozygous TE insertion
*  --h float; the heterozygous effect of a TE insertion
*  --min-w float the minimum frequency of an average individual in the population (default 0.1)

## optional for population structure
* --mate-barrier string; a comma separated list about x-coordiantes where mating barriers will be introduced eg. '50,100' will introduce mate barriers at x-coordinates 50 and 100
* --barrier-strength float
    	the strength of the mate barriers; between 0.0 and 1.0; 1.0= total exclusion of any mating across boundary and 0.0= barrier has no effect 0.5= 50% reduced mating across boundary

## cosmetic parameters

* --sampleid string; the ID of the sample; will be a help in R to group samples like with facete_grid()
* --file-debug string optional output file for debugging various aspects
* --replicate-offset int starting index of the replicates; may be used for pseudo-parallelization) (default 1)
* -threads int  number of threads (default 1); currently multithreading is not supported (it probably never will be)
* --silent suppress output


 










