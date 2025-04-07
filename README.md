# Gator

Gator is a CLI RSS feed aggregator written in Go. Still actively working on this so expect changes, but it works.

You will need Postgres and Go installed to run Gator.

### Setup

Clone the repo:
```git clone https://github.com/drakedeloz/gator.git```

Install gator with:
```cd gator && go install .```

Create a config file in your home directory titled ```.gatorconfig.json```

Use psql or whatever postgres client you like to create a gator database and run the migrations using goose.