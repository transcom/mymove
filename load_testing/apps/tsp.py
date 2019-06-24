# -*- coding: utf-8 -*-
import random
from urllib.parse import urljoin

from locust import TaskSet
from locust import seq_task
from locust import task
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient

from .base import BaseTaskSequence
from .base import InternalAPIMixin
from .base import PublicAPIMixin
from .base import get_swagger_config
from .base import swagger_request


class TSPQueue(BaseTaskSequence, InternalAPIMixin, PublicAPIMixin):

    login_gov_user = None
    session_token = None
    user = {}

    def update_user(self):
        self.user = swagger_request(self.swagger_internal.users.showLoggedInUser)

    @seq_task(1)
    def login(self):
        resp = self.client.post("/devlocal-auth/create", data={"userType": "tsp"})
        try:
            self.login_gov_user = resp.json()
        except Exception as e:
            print(e)
            print("login could not be parsed", resp.content)
            self.interrupt()

        try:
            self.session_token = self.client.cookies.get("tsp_session_token")
        except Exception as e:
            print(e)
            print("missing session token")
            self.interrupt()

        self.requests_client = RequestsClient()
        # Set the session to be the same session as locust uses
        self.requests_client.session = self.client

        try:
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0
            # which doesn't respect nullable sub-definitions
            self.swagger_internal = SwaggerClient.from_url(
                urljoin(self.parent.parent.host, "internal/swagger.yaml"),
                request_headers={"x-csrf-token": self.csrf},
                http_client=self.requests_client,
                config=get_swagger_config(),
            )
            self.swagger_public = SwaggerClient.from_url(
                urljoin(self.parent.parent.host, "api/v1/swagger.yaml"),
                request_headers={"x-csrf-token": self.csrf},
                http_client=self.requests_client,
                config=get_swagger_config(),
            )
            # If either of these fails we can't continue
            if not (self.swagger_internal and self.swagger_public):
                self.interrupt()
        except Exception as e:
            print(e)
            self.interrupt()

    @seq_task(2)
    def retrieve_user(self):
        self.update_user()

    @seq_task(3)
    @task(10)
    def view_shipment_in_random_queue(self):
        """
        Choose a random queue to visit and pick a random move to view

        This task pretents to be a user who has work to do in a specific queue.
        """
        queue_types = ["AWARDED", "ACCEPTED", "APPROVED", "IN_TRANSIT", "DELIVERED"]
        q_type = random.choice(queue_types)

        queue = swagger_request(
            self.swagger_public.shipments.indexShipments, status=[q_type]
        )

        if not queue:
            return
        if len(queue) == 0:
            return

        # Pick a random shipment
        item = random.choice(queue)

        # These are all the requests loaded in a single move in rough order of execution

        shipment_id = item["id"]
        swagger_request(
            self.swagger_public.shipments.getShipment, shipmentId=shipment_id
        )

        swagger_request(
            self.swagger_public.service_agents.indexServiceAgents,
            shipmentId=shipment_id,
        )

        swagger_request(
            self.swagger_public.transportation_service_provider.getTransportationServiceProvider,
            shipmentId=shipment_id,
        )

        swagger_request(
            self.swagger_public.move_docs.indexMoveDocuments, shipmentId=shipment_id
        )

        swagger_request(
            self.swagger_public.accessorials.getTariff400ngItems,
            requires_pre_approval=True,
        )

        swagger_request(
            self.swagger_public.accessorials.getShipmentLineItems,
            shipmentId=shipment_id,
        )

        swagger_request(
            self.swagger_public.shipments.getShipmentInvoices, shipmentId=shipment_id
        )

        swagger_request(
            self.swagger_public.storage_in_transits.indexStorageInTransits,
            shipmentId=shipment_id,
        )

    @seq_task(4)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}


class TSPUserBehavior(TaskSet):
    tasks = {TSPQueue: 1}
