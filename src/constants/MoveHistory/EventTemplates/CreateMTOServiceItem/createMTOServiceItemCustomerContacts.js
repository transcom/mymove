import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.createMTOServiceItem,
  tableName: t.mto_service_item_customer_contacts,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Requested service item',
  getDetailsLabeledDetails: ({ changedValues }) => {
    const { type, time_military: timeMilitary } = changedValues;

    const deliveryTimeOrder = type === 'FIRST' ? 'first_available_delivery_time' : 'second_available_delivery_time';

    const newChangedValues = {
      ...changedValues,
    };

    newChangedValues[deliveryTimeOrder] = timeMilitary;

    return newChangedValues;
  },
};
