# -*- coding: utf-8 -*-
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

    # User is the LoggedInUserPayload
    user = None
    duty_stations = []

    def update_user(self):
        self.user = swagger_request(self.swagger_internal.users.showLoggedInUser)

    def update_service_member(self, service_member):
        self.user.service_member = service_member

    def update_duty_stations(self, duty_stations):
        self.duty_stations = duty_stations

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
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0
            # which doesn't respect nullable sub-definitions
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
        self.update_user()

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
            social_security_number="333-33-3333",  # Random
            rank="E_5",  # Rotate
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
            middle_name="Carol",
            last_name="Bob",  # Random
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
            personal_email=self.user["email"],  # Email is derived from logging in
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
        payload = model(current_station_id=self.duty_stations[0].id)
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(9)
    def current_residence_address(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        address = self.swagger_internal.get_model("Address")
        payload = model(
            residential_address=address(
                street_address_1="12345 Fake St",
                city="Aurora",
                state="CO",
                postal_code="80013",
            )
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(10)
    def backup_mailing_address(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        address = self.swagger_internal.get_model("Address")
        payload = model(
            backup_mailing_address=address(
                street_address_1="12345 Fake St",
                city="Aurora",
                state="CO",
                postal_code="80013",
            )
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(11)
    def backup_contact(self):
        model = self.swagger_internal.get_model(
            "CreateServiceMemberBackupContactPayload"
        )
        payload = model(
            name="Alice",
            email="alice@example.com",
            permission="NONE",
            telephone="333-333-3333",
        )
        swagger_request(
            self.swagger_internal.backup_contacts.createServiceMemberBackupContact,
            serviceMemberId=self.user["service_member"].id,
            createBackupContactPayload=payload,
        )

    #
    # At this point the user profile is complete so let's refresh our knowledge
    # of the user's profile.
    #

    @seq_task(12)
    def refresh_user_profile(self):
        self.update_user()

    @seq_task(13)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}
