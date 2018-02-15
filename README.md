# gearbox-shim-calculator
Calculate what shims to use in you gearbox on assembly

expected input:
- list of shims available (thinkness in integers ex: in 1/100 mm)
- target thikness
- margin of error (tolerance)

## usage

```
Usage of ./bin/calculate:
  -M int
    	Maximums iteration depht (max shims in one set) (default 6)
  -m int
    	Margin on target (default 1)
  -shimlist string
    	Comma seperated list of shims (default "24,27,30,33,36,39,42,45,69,93,111,117,141")
  -t int
    	Target thickness (default 176)
  -threads int
    	Threads to use (default 2)
```

### Build

```
./make
```