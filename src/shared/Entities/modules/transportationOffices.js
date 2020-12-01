import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { get, isNil } from 'lodash';

export const loadDutyStationTransportationOfficeOperation = 'TransportationOffice.loadDutyStationTransportationOffice';
const dutyStationTransportationOfficeSchemaKey = 'transportationOffice';

export function loadDutyStationTransportationOffice(
  dutyStationId,
  label = loadDutyStationTransportationOfficeOperation,
  schemaKey = dutyStationTransportationOfficeSchemaKey,
) {
  const swaggerTag = 'transportation_offices.showDutyStationTransportationOffice';
  return swaggerRequest(getClient, swaggerTag, { dutyStationId }, { label, schemaKey });
}

function selectCurrentDutyStation(state) {
  // TODO: change when service member is refactored
  return get(state, 'serviceMember.currentServiceMember.current_station');
}

export function selectDutyStationTransportationOffice(state) {
  // check for the service member's duty station outside of entities until refactored to be in entities
  const dutyStation = selectCurrentDutyStation(state);
  const transportationOffice = dutyStation.transportation_office;
  // check in entities for the loaded transportation office
  if (isNil(state.entities.transportationOffices)) {
    return {};
  }
  const offices = Object.values(state.entities.transportationOffices);
  const officesOfDutyStation = offices.find((office) => office.id === transportationOffice.id);

  return officesOfDutyStation;
}
