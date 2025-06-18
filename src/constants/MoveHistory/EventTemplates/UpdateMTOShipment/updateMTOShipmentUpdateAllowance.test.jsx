import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipmentUpdateAllowance';

describe('when given an update to the allowance due to MTOShipment update, update MTO shipment history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOShipment',
    tableName: 'entitlements',
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
        shipment_locator: 'ABC123-01',
      },
    ],
    changedValues: { authorized_weight: 1650 },
  };

  const historyRecordForGunSafeAllowanceUpdate = {
    action: 'UPDATE',
    eventName: 'updateMTOShipment',
    tableName: 'entitlements',
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
        shipment_locator: 'ABC123-01',
      },
    ],
    changedValues: {
      gun_safe_weight: 222,
      gun_safe: true,
    },
  };

  it('correctly matches the update to the allowance, update MTO shipment event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly matches the event name when shipment is updated', () => {
    const template = getTemplate(historyRecord);
    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated shipment')).toBeInTheDocument();
  });

  it('displays the proper update MTO shipment record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Max billable weight')).toBeInTheDocument();
    expect(screen.getByText(': 1,650 lbs')).toBeInTheDocument();
  });

  it('correctly matches the event name when allowance is updated', () => {
    const template = getTemplate(historyRecordForGunSafeAllowanceUpdate);
    render(template.getEventNameDisplay(historyRecordForGunSafeAllowanceUpdate));
    expect(screen.getByText('Updated allowances')).toBeInTheDocument();
  });

  it('displays the proper updates to entitlement record with gun safe', () => {
    const template = getTemplate(historyRecordForGunSafeAllowanceUpdate);
    render(template.getDetails(historyRecordForGunSafeAllowanceUpdate));
    expect(screen.getByText('Gun safe weight allowance')).toBeInTheDocument();
    expect(screen.getByText(': 222 lbs')).toBeInTheDocument();
    expect(screen.getByText('Gun safe authorized')).toBeInTheDocument();
    expect(screen.getByText(': Yes')).toBeInTheDocument();
  });
});
