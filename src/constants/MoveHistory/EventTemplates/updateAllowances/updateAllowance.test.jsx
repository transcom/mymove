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
  it('correctly renders the details component', () => {
    const result = getTemplate(item);
    render(result.getDetails(item));
    expect(screen.getByText('Authorized weight')).toBeInTheDocument();
    expect(screen.getByText(': 4000 lbs')).toBeInTheDocument();
    expect(screen.getByText('Storage in transit (SIT)')).toBeInTheDocument();
    expect(screen.getByText(': 80 days')).toBeInTheDocument();
    expect(screen.getByText('Dependents')).toBeInTheDocument();
    expect(screen.getByText(': false')).toBeInTheDocument();
    expect(screen.getByText('Pro-gear')).toBeInTheDocument();
    expect(screen.getByText(': 10 lbs')).toBeInTheDocument();
    expect(screen.getByText('Spouse pro-gear')).toBeInTheDocument();
    expect(screen.getByText(': 80 lbs')).toBeInTheDocument();
    expect(screen.getByText('RME')).toBeInTheDocument();
    expect(screen.getByText(': 100 lbs')).toBeInTheDocument();
  });
});
