import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

export const getTspForShipmentLabel = 'Shipments.getTspForShipment';

export function getTspForShipment(label, shipmentId) {
  return swaggerRequest(
    getPublicClient,
    'transportation_service_provider.getTransportationServiceProvider',
    { shipmentId },
    { label },
  );
}
