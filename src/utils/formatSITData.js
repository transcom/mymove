import { formatMoveHistoryFullAddress } from './formatters';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import { SIT_ADDRESS_UPDATE_STATUS } from 'constants/sitUpdates';

const formatStatus = (changedValues, eventName) => {
  if (
    (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.APPROVED && eventName === o.createSITAddressUpdate) ||
    (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.REQUESTED && eventName === o.createSITAddressUpdateRequest)
  ) {
    return 'Updated';
  }

  if (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.REJECTED && eventName === o.rejectSITAddressUpdate) {
    return 'Update request rejected';
  }

  if (changedValues.status === SIT_ADDRESS_UPDATE_STATUS.APPROVED && eventName === o.approveSITAddressUpdate) {
    return 'Update request approved';
  }

  return '';
};

export function formatSITData({ context, changedValues, eventName }) {
  if (!context) return {};

  const sitValues = {};
  if (context[0].sit_destination_address_final) {
    sitValues.sit_destination_address_final = formatMoveHistoryFullAddress(
      JSON.parse(context[0].sit_destination_address_final),
    );
  }

  if (context[0].sit_destination_address_initial) {
    sitValues.sit_destination_address_initial = formatMoveHistoryFullAddress(
      JSON.parse(context[0].sit_destination_address_initial),
    );
  }

  if (context[0].contractor_remarks) {
    sitValues.contractor_remarks = context[0].contractor_remarks;
  }

  if (changedValues?.status) {
    sitValues.status = formatStatus(changedValues, eventName);
  }

  return sitValues;
}

export default { formatSITData };
