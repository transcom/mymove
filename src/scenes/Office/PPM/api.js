import { getClient, checkResponse } from 'shared/api';

export async function GetPpmIncentive(moveDate, originZip, destZip, weight) {
  const client = await getClient();
  const response = await client.apis.ppm.showPPMObligation({
    planned_move_date: moveDate,
    origin_zip: originZip,
    destination_zip: destZip,
    weight: weight,
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}
