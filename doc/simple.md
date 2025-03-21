simple simulation scenarios
================

## Simplest simulation - TE invasion without host defence

``` bash
# simulation code
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit none --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --file-spatial simple.txt 
```

**parameters explained**

- –seed 5: seed the random seed
- –genome MB:1,1: a genome of 2 chromsomes each with a length of 1MB
- –rr 4,4: a recombination rate of 4cM/Mb per chromosome
- –u 0.1: the transposition rate is u=0.1
- –grid-x 50: the x-size of the 2D grid is 50
- –grid-y 20: the y-size of the 2D grid is 20
- –mate-radius 1: the mate radius is 1; ie when generating the next
  generation, potential mates will be found in a radius of 1 of the
  focal individual. eg focal individual has y-x coordinates 3-3, then
  the potential mates have coordinates 2-2, 2-3, 2-4, 3-2, 3-3, 3-4,
  4-2, 4-3, 4-4
- –epi-inherit none: so far we assume no epigenetic inheritance of the
  TE silencing
- –basepop 9-11,1-2,10: all individuals in the y-range 9-11 and the
  x-range 1-2 will have exactly 10
- –rep 1: we simulate a single replicate
- –steps 5: output will be generated all 5 generations; note that we
  perform additional filtering in R, hence the final figure does not
  show the results for all 20 time points with an output –gen 100:
  simulations will be performed for 100 generations –file-spatial
  simple.txt: the output file for generating the figure

![](simple_files/figure-gfm/unnamed-chunk-2-1.png)<!-- --> **Note** the
TE starts spreading at the small rectangular patch at generation 1; it
is invading the population uninhibited by any host defence. furthermore
very high TE copy numbers per individual are reached by generation 100

## Host defence - no epigenetic inherited silencing

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit none --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --file-spatial simple-defence.txt
```

**parameters explained**

- all parameters are like before except
- –trigger-defense 40: host defence will be triggered at 40 copies;
  i.e. individuals with 40 TE insertions will have a transposition rate
  u=0.0

![](simple_files/figure-gfm/unnamed-chunk-4-1.png)<!-- --> **Note** the
TE starts spreading at the small patch by generation 1. By generation 50
already several indivdiuals have triggered the host defence (no
epigenetic inheritance in this scenario); By generation 100, bot the TE
and the host defence is spreading; **important** only individuals with
40 or more TE copies can have the silencing state

## Host defence - uniparental epigenetic silencing inheritance

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit dros --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --file-spatial simple-defence-dros.txt 
```

**parameters explained**

- –epi-inherit dros: uniparental epigenetic inherited silencing;
  silencing is transmitted maternally as in Drosophila. since
  hermaphrodites are simulated one mate is randomly chosen as the
  mother; if a TE is silenced in the mother any offspring with at least
  one TE insertion will inherit the silencing-status. offspring that per
  chance do not have a single insertion will not inherit the silencing
  status; paternal offspring will have active TEs (unless the mother has
  silenced TEs)

![](simple_files/figure-gfm/unnamed-chunk-6-1.png)<!-- --> **Note** in
this scenario even individuals with \<40 TEs may have silenced TEs
(compare to previous scenario with –epi-inheritance none )

## Host defence - biparental epigenetic silencing inheritance

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit ara --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --file-spatial simple-defence-ara.txt 
```

**parameters explained**

- –epi-inherit ara: biparental epigenetic inherited silencing; if a TE
  is silenced in any parent, the TE will be silenced in all offspring
  (even if the TE is still active in the other parent)

![](simple_files/figure-gfm/unnamed-chunk-8-1.png)<!-- --> **Note** how
the epigenetic silencing is spreading faster than the TE such that by
generation 100 all TE insertions are silenced; furthermore
