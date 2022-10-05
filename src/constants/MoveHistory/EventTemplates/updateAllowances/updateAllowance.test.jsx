import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateAllowance from 'constants/MoveHistory/EventTemplates/updateAllowances/updateAllowance';

describe('When a service counselor updates shipping allowances', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateAllowance,
    tableName: t.entitlements,
    eventNameDisplay: 'Updated allowances',
    changedValues: {
      authorized_weight: '4000',
      dependents_authorized: 'false',
      pro_gear_weight: '10',
      pro_gear_weight_spouse: '80',
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
      ['Authorized weight', ': 4000 lbs'],
      ['Storage in transit (SIT)', ': 80 days'],
      ['Dependents', ': false'],
      ['Pro-gear', ': 10 lbs'],
      ['Spouse pro-gear', ': 80 lbs'],
      ['RME', ': 100 lbs'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(item);
      render(result.getDetails(item));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
