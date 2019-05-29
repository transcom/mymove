from urllib.parse import urljoin

from locust import seq_task
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient

from .base import BaseTaskSequence
from .base import InternalAPIMixin
from .base import get_swagger_config
from .base import swagger_request


class TSPUserBehavior(BaseTaskSequence, InternalAPIMixin):

    login_gov_user = None
    session_token = None
    user = {}

    @seq_task(1)
    def login(self):
        resp = self.client.post('/devlocal-auth/create', data={"userType": "tsp"})
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get('tsp_session_token')
            self.requests_client = RequestsClient()
            # Set the session to be the same session as locust uses
            self.requests_client.session = self.client
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0 which doesn't respect
            # nullable sub-definitions
            self.swagger_internal = SwaggerClient.from_url(
                urljoin(self.parent.host, "internal/swagger.yaml"),
                request_headers={'x-csrf-token': self.csrf},
                http_client=self.requests_client,
                config=get_swagger_config())
        except Exception:
            print(resp.content)

    @seq_task(2)
    def retrieve_user(self):
        resp = self.client.get("/internal/users/logged_in")
        self.user = resp.json()
        # check response for 200

    @seq_task(3)
    def view_new_moves_queue(self):
        swagger_request(
            self.swagger_internal.queues.showQueue,
            queueType="new")

    @seq_task(4)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}
