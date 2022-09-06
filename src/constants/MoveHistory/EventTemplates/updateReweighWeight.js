import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateReweigh,
  tableName: t.reweighs,
  detailsType: d.LABELED,
  getEventNameDisplay: () => `Updated shipment`,
  getDetailsLabeledDetails: ({ changedValues, context }) => {
    return {
      shipment_type: context[0]?.shipment_type,
      reweigh_weight: changedValues.weight,
    };
  },
};
