import { formatMoveHistoryFullAddress } from 'utils/formatters';
import a from 'constants/MoveHistory/Database/Actions';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.createMTOServiceItem,
  tableName: t.addresses,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated service item request',
  getDetailsLabeledDetails: ({ oldValues, changedValues, context }) => {
    const address = formatMoveHistoryFullAddress({ ...oldValues, ...changedValues });

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
