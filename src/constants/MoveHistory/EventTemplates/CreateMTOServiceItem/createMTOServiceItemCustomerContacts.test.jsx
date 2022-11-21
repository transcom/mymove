import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItemCustomerContacts';

describe('when given a Create basic service item customer contacts history record', () => {
  const historyRecord = {
    action: 'INSERT',
    changedValues: {
      first_available_delivery_date: '2022-06-30T00:00:00+00:00',
      time_military: '1500Z',
      type: 'SECOND',
    },
    eventName: 'createMTOServiceItem',
    tableName: 'mto_service_item_customer_contacts',
  };
  const template = getTemplate(historyRecord);
  it('correctly matches the create service item customer contacts event', () => {
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Requested service item');
  });
  describe('when given a specific set of details', () => {
    it.each([
      ['first_available_delivery_date', '30 Jun 2022'],
      ['second_available_delivery_time', '1500Z'],
    ])('for label %s it displays the proper details value %s', async (label, value) => {
      render(template.getDetails(historyRecord));
      expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
    });
  });
});
