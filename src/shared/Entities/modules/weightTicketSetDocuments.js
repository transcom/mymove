import { getClient } from 'shared/Swagger/api';
import { swaggerRequest } from 'shared/Swagger/request';

export const createWeightTicketSetDocumentLabel = 'weightTicketDocumentSet.createWeightTicketDocument';
// payload shape
// {
//     personally_procured_move_id,
//     upload_ids,
//     weight_ticket_set_type,
//     vehicle_nickname,
//     vehicle_make,
//     vehicle_model,
//     empty_weight_ticket_missing,
//     empty_weight,
//     full_weight_ticket_missing,
//     full_weight,
//     weight_ticket_date,
//     trailer_ownership_missing
// },
export function createWeightTicketSetDocument(moveId, payload, label = createWeightTicketSetDocumentLabel) {
  const swaggerTag = 'move_docs.createWeightTicketDocument';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      moveId,
      createWeightTicketDocument: { ...payload },
    },
    { label },
  );
}
