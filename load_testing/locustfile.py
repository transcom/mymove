from locust import HttpLocust

from apps import AnonBehavior
from apps import MilMoveUserBehavior
from apps import OfficeUserBehavior
from apps import TSPUserBehavior


class AnonUser(HttpLocust):
    host = "http://milmovelocal:8080"
    # weight = 5  # 5x more likely than other users
    weight = 1
    task_set = AnonBehavior


class MilMoveUser(HttpLocust):
    host = "http://milmovelocal:8080"
    weight = 1
    task_set = MilMoveUserBehavior
    min_wait = 1000
    max_wait = 5000


class OfficeUser(HttpLocust):
    host = "http://officelocal:8080"
    weight = 1
    task_set = OfficeUserBehavior
    min_wait = 1000
    max_wait = 5000


class TSPUser(HttpLocust):
    host = "http://tsplocal:8080"
    weight = 1
    task_set = TSPUserBehavior
    min_wait = 1000
    max_wait = 5000
