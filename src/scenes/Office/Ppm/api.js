import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatDateForSwagger } from 'shared/dates';

export async function GetPpmIncentive(moveDate, originZip, originDutyLocationZip, ordersID, weight) {
  const client = await getClient();
  const response = await client.apis.ppm.showPPMIncentive({
    original_move_date: formatDateForSwagger(moveDate),
    origin_zip: originZip,
    origin_duty_location_zip: originDutyLocationZip,
    orders_id: ordersID,
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
