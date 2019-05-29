from urllib.parse import urljoin

from locust import seq_task
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient

from .base import BaseTaskSequence
from .base import InternalAPIMixin
from .base import get_swagger_config
from .base import swagger_request


class MilMoveUserBehavior(BaseTaskSequence, InternalAPIMixin):

    login_gov_user = None
    session_token = None

    # User is where service member data is stored
    user = {}

    def update_service_member(self, service_member):
        self.user["service_member"] = service_member

    def update_duty_stations(self, duty_stations):
        self.user["duty_stations"] = duty_stations

    @seq_task(1)
    def login(self):
        resp = self.client.post("/devlocal-auth/create", data={"userType": "milmove"})
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get("mil_session_token")
            self.requests_client = RequestsClient()
            # Set the session to be the same session as locust uses
            self.requests_client.session = self.client
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0 which doesn't respect
            # nullable sub-definitions
            self.swagger_internal = SwaggerClient.from_url(
                urljoin(self.parent.host, "internal/swagger.yaml"),
                request_headers={"x-csrf-token": self.csrf},
                http_client=self.requests_client,
                config=get_swagger_config(),
            )
        except Exception:
            print(resp.content)

    @seq_task(2)
    def retrieve_user(self):
        resp = self.client.get("/internal/users/logged_in")
        self.user = resp.json()
        # check response for 200

    @seq_task(3)
    def create_service_member(self):
        model = self.swagger_internal.get_model("CreateServiceMemberPayload")
        payload = model(user_id=self.user["id"])
        service_member = swagger_request(
            self.swagger_internal.service_members.createServiceMember,
            createServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(4)
    def create_your_profile(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            affiliation="NAVY",  # Rotate
            edipi="3333333333",  # Random
            rank="E_5",  # Rotate
            social_security_number="333-33-3333",  # Random
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(5)
    def create_your_name(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            first_name="Alice",  # Random
            last_name="Bob",  # Random
            middle_name="Carol",
            suffix="",
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(6)
    def create_your_contact_info(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            email_is_preferred=True,
            personal_email="20190321164732@example.com",
            phone_is_preferred=True,
            secondary_telephone="333-333-3333",
            telephone="333-333-3333",
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(7)
    def search_for_duty_station(self):
        station_list = ["b", "buck", "buckley"]
        for station in station_list:
            duty_stations = swagger_request(
                self.swagger_internal.duty_stations.searchDutyStations, search=station
            )
            self.update_duty_stations(duty_stations)

    @seq_task(8)
    def current_duty_station(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(current_station_id=self.user["duty_stations"][0].id)
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(9)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}
