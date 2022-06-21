import { formatMoveHistoryFullAddress } from 'utils/formatters';
import a from 'constants/MoveHistory/Database/Actions';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: '*',
  tableName: t.addresses,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ changedValues, context }) => {
    const address = formatMoveHistoryFullAddress(changedValues);

    const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;

    let addressLabel = '';
    if (addressType === 'pickupAddress') {
      addressLabel = 'pickup_address';
    } else if (addressType === 'destinationAddress') {
      addressLabel = 'destination_address';
    }

    const newChangedValues = {
      shipment_type: context[0]?.shipment_type,
      ...changedValues,
    };

    newChangedValues[addressLabel] = address;

    return newChangedValues;
  },
};
