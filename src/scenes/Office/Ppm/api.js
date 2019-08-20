import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatDateString } from 'shared/utils';

export async function GetPpmIncentive(moveDate, originZip, originDutyStationZip, destZip, weight) {
  const client = await getClient();
  const response = await client.apis.ppm.showPPMIncentive({
    original_move_date: formatDateString(moveDate),
    origin_zip: originZip,
    origin_duty_station_zip: originDutyStationZip,
    destination_zip: destZip,
    weight: weight,
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}

export async function GetExpenseSummary(personallyProcuredMoveId) {
  const client = await getClient();
  const response = await client.apis.ppm.requestPPMExpenseSummary({
    personallyProcuredMoveId,
  });
  checkResponse(response, 'failed to retrieve summary due to server error');
  return response.body;
}
