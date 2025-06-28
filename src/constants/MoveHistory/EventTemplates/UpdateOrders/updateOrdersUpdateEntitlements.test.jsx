import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateOrders/updateOrdersUpdateEntitlements';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a update to an order, update allowance history record', () => {
  const historyRecord = {
    action: a.INSERT,
    eventName: o.updateOrders,
    tableName: t.entitlements,
    changedValues: {
      accompanied_tour: null,
      authorized_weight: 13000,
      dependents_authorized: true,
      gun_safe: false,
      gun_safe_weight: 200,
      pro_gear_weight: 2000,
      pro_gear_weight_spouse: 500,
      storage_in_transit: 90,
    },
  };

  it('correctly matches the update order, update allowance event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper update allowances record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Authorized weight')).toBeInTheDocument();
    expect(screen.getByText(': 13,000 lbs')).toBeInTheDocument();
    expect(screen.getByText('Dependents')).toBeInTheDocument();
    expect(screen.getByText(': Yes')).toBeInTheDocument();
    expect(screen.getByText('Gun safe')).toBeInTheDocument();
    expect(screen.getByText(': No')).toBeInTheDocument();
    expect(screen.getByText('Gun safe weight')).toBeInTheDocument();
    expect(screen.getByText(': 200 lbs')).toBeInTheDocument();
    expect(screen.getByText('Pro-gear weight')).toBeInTheDocument();
    expect(screen.getByText(': 2,000 lbs')).toBeInTheDocument();
    expect(screen.getByText('Spouse pro-gear weight')).toBeInTheDocument();
    expect(screen.getByText(': 500 lbs')).toBeInTheDocument();
    expect(screen.getByText('Storage in transit (SIT)')).toBeInTheDocument();
    expect(screen.getByText(': 90 days')).toBeInTheDocument();
  });
});
