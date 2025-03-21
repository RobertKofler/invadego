simple simulation scenarios
================

## Simplest simulation - TE invasion without host defence

``` bash
# simulation code
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit none --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --file-spatial simple.txt 
```

**parameters explained**

- seed the random seed
- a genome of 2 chromsomes each with a length of 1MB
- a recombination rate of 4cM/Mb per chromosome
- the transposition rate is u=0.1

![](simple_files/figure-gfm/unnamed-chunk-2-1.png)<!-- -->
