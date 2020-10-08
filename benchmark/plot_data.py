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

bar_color = ("#09BAC8", "#FFF200", "#F72E5C", "#17B838")
ecolor = ""
grid_color = "#B0B0B0"
axis_heading = "#585858"
label_color = "#777777"

def plot_data(
    datafile,
    title,
    xlabel,
    ylabel,
    file,
    log=False,
):
    (avg, err) = get_plotable_tables(read_csv(datafile))
    plt.figure()

    df = pd.DataFrame(avg, columns=filters)
    p = df.plot.bar(yerr=err, width=0.8, color=bar_color)
    p.grid(axis="y", linestyle="dashed", color=grid_color)
    p.set_axisbelow(True)

    p.xaxis.label.set_color(axis_heading)
    p.yaxis.label.set_color(axis_heading)
    p.tick_params(axis='x', colors=label_color)
    p.tick_params(axis='y', colors=label_color)
    p.title.set_color(axis_heading)
    p.spines['bottom'].set_color(label_color)
    p.spines['left'].set_color(label_color)
    p.spines['top'].set_color((0,0,0,0))
    p.spines['right'].set_color((0,0,0,0))

    # Shrink current axis by 10%
    box = p.get_position()
    p.set_position([box.x0, box.y0, box.width * 0.9, box.height])

    leg = p.legend(loc="center left", bbox_to_anchor=(1, 0.5))
    for text in leg.get_texts():
        text.set_color(label_color)

    p.set_title(title)
    p.set_xticklabels(number_rules, rotation="horizontal")
    p.set_xlabel(xlabel)
    p.set_ylabel(ylabel)
    if log:
        p.set_yscale("log")

    fig = plt.gcf()
    fig.set_size_inches(12, 5)
    plt.savefig(file)

plot_data("throughput.csv", "Throughput", "Number of rules",
    "Throughput [Gbps]", "throughput.svg")

plot_data("cpu.csv", "CPU Usage", "Number of rules",
    "CPU Usage [Percent]", "cpu.svg")

plot_data("latency.csv", "Latency", "Number of rules",
    "Latency [us]", "latency.svg", log=True)

plot_data("setup.csv", "Rules Setup Time", "Number of rules",
    "Setup Time [us]", "setup.svg", log=True)
