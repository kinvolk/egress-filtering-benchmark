import csv

import numpy as np
import pandas as pd
import matplotlib
import matplotlib.pyplot as plt

# defined in parameters.py
from parameters import  (
    filters,
    number_rules,
    iterations,
)

def read_csv(filename):
    data = {}
    with open(filename, newline="") as csvfile:
        reader = csv.reader(csvfile, delimiter="\t")
        header_found = False
        for row in reader:
            # skip header
            if not header_found:
                header_found = True
                continue
            filter = row[0]
            rules = row[1]

            if not filter in data:
                data[filter] = {}

            if not rules in data[filter]:
                data[filter][rules] = {}

            data[filter][rules] = list(float(i) for i in row[2:])

    return data

def get_plotable_tables(data):
    avg = [[ 0 for x in filters ] for y in number_rules]
    err = [[ 0 for x in number_rules ] for y in filters]

    for (filter_index, filter) in enumerate(filters):
        for (rules_index, rule) in enumerate(number_rules):
            # TODO: add check if value if not found
            values = data[filter][rule]
            avg[rules_index][filter_index] = np.average(values)
            err[filter_index][rules_index] = t_err = np.std(values, ddof=1)/np.sqrt(len(values))
    return (avg, err)

########### plot throuhgput ####################
(avg, err) = get_plotable_tables(read_csv("throughput.csv"))
plt.figure()

df = pd.DataFrame(avg, columns=filters)
p = df.plot.bar(yerr=err, width=0.8)
p.grid(axis="y", linestyle="dashed")
p.set_axisbelow(True)

# Shrink current axis by 10%
box = p.get_position()
p.set_position([box.x0, box.y0, box.width * 0.9, box.height])

p.legend(loc="center left", bbox_to_anchor=(1, 0.5))
p.set_title("Throughput")
p.set_xticklabels(number_rules, rotation="horizontal")
p.set_xlabel("Number of rules")
p.set_ylabel("Throughput [Gbps]")
#p.set_yscale("log")

fig = plt.gcf()
fig.set_size_inches(12, 5)
plt.savefig("throughput.svg")
#plt.show()

########### plot cpu ####################
(avg, err) = get_plotable_tables(read_csv("cpu.csv"))
plt.figure()
df = pd.DataFrame(avg, columns=filters)
p = df.plot.bar(yerr=err, width=0.8)
p.grid(axis="y", linestyle="dashed")
p.set_axisbelow(True)

# Shrink current axis by 10%
box = p.get_position()
p.set_position([box.x0, box.y0, box.width * 0.9, box.height])

p.legend(loc="center left", bbox_to_anchor=(1, 0.5))
p.set_title("CPU Usage")
p.set_xticklabels(number_rules, rotation="horizontal")
p.set_xlabel("Number of rules")
p.set_ylabel("CPU Usage [Percent]")
#p.set_yscale("log")

fig = plt.gcf()
fig.set_size_inches(12, 5)
plt.savefig("cpu.svg")

########### plot latency #####################
(avg, err) = get_plotable_tables(read_csv("latency.csv"))
plt.figure()
df = pd.DataFrame(avg, columns=filters)
p = df.plot.bar(yerr=err, width=0.8)
p.grid(axis="y", linestyle="dashed")
p.set_axisbelow(True)

# Shrink current axis by 10%
box = p.get_position()
p.set_position([box.x0, box.y0, box.width * 0.9, box.height])

p.legend(loc="center left", bbox_to_anchor=(1, 0.5))
p.set_title("Latency")
p.set_xticklabels(number_rules, rotation="horizontal")
p.set_xlabel("Number of rules")
p.set_ylabel("Latency [ms]")
#p.set_yscale("log")

fig = plt.gcf()
fig.set_size_inches(12, 5)
plt.savefig("latency.svg")

########### plot setup time ####################
(avg, err) = get_plotable_tables(read_csv("setup.csv"))
plt.figure()

df = pd.DataFrame(avg, columns=filters)
p = df.plot.bar(yerr=err, width=0.8)
p.grid(axis="y", linestyle="dashed")
p.set_axisbelow(True)

# Shrink current axis by 10%
box = p.get_position()
p.set_position([box.x0, box.y0, box.width * 0.9, box.height])

p.legend(loc="center left", bbox_to_anchor=(1, 0.5))
p.set_title("Rules Setup Time")
p.set_xticklabels(number_rules, rotation="horizontal")
p.set_xlabel("Number of rules")
p.set_ylabel("Setup Time [us]")
p.set_yscale("log")

fig = plt.gcf()
fig.set_size_inches(12, 5)
plt.savefig("setup.svg")
#plt.show()
