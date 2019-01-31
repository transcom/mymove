import { get } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const approvePpmLabel = 'PPMs.approvePPM';

export function approvePPM(personallyProcuredMoveId) {
  const label = approvePpmLabel;
  const swaggerTag = 'office.approvePPM';
  return swaggerRequest(getClient, swaggerTag, { personallyProcuredMoveId }, { label });
}

export function selectPpmStatus(state, id) {
  const ppm = get(state, `entities.personallyProcuredMove.${id}`);
  if (ppm) {
    return ppm.status;
  } else {
    return get(state, 'office.officePPMs.0.status', '');
  }
}
