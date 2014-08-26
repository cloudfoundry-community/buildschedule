Build schedule
==============

This CLI converts a YAML description of a schedule (such as for multi-day training) into a single HTML page.

An example 3-day training schedule might look like:

![3-day](http://cl.ly/image/1d2D3F3h1F1G/Example_3_day_training.png)

Usage
-----

The command is to be run from the root of your training project. An example tree view of your project might be:

```
.
├── public
│   ├── decks
│   │   ├── day-1-afternoon.md
│   │   ├── day-1-morning.md
│   │   ├── day-3-afternoon.md
│   │   └── day-3-morning.md
│   ├── index.html
│   └── labs
│       └── do-a-lab.md
└── schedules
    └── 3-day.yml
```

The schedule file `3-day.yml` would reference the slide decks & labs.

To generate the `public/index.html`:

```
buildschedule schedules/3-day.yml 2> public/index.html
```

See the [3-day.yml example](https://github.com/cloudfoundry-community/buildschedule/blob/master/examples/3-day/schedules/3-day.yml).

The generated `public/index.html` assumes that Bootstrap is included in the `public/bootstrap-3.2.0-dist` folder.

Installation
------------

```
go get -u github.com/cloudfoundry-community/buildschedule
```

Development
-----------

To locally build the tool:

```
go get ./...
```
