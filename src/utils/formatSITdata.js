import { formatMoveHistoryFullAddress } from './formatters';

export function formatSITData({ context }) {
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

  return sitValues;
}

export default { formatSITData };
