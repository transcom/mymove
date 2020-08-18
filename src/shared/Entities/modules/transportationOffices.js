import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { get } from 'lodash';

export const loadDutyStationTransporationOfficeLabel = 'TransportationOffice.loadDutyStationTransporationOffice';

export function loadDutyStationTransportationOffice(dutyStationId, label = loadDutyStationTransporationOfficeLabel) {
  const swaggerTag = 'transportation_offices.showDutyStationTransportationOffice';
  return swaggerRequest(getClient, swaggerTag, { dutyStationId }, { label });
}

export function selectDutyStationTransportationOffice(state, dutyStationId) {
  // TODO: Add the dutyStationId to the return payload
  // check for Transportation office that has the dutyStationId
  const offices = get(state, 'entities.TransportationOffices');
  const officesOfDutyStation = offices.filter((office) => office.duty_station_id === dutyStationId);

  return officesOfDutyStation[0];
}
