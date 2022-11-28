import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateCustomer from 'constants/MoveHistory/EventTemplates/UpdateCustomer/updateCustomer';

describe('When a service counselor updates a customer name', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateCustomer',
    tableName: 'service_members',
    eventNameDisplay: 'Updated profile',
    changedValues: {
      first_name: 'Chuck',
      last_name: 'Bartowski',
    },
  };

  it('correctly matches the updateCustomer event to the template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(updateCustomer);
  });

  describe('Renders the correct details for updated customer name', () => {
    it.each([
      ['First name', ': Chuck'],
      ['Last name', ': Bartowski'],
    ])('expect label %s to have value %s', async (label, value) => {
      const result = getTemplate(historyRecord);

      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
