import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItemCustomerContacts';

describe('when given a Create basic service item customer contacts history record', () => {
  const firstHistoryRecord = {
    action: 'INSERT',
    changedValues: {
      first_available_delivery_date: '2022-11-15',
      time_military: '1400Z',
      type: 'FIRST',
    },
    context: [
      {
        name: 'Domestic destination 1st day SIT',
        shipment_id_abbr: 'c3a9e',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'createMTOServiceItem',
    tableName: 'mto_service_item_customer_contacts',
  };

  const secondHistoryRecord = {
    action: 'INSERT',
    changedValues: {
      first_available_delivery_date: '2022-11-16',
      time_military: '1500Z',
      type: 'SECOND',
    },
    context: [
      {
        name: 'Domestic destination 1st day SIT',
        shipment_id_abbr: 'c3a9e',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'createMTOServiceItem',
    tableName: 'mto_service_item_customer_contacts',
  };

  it('correctly matches the create service item customer contacts event', () => {
    const template = getTemplate(firstHistoryRecord);
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Requested service item');
  });

  it('should display the correct shipment label', () => {
    const template = getTemplate(firstHistoryRecord);

    render(template.getDetails(firstHistoryRecord));
    expect(screen.getByText('HHG shipment #C3A9E, Domestic destination 1st day SIT'));
  });

  describe('when given a specific set of details', () => {
    it.each([
      ['First available delivery date', ': 15 Nov 2022', firstHistoryRecord],
      ['First available delivery time', ': 1400z', firstHistoryRecord],
      ['Second available delivery date', ': 16 Nov 2022', secondHistoryRecord],
      ['Second available delivery time', ': 1500Z', secondHistoryRecord],
    ])('for label %s it displays the proper details value %s', async (label, value, record) => {
      const template = getTemplate(record);

      render(template.getDetails(record));
      expect(screen.getByText(label, { exact: false })).toBeInTheDocument();
      expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
    });
  });
});
