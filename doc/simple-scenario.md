# simples scenario
spread of TE, no host defence

**simulation code**
```
/invade --seed 5 --genome MB:1,1 --rr 4,4 --u 0.1 --grid-x 50 --grid-y 20 --mate-radius 1 --epi-inherit none --basepop 8-12,1-3,10 --rep 1 --steps 5 --gen 100 --file-spatial simple.txt 
```



# visualization howto

requires `--file-spatial` output and R

```
library(tidyverse)
df <- read_tsv("/Users/rokofler/analysis/2024-spatialsimulations/2025-overview/2025-03-documentation/simple.txt") %>% select(-1)

# filter generations
df <-  df %>% filter(gen %in% c(0, 25, 50, 75, 100))

# bin
df$TEbin <- cut(df$TEs, c(seq(-1, 1500, 250),Inf))

# assign colors
bin_colors <- c("darkgrey", "#3DA53B") # "#db2c3fff",

# generate plot
pse <- ggplot(df, aes(x = X, y = Y)) +
  geom_point(aes(alpha = TEbin, color = silenced), size = 2) +
  scale_color_manual(values = bin_colors) +  # Use discrete bins
  facet_grid(rows = vars(gen), switch = "y") +
  theme(axis.text = element_blank(),
        axis.ticks = element_blank(),
       axis.title = element_blank(),
        strip.text.y = element_text(size = 14),
        strip.text.x = element_text(size = 14),
        legend.title = element_text(size = 14),  # Legend title size
       legend.text = element_text(size = 12)) +
  labs(color = "TE Status", alpha = "Copynumber")
# save to file
plot(pse)
ggsave("/Users/rokofler/analysis/2024-spatialsimulations/2025-overview/2025-03-documentation/simple.png",dpi=100)
```


