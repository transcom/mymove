import { formatMoveHistoryFullAddress } from './formatters';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import { SIT_ADDRESS_UPDATE_STATUS, DESTINATION_SIT_ADDRESS_UPDATE_STATUS_FOR_UI } from 'constants/sitUpdates';

const formatStatus = (changedValues, eventName) => {
  if (
    (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.APPROVED && eventName === o.createSITAddressUpdate) ||
    (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.REQUESTED && eventName === o.createSITAddressUpdateRequest)
  ) {
    return DESTINATION_SIT_ADDRESS_UPDATE_STATUS_FOR_UI.UPDATED;
  }

  if (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.REJECTED && eventName === o.rejectSITAddressUpdate) {
    return DESTINATION_SIT_ADDRESS_UPDATE_STATUS_FOR_UI.REJECTED;
  }

  if (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.APPROVED && eventName === o.approveSITAddressUpdate) {
    return DESTINATION_SIT_ADDRESS_UPDATE_STATUS_FOR_UI.APPROVED;
  }

  return '';
};

// parse our address record and pass it off to existing helper function
// to format for Move history
const formatAddress = (address) => {
  return formatMoveHistoryFullAddress(JSON.parse(address));
};

// To get the information requested for destination SIT address
// updates, we need to utilize the context object in our
// move_history_fetcher query. This requires us to return an
// address record in JSON.
export function formatSITData({ context, changedValues, eventName }) {
  if (!context) return {};

  const sitValues = {};

  // use helpers to convert address from json to string for UI
  if (context[0]?.sit_destination_address_final) {
    sitValues.sit_destination_address_final = formatAddress(context[0]?.sit_destination_address_final);
  }

  if (context[0]?.sit_destination_address_initial) {
    sitValues.sit_destination_address_initial = formatAddress(context[0]?.sit_destination_address_initial);
  }

  // add contractor_remarks to return object, since it won't
  // exist in the changedValues in some cases
  if (context[0]?.contractor_remarks) {
    sitValues.contractor_remarks = context[0].contractor_remarks;
  }

  // format the status values to reflect client request
  if (changedValues?.status) {
    sitValues.status = formatStatus(changedValues, eventName);
  }

  return sitValues;
}

export default { formatSITData };
