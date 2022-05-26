import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateMTOShipmentAddresses';

describe('when given an mto shipment update with address table history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipment,
    tableName: 'addresses',
    detailsType: d.LABELED,
    changedValues: {
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
    },
    oldValues: {
      city: 'Beverly Hills',
      postal_code: '90211',
      state: 'CA',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
    },
    context: [{ shipment_type: 'HHG', address_type: 'pickupAddress' }],
  };

  it('correctly matches the Update mto shipment address event for pickup addresses', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the adddresses correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        oldValues: item.oldValues,
        context: item.context,
      }),
    ).toEqual({
      pickup_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
      shipment_type: 'HHG',
    });
  });

  it('correctly matches the Update mto shipment address event for destination addresses', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the adddresses correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        oldValues: item.oldValues,
        context: [{ shipment_type: 'HHG', address_type: 'destinationAddress' }],
      }),
    ).toEqual({
      destination_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
      shipment_type: 'HHG',
    });
  });
});
