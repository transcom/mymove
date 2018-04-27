import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateServiceMember(service_member) {
  // we create service members with no data associated with them.
  const client = await getClient();
  const response = await client.apis.service_members.createServiceMember({
    createServiceMemberPayload: service_member,
  });
  checkResponse(
    response,
    'failed to create a blank service member profile due to server error',
  );
  return response.body;
}

export async function GetServiceMember(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.service_members.showServiceMember({
    serviceMemberId,
  });
  checkResponse(response, 'failed to get service member due to server error');
  return response.body;
}

export async function UpdateServiceMember(
  serviceMemberId,
  serviceMemberPayload,
) {
  const client = await getClient();
  const response = await client.apis.service_members.patchServiceMember({
    serviceMemberId,
    patchServiceMemberPayload: serviceMemberPayload,
  });
  checkResponse(
    response,
    'failed to update service member due to server error',
  );
  return response.body;
}

export async function SearchDutyStations(affiliation, query) {
  const client = await getClient();
  const response = await client.apis.duty_stations.searchDutyStations({
    affiliation: affiliation,
    search: query,
  });
  checkResponse(response, 'failed to query duty stations due to server error');
  return response.body;
}
