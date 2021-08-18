# RemGo

Since I'm running a personal computer - somewhere in the early nineties of last century - I'm in search for the perfect calendar handling my appointments. I frequently changed my apps, but was never satisfied. Now a few weeks ago I switched on all my computers from various Debian and Ubuntu distributions to the rolling Arch Linux. That included a search for a perfect calendar again. On this way I came across a tool with a different approach. It's called _remind_ and is actively developed by the Canadian developer Dianne Skoll since 2005. Since I prefer running my calendars mostly in agenda mode I was very elated by the functionalities of this tool. Nevertheless there are some insufficiences regarding my personal workflows. In the simple terminal output times are shown only in am/pm format. One of the tools in the pipe chain of _remind_, _rem2ps_ and _ps2pdf_ can't handle unicode, so at least all German umlauts are missinterpreted. Finally there's not a possibility to display a sequence of days in the terminal agenda.- Because I really like the idea behind _remind_, these are enough reasons for me to build an own application based on the excellent reminder syntax.

## Features

These are the implemented features so far ...

* CLI
* separate configuration via _~/.config/remgo/config.yaml_
* useful default values for most parameters
* processing of basic REM constructs
* processing of the AT and DURATION attributes
* colorized output in terminals
* output for ...
  * today
  * a week
  * next week

There's just one ToDo for the moment. I'm still missing a simple output for postprocessing e.g. via mail. Should be done in a few weeks.

## Install

### Linux

* clone the repository
* move to the created directory
* run ```go install```
* copy the _sample.yaml_ to _~/.config/remgo/config.yaml_ and change it to your requirements

### Windows

* installation on Windows not tried so far
