import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/updateMTOServiceItemAddress';

describe('when given a Update basic service item address history record', () => {
  const item = {
    action: a.UPDATE,
    changedValues: {
      city: 'San Diego',
      postal_code: '92134',
      street_address_1: '123 Test Street',
      street_address_2: '#19',
    },
    oldValues: {
      city: 'Beverly Hills',
      country: 'US',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    context: [
      {
        address_type: 'pickupAddress',
        shipment_type: 'HHG',
      },
    ],
    eventName: o.createMTOServiceItem,
    tableName: t.addresses,
  };
  it('correctly matches the update service item address event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay()).toEqual('Updated service item request');
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      shipment_type: 'HHG',
      city: 'San Diego',
      postal_code: '92134',
      street_address_1: '123 Test Street',
      street_address_2: '#19',
      pickup_address: '123 Test Street, #19, San Diego, CA 92134',
    });
  });
});
