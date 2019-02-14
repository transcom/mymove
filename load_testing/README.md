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

Then run tests:

```sh
locust --host=http://milmovelocal:8080 -f load_testing/locustfile.py
```

Then open [http://localhost:8089](http://localhost:8089/) and enter the number of users to simulate and the hatch rate.
Finally, hit the `Start swarming` button and wait for the tests to finish.
