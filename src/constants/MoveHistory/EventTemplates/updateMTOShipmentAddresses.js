import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { formatMoveHistoryFullAddress } from 'utils/formatters';
import AddressTypes from 'constants/MoveHistory/Database/AddressTypes';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.addresses,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ oldValues, changedValues, context }) => {
    let newChangedValues = {
      street_address_1: oldValues.street_address_1,
      street_address_2: oldValues.street_address_2,
      city: oldValues.city,
      state: oldValues.state,
      postal_code: oldValues.postal_code,
      ...changedValues,
    };

    const address = formatMoveHistoryFullAddress(newChangedValues);

    const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
    const addressLabel = AddressTypes[addressType];

    newChangedValues = {
      shipment_type: context[0]?.shipment_type,
      ...changedValues,
    };

    newChangedValues[addressLabel] = address;

    return newChangedValues;
  },
};
