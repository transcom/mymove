import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.createMTOServiceItem,
  tableName: t.mto_service_items,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Requested service item',
  getDetailsLabeledDetails: ({ changedValues, context }) => {
    const newChangedValues = {
      shipment_type: context[0]?.shipment_type,
      service_item_name: context[0]?.name,
      ...changedValues,
    };

    return newChangedValues;
  },
};
