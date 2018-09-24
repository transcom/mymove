import { getClient, checkResponse } from 'shared/Swagger/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateServiceMember(serviceMember) {
  // we create service members with no data associated with them.
  const client = await getClient();
  const response = await client.apis.service_members.createServiceMember({
    createServiceMemberPayload: serviceMember,
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

export async function CreateBackupContactAPI(serviceMemberId, backupContact) {
  // we create service members with no data associated with them.
  const client = await getClient();
  const response = await client.apis.backup_contacts.createServiceMemberBackupContact(
    {
      serviceMemberId: serviceMemberId,
      createBackupContactPayload: backupContact,
    },
  );
  checkResponse(
    response,
    'failed to create a backup contact due to server error',
  );
  return response.body;
}

export async function IndexBackupContactsAPI(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.backup_contacts.indexServiceMemberBackupContacts(
    {
      serviceMemberId,
    },
  );
  checkResponse(response, 'failed to get backup contacts due to server error');
  return response.body;
}

export async function UpdateBackupContactAPI(
  backupContactId,
  backupContactPayload,
) {
  const client = await getClient();
  const response = await client.apis.backup_contacts.updateServiceMemberBackupContact(
    {
      backupContactId,
      updateServiceMemberBackupContactPayload: backupContactPayload,
    },
  );
  checkResponse(
    response,
    'failed to update backup contact due to server error',
  );
  return response.body;
}

export async function SearchDutyStations(query) {
  const client = await getClient();
  const response = await client.apis.duty_stations.searchDutyStations({
    search: query,
  });
  checkResponse(response, 'failed to query duty stations due to server error');
  return response.body;
}
