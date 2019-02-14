# Load Testing

## locust.io

Getting started

```sh
cd load_testing/
brew install libev
virtualenv .venv -p python3
source .venv/bin/activate
pip install -r requirements.txt
```

In a separate window ensure that the app server is running with `make server_run`.

## Running tests with Web UI

```sh
locust --host=http://milmovelocal:8080 -f load_testing/locustfile.py
```

Then open [http://localhost:8089](http://localhost:8089/) and enter the number of users to simulate and the hatch rate.
Finally, hit the `Start swarming` button and wait for the tests to finish.

## Running tests from the CLI

You can run the test suite without the Web UI with a command similar to this:

```sh
locust --host=http://milmovelocal:8080 -f load_testing/locustfile.py --clients=50 --hatch-rate=5 --no-web --run-time=60s
```
