import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateAllowanceByCounselor from 'constants/MoveHistory/EventTemplates/UpdateAllowances/updateAllowanceByCounselor';

describe('When a service counselor updates shipping allowances', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'counselingUpdateAllowance',
    tableName: 'entitlements',
    eventNameDisplay: 'Updated allowances',
    changedValues: {
      authorized_weight: '8000',
      dependents_authorized: true,
      pro_gear_weight: '100',
      pro_gear_weight_spouse: '85',
      required_medical_equipment_weight: '10',
      storage_in_transit: '80',
    },
  };
  it('correctly matches the update allowances event template', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(updateAllowanceByCounselor);
    expect(result.getEventNameDisplay()).toMatch(historyRecord.eventNameDisplay);
  });
  describe('it correctly renders the details component', () => {
    it.each([
      ['Authorized weight', ': 8,000 lbs'],
      ['Dependents', ': Yes'],
      ['Pro-gear', ': 100 lbs'],
      ['Spouse pro-gear', ': 85 lbs'],
      ['RME', ': 10 lbs'],
      ['Storage in transit (SIT)', ': 80 days'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
