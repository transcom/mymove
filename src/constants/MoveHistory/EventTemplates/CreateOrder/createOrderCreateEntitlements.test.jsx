import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateOrder/createOrderCreateEntitlements';

describe('When given a created order event for the entitlements table from the office side', () => {
  const item = {
    action: 'INSERT',
    eventName: 'createOrder',
    tableName: 'entitlements',
    eventNameDisplay: 'Created allowances',
    changedValues: {
      authorized_weight: 8000,
      dependents_authorized: true,
      storage_in_transit: 90,
    },
  };
  it('correctly matches the created orders template', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
  describe('When given a specific set of details', () => {
    it.each([
      ['Authorized weight', ': 8,000 lbs'],
      ['Storage in transit (SIT)', ': 90 days'],
      ['Dependents', ': Yes'],
    ])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(item);
      render(result.getDetails(item));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
