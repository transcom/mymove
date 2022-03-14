import { isNil } from 'lodash';

import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { selectCurrentDutyLocation } from 'store/entities/selectors';

export const loadDutyLocationTransportationOfficeOperation =
  'TransportationOffice.loadDutyLocationTransportationOffice';
const dutyLocationTransportationOfficeSchemaKey = 'transportationOffice';

export function loadDutyLocationTransportationOffice(
  dutyLocationId,
  label = loadDutyLocationTransportationOfficeOperation,
  schemaKey = dutyLocationTransportationOfficeSchemaKey,
) {
  const swaggerTag = 'transportation_offices.showDutyLocationTransportationOffice';
  return swaggerRequest(getClient, swaggerTag, { dutyLocationId }, { label, schemaKey });
}

export function selectDutyLocationTransportationOffice(state) {
  // check for the service member's duty location outside of entities until refactored to be in entities
  const dutyLocation = selectCurrentDutyLocation(state);
  const transportationOffice = dutyLocation.transportation_office;
  // check in entities for the loaded transportation office
  if (isNil(state.entities.transportationOffices)) {
    return {};
  }
  const offices = Object.values(state.entities.transportationOffices);
  const officesOfDutyLocation = offices.find((office) => office.id === transportationOffice.id);

  return officesOfDutyLocation;
}
