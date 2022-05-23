import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/updateMTOShipment';

describe('when given an mto shipment update with mto shipment table history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipment,
    tableName: 'mto_shipments',
    detailsType: d.LABELED,
    changedValues: {
      destination_address_type: 'HOME_OF_SELECTION',
      requested_delivery_date: '2020-04-14',
      requested_pickup_date: '2020-03-23',
    },
  };
  it('correctly matches the Update mto shipment event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
});
