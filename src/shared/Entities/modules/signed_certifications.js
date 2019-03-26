import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const Label = 'SignedCertifications.createSignedCertification';

export function createSignedCertification(
  moveId,
  payload /*shape: {personally_procured_move_id, shipment_id, certification_text, signature, date, certification_type}*/,
) {
  return swaggerRequest(
    getClient,
    'certification.createSignedCertification',
    {
      moveId: moveId,
      createSignedCertificationPayload: {
        ...payload,
      },
    },
    { Label },
  );
}
