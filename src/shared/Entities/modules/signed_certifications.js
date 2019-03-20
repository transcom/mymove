import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const Label = 'SignedCertifications.createSignedCertification';

export function createSignedCertification(
  moveId,
  payload /*shape: {personally_procured_move_id, shipment_id, certification_text, signature, date, certification_type}*/,
) {
  const { certification_text, signature, date, personally_procured_move_id, certification_type, shipment_id } = payload;
  return swaggerRequest(
    getClient,
    'certification.createSignedCertification',
    {
      moveId: moveId,
      createSignedCertificationPayload: {
        certification_text: certification_text,
        signature: signature,
        date: date,
        shipment_id: shipment_id,
        personally_procured_move_id: personally_procured_move_id,
        certification_type: certification_type,
      },
    },
    { Label },
  );
}
