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

export function selectTspForShipment(state, shipmentId) {
  const tsp = Object.values(state.entities.transportationServiceProviders).find(tsp => tsp.shipment_id === shipmentId);
  return tsp || {};
}
