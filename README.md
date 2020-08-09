Horologium
===

Horologium is a simple Go library and command line tool that helps to keep track
of consumption and costs. Currently, it mainly addresses consumables billed by meters,
for example power, water, … .

App
---

The command line app is a simple tool to get an overview over costs and consumption in the
last months:

```shell script
$> horologium -lastMonths 6 powerConsumption.yml 
|   MONTH   | YEAR | CONSUMPTION |  COSTS  |
|-----------|------|-------------|---------|
| January   | 2020 |       22.11 | 1987.52 |
| February  |      |       33.59 | 2292.43 |
| March     |      |       48.41 | 2686.05 |
| April     |      |       44.59 | 2584.59 |
| May       |      |       48.61 | 2691.36 |
| June      |      |       66.58 | 3168.64 |
| July      |      |        0.00 | 1400.28 |

```
The above example was executed in July 2020, thus the last six months are evaluated.

The meter readings and pricing plans have to be given in a yaml file having
the following format:

```yaml
name: "A pseudo power consumption for testing"
consumptionFormat: "%.2f" # Optional, defaults to %.2f
currencyFormat: "%.2f"    # Optional, defaults to %.2f
plans:
  - {name: 2018, basePrice: 1241.34, unitPrice: 26.32, validFrom: "2018-01-01", validTo: "2018-01-01"}
  - {name: 2019, basePrice: 1341.12, unitPrice: 27.28, validFrom: "2019-01-01", validTo: "2019-01-01"}
  - {name: 2020, basePrice: 1400.28, unitPrice: 26.56, validFrom: "2020-01-01", validTo: "2020-01-01"}
readings:
  - {date: 2019-12-01, count: 1104.25}
  - {date: 2020-01-01, count: 1201.23}
  - {date: 2020-02-01, count: 1223.34}
  - {date: 2020-03-01, count: 1256.93}
  - {date: 2020-04-01, count: 1305.34}
  - {date: 2020-05-01, count: 1349.93}
  - {date: 2020-06-01, count: 1398.54}
  - {date: 2020-07-01, count: 1465.12}
```

Typically, every month there is a meter reading added to the file (with an external editor).

Date Interpretation
---
This app interpretes dates as being at the beginning of the day. Therefore, the range
2020-01-04 to 2020-01-06 spans 48 hours, from the January 4th 0:00 to January, 5th 24:00 or January, 6th 0:00.

Library
---
The code can also be used in another Go program as library. The exported functions are fully documented.

Future
---
This is a hobby project of mine. There are some things I like to improve and to add, I don't know
when or whether this will happen.

