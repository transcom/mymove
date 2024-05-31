/* eslint-disable import/prefer-default-export */
import { getParamByName } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';

let standaloneCrateCap;
await getParamByName('standaloneCrateCap')
  .then((response) => {
    if (response.parameterValue != null) {
      standaloneCrateCap = parseFloat(response.parameterValue).toFixed(2);
    } else {
      standaloneCrateCap = 0;
    }
  })
  .catch((error) => {
    milmoveLogger.error(error);
    standaloneCrateCap = 0;
  });

export const TXO_PARAMS = {
  STANDALONE_CRATE_CAP: standaloneCrateCap,
};
