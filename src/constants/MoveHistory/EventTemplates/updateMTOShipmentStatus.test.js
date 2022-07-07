import t from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/updateMTOShipmentStatus';

describe('when Prime user cancels a shipment', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipmentStatus,
    tableName: 'mto_shipments',
    detailsType: d.LABELED,
  };
  const historyRecord = {
    oldValues: {
      shipment_type: 'HHG shipment',
    },
    changedValues: {
      status: 'CANCELED',
    },
  };
  it('correctly displays shipment type and status', () => {
    const result = t(item);
    result.getDetailsLabeledDetails(historyRecord);
    expect(result).toMatchObject(e);
  });
});
