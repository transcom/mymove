import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { get } from 'lodash';

export const getTspForShipmentLabel = 'Shipments.getTspForShipment';

export function getTspForShipment(shipmentId, label = getTspForShipmentLabel) {
  return swaggerRequest(
    getPublicClient,
    'transportation_service_provider.getTransportationServiceProvider',
    { shipmentId },
    { label },
  );
}

// Selectors

export const selectTspById = (state, tspId) => state.entities.transportationServiceProviders[tspId] || {}; // eslint-disable-line security/detect-object-injection

export function selectTransportationServiceProviderForShipment(state, shipmentId) {
  const transportationServiceProviderId = get(
    state,
    `entities.shipments.${shipmentId}.transportation_service_provider_id`,
  );
  if (transportationServiceProviderId) {
    return selectTspById(state, transportationServiceProviderId);
  } else {
    return {};
  }
}
