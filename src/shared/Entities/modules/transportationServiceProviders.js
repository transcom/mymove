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

// Selectors

export const selectTspById = (state, tspId) => state.entities.transportationServiceProviders[tspId] || {}; // eslint-disable-line security/detect-object-injection
