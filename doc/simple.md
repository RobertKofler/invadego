simple simulation scenarios
================

# Simplest simulation - TE invasion without host defence

## simple - mate radius 1

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

![](simple_files/figure-gfm/unnamed-chunk-2-1.png)<!-- -->

**Note** the TE starts spreading at the small rectangular patch at
generation 1; it is invading the population uninhibited by any host
defence. furthermore very high TE copy numbers per individual are
reached by generation 100

## simple - mate radius 2

``` bash
# simulation code
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 2 --epi-inherit none --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --file-spatial simple-mr2.txt 
```

**parameters explained**

- all parameters like before except
- –mate-radius 2: hence potential mates can be found more distantly of
  the focal species (as compared to mate-radius 1)

![](simple_files/figure-gfm/unnamed-chunk-4-1.png)<!-- -->

**Note** with a larger mate radius the TE will spread faster

# Host defence

## host defence - no epigenetic inherited silencing

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit none --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --file-spatial simple-defence.txt
```

**parameters explained**

- parameters are like before except mate-radius set to 1 and
- –trigger-defense 40: host defence will be triggered at 40 copies;
  i.e. individuals with 40 TE insertions will have a transposition rate
  u=0.0

![](simple_files/figure-gfm/unnamed-chunk-6-1.png)<!-- -->

**Note** the TE starts spreading at the small patch by generation 1. By
generation 50 already several indivdiuals have triggered the host
defence (no epigenetic inheritance in this scenario); By generation 100,
bot the TE and the host defence is spreading; **important** only
individuals with 40 or more TE copies can have the silencing state

## host defence - uniparental epigenetic silencing inheritance

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

![](simple_files/figure-gfm/unnamed-chunk-8-1.png)<!-- -->

**Note** in this scenario even individuals with \<40 TEs may have
silenced TEs (compare to previous scenario with –epi-inheritance none )

## host defence - biparental epigenetic silencing inheritance

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit ara --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --file-spatial simple-defence-ara.txt 
```

**parameters explained**

- –epi-inherit ara: biparental epigenetic inherited silencing; if a TE
  is silenced in any parent, the TE will be silenced in all offspring
  (even if the TE is still active in the other parent)

![](simple_files/figure-gfm/unnamed-chunk-10-1.png)<!-- -->

**Note** how the epigenetic silencing is spreading faster than the TE
such that by generation 100 all TE insertions are silenced; furthermore

## host defence - biparental epigenetic silencing and selfing

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit ara --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --selfing-rate 0.95 --condinv --file-spatial simple-defence-ara-self.txt 
```

**parameters explained**

- –selfing-rate 0.95: there is a 95% chance that the mother and the
  father will be identical
- –condinv: conditional on invasion; this is just a convenience method
  (rescuing me from finding seeds where it does work)

![](simple_files/figure-gfm/unnamed-chunk-12-1.png)<!-- -->

# Negative selection against inactive TEs

## negative selection - biparental inheritance

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit ara --basepop "1-20,1-50,40'" --rep 1 --steps 5 --gen 400 --trigger-defense 40 --s 0.1 --h 0.5 --file-spatial simple-negsel-ara.txt
```

**parameters explained**

- –s 0.1: selection against homozygous insertions, fitness w=1-s
- –h 0.5: selection against heterozygous insertions, fitness is w=1-hs
- –basepop “1-20,1-50,40’” all individuals in the entire population will
  have 40 **silenced** TE insertions (the dash indicates the silencing)
  ![](simple_files/figure-gfm/unnamed-chunk-14-1.png)<!-- -->

**Note** selection against TEs may lead to a patchy distribution **when
TEs are silenced biparental**

## negative selection - uniparental inheritance

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit dros --basepop "1-20,1-50,40'" --rep 1 --steps 5 --gen 400 --trigger-defense 40 --s 0.1 --h 0.5 --file-spatial simple-negsel-dros.txt 
```

![](simple_files/figure-gfm/unnamed-chunk-16-1.png)<!-- -->

**Now this is a surprise to me! negative selection against TEs with
uniparental silencing does not lead to patchyness** actually it
reactivates the TE on a broad scale, even when the TE was silenced in
the entire base population! thats truly fascinating! but makes sense, as
soon as some individuals have zero TEs (thanks to neg.sel.) the TE can
be reactivated by paternal crosses

# Population structure - uniparental inheritance

``` bash
./invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit dros --basepop 9-11,1-2,10 --rep 1 --steps 5 --gen 100 --trigger-defense 40 --mate-barrier 25 --barrier-strength 0.99 --file-spatial simple-matebar.txt 
```

**parameters explained**

- –mate-barrier 25: position of mating barrier in the x-grid; it will be
  more difficult for individuals at opposite ends of the barrier to mate
  eg with x-coordinates 25 and 26
- –barrier-strength 0.99: a 99% reduction in mating-probabilty for
  individuals on opposite ends of the barrier

![](simple_files/figure-gfm/unnamed-chunk-18-1.png)<!-- -->

**Note** the coordinates of the mate barrier can be clearly discerned at
generation 100; also not how all individuals in the left subpopulation
have about 40 TE copies (ie the silencing threshold), whereas most
individuals in the right subpopulation have very few insertions
