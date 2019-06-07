import { getClient } from 'shared/Swagger/api';
import { swaggerRequest } from 'shared/Swagger/request';
export const createWeightTicketSetDocumentLabel = 'weightTicketDocumentSet.createWeightTicketDocument';
// payload shape
// {
//     personally_procured_move_id,
//     upload_ids,
//     vehicle_options,
//     vehicle_nickname,
//     empty_weight_ticket_missing,
//     empty_weight,
//     full_weight_ticket_missing,
//     full_weight,
//     weight_ticket_date,
//     trailer_ownership_missing
// },
// "operation MoveDocsCreateWeightTicketDocument has not yet been implemented"
// operation move_docs.createWeightTicketDocument failed: Error: Not Implemented (501)
export function createWeightTicketSetDocument(moveId, payload, label = createWeightTicketSetDocumentLabel) {
  return swaggerRequest(
    getClient,
    'move_docs.createWeightTicketDocument',
    {
      moveId,
      createWeightTicketDocument: { ...payload },
    },
    { label },
  );
}
