import random
from urllib.parse import urljoin

from locust import seq_task
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient

from .base import BaseTaskSequence
from .base import InternalAPIMixin
from .base import PublicAPIMixin
from .base import get_swagger_config
from .base import swagger_request


class OfficeUserBehavior(BaseTaskSequence, InternalAPIMixin, PublicAPIMixin):

    login_gov_user = None
    session_token = None
    user = {}

    @seq_task(1)
    def login(self):
        resp = self.client.post('/devlocal-auth/create', data={"userType": "office"})
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get('office_session_token')
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
            self.swagger_public = SwaggerClient.from_url(
                urljoin(self.parent.host, "api/v1/swagger.yaml"),
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
    def view_moves_in_each_queue(self):
        queue_types = ["new", "ppm", "hhg_accepted", "hhg_delivered", "all"]
        for q_type in queue_types:
            queue = swagger_request(
                self.swagger_internal.queues.showQueue,
                queueType=q_type)
            # Only look at up to 5 randomly chosen items
            for item in random.sample(queue, min(len(queue), 5)):
                # These are all the requests you'd see loaded in a single move

                move_id = item["id"]
                # http://officelocal:3001/internal/moves/60098ff1-8dc9-4318-a2e8-47bc8aac11a4
                move = swagger_request(
                    self.swagger_internal.moves.showMove,
                    moveId=move_id)
                # http://officelocal:3001/internal/moves/60098ff1-8dc9-4318-a2e8-47bc8aac11a4/move_documents
                swagger_request(
                    self.swagger_internal.move_docs.indexMoveDocuments,
                    moveId=move_id)
                # http://officelocal:3001/api/v1/tariff_400ng_items?requires_pre_approval=true
                swagger_request(
                    self.swagger_public.accessorials.getTariff400ngItems,
                    requires_pre_approval=True)
                # http://officelocal:3001/internal/moves/60098ff1-8dc9-4318-a2e8-47bc8aac11a4/personally_procured_move
                swagger_request(
                    self.swagger_internal.ppm.indexPersonallyProcuredMoves,
                    moveId=move_id)

                orders_id = move["orders_id"]
                # http://officelocal:3001/internal/orders/a680accb-cc73-4af2-bd9c-283826cdd88f
                orders = swagger_request(
                    self.swagger_internal.orders.showOrders,
                    ordersId=orders_id)

                service_member_id = move["service_member_id"]
                # http://officelocal:3001/internal/service_members/61473913-36b8-425d-b46a-cee488a4ae71
                swagger_request(
                    self.swagger_internal.service_members.showServiceMember,
                    serviceMemberId=service_member_id)

                # http://officelocal:3001/internal/service_members/61473913-36b8-425d-b46a-cee488a4ae71/backup_contacts
                swagger_request(
                    self.swagger_internal.backup_contacts.indexServiceMemberBackupContacts,
                    serviceMemberId=service_member_id)

                shipment_id = orders["moves"][0]["shipments"][0]["id"]
                # http://officelocal:3001/api/v1/shipments/0ffa777e-3ea0-4f8f-b44b-ae0af309ebf4
                swagger_request(
                    self.swagger_public.shipments.getShipment,
                    shipmentId=shipment_id)

                # http://officelocal:3001/api/v1/shipments/0ffa777e-3ea0-4f8f-b44b-ae0af309ebf4/transportation_service_provider
                swagger_request(
                    self.swagger_public.transportation_service_provider.getTransportationServiceProvider,
                    shipmentId=shipment_id)

                # http://officelocal:3001/api/v1/shipments/0ffa777e-3ea0-4f8f-b44b-ae0af309ebf4/accessorials
                swagger_request(
                    self.swagger_public.accessorials.getShipmentLineItems,
                    shipmentId=shipment_id)

                # http://officelocal:3001/api/v1/shipments/0ffa777e-3ea0-4f8f-b44b-ae0af309ebf4/invoices
                swagger_request(
                    self.swagger_public.shipments.getShipmentInvoices,
                    shipmentId=shipment_id)

                # http://officelocal:3001/api/v1/shipments/0ffa777e-3ea0-4f8f-b44b-ae0af309ebf4/service_agents
                swagger_request(
                    self.swagger_public.service_agents.indexServiceAgents,
                    shipmentId=shipment_id)

                # http://officelocal:3001/api/v1/shipments/0ffa777e-3ea0-4f8f-b44b-ae0af309ebf4/storage_in_transits
                swagger_request(
                    self.swagger_public.storage_in_transits.indexStorageInTransits,
                    shipmentId=shipment_id)

    @seq_task(4)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}
