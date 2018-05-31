import { getClient, checkResponse } from 'shared/api';

export async function showDutyStationTransportationOffice(dutyStationId) {
  const client = await getClient();
  const response = await client.apis.transportation_offices.showDutyStationTransportationOffice(
    {
      dutyStationId,
    },
  );
  checkResponse(response, 'failed to get orders due to server error');
  return response.body;
}
