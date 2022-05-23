import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import e from 'constants/MoveHistory/EventTemplates/createMTOShipmentAddresses';

describe('when given an mto shipment insert with address table history record', () => {
  const item = {
    action: 'INSERT',
    eventName: '',
    tableName: 'addresses',
    detailsType: d.LABELED,
    changedValues: {
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
      state: 'CA',
    },
    context: [{ shipment_type: 'HHG', address_type: 'pickupAddress' }],
  };

  it('correctly matches the insert mto shipment address event for pickup addresses', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the adddresses correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        context: item.context,
      }),
    ).toEqual({
      pickup_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
      state: 'CA',
      shipment_type: 'HHG',
    });
  });

  it('correctly matches the insert mto shipment address event for destination addresses', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the adddresses correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        context: [{ shipment_type: 'HHG', address_type: 'destinationAddress' }],
      }),
    ).toEqual({
      destination_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
      state: 'CA',
      shipment_type: 'HHG',
    });
  });
});
