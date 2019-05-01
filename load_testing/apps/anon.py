from locust import TaskSet
from locust import task


class AnonBehavior(TaskSet):

    @task(1)
    def index(self):
        self.client.get("/")
