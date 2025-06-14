import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateAllowance from 'constants/MoveHistory/EventTemplates/UpdateAllowances/updateAllowance';

describe('When a service counselor updates shipping allowances', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'updateAllowance',
    tableName: 'entitlements',
    eventNameDisplay: 'Updated allowances',
    changedValues: {
      authorized_weight: '4000',
      dependents_authorized: true,
      pro_gear_weight: '10',
      pro_gear_weight_spouse: '80',
      gun_safe_weight: '222',
      gun_safe: false,
      required_medical_equipment_weight: '100',
      storage_in_transit: '80',
    },
  };
  it('correctly matches the update allowances event template', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateAllowance);
  });
  describe('it correctly renders the details component', () => {
    it.each([
      ['Authorized weight', ': 4,000 lbs'],
      ['Storage in transit (SIT)', ': 80 days'],
      ['Dependents', ': Yes'],
      ['Pro-gear weight', ': 10 lbs'],
      ['Spouse pro-gear weight', ': 80 lbs'],
      ['Gun safe weight allowance', ': 222 lbs'],
      ['Gun safe authorized', ': No'],
      ['Required medical equipment', ': 100 lbs'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(item);
      render(result.getDetails(item));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
