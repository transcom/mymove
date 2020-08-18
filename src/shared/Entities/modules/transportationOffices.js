import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { get } from 'lodash';

export const loadDutyStationTransporationOfficeLabel = 'TransportationOffice.loadDutyStationTransporationOffice';

export function loadDutyStationTransportationOffice(dutyStationId, label = loadDutyStationTransporationOfficeLabel) {
  const swaggerTag = 'transportation_offices.showDutyStationTransportationOffice';
  return swaggerRequest(getClient, swaggerTag, { dutyStationId }, { label });
}

function selectCurrentDutyStation(state) {
  // TODO: change when service member is refactored
  return get(state, 'serviceMember.currentServiceMember.current_station');
}

export function selectDutyStationTransportationOffice(state) {
  // check for the service member's duty station outside of entities until refactored to be in entities
  const dutyStation = selectCurrentDutyStation(state);
  const transportationOffice = dutyStation.transportation_office;
  // check in entities for the loaded transporation office
  const offices = get(state, 'entities.TransportationOffices');
  const officesOfDutyStation = offices.filter((office) => office.id === transportationOffice.id);

  return officesOfDutyStation[0];
}
