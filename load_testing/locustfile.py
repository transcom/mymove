# -*- coding: utf-8 -*-
from locust import HttpLocust

from apps import AnonBehavior
from apps import ServiceMemberUserBehavior
from apps import OfficeUserBehavior
from apps import TSPUserBehavior


class AnonUser(HttpLocust):
    host = "http://milmovelocal:8080"
    # host = "https://my.experimental.move.mil"
    # weight = 5  # 5x more likely than other users
    weight = 1
    task_set = AnonBehavior


class ServiceMemberUser(HttpLocust):
    host = "http://milmovelocal:8080"
    # host = "https://my.experimental.move.mil"
    weight = 1
    task_set = ServiceMemberUserBehavior


class OfficeUser(HttpLocust):
    host = "http://officelocal:8080"
    # host = "https://office.experimental.move.mil"
    weight = 1
    task_set = OfficeUserBehavior


class TSPUser(HttpLocust):
    host = "http://tsplocal:8080"
    # host = "https://tsp.experimental.move.mil"
    weight = 1
    task_set = TSPUserBehavior
