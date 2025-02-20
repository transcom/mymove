import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given an UpdateMTOServiceItem history record with pricing/weight changes', () => {
  const context = [
    {
      name: 'International shipping & linehaul',
      shipment_type: 'HHG',
      shipment_locator: 'RQ38D4-01',
      shipment_id_abbr: 'a1b2c',
    },
  ];

  it('correctly matches the service item price update event', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        pricing_estimate: 150000,
      },
      context,
      eventName: o.updateMTOShipment,
      tableName: t.mto_service_items,
    };

    const template = getTemplate(historyRecord);
    expect(template.getEventNameDisplay(historyRecord)).toEqual('Service item estimated price updated');

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Service item')).toBeInTheDocument();
    expect(screen.getByText(/International shipping & linehaul/)).toBeInTheDocument();
    expect(screen.getByText('Estimated Price')).toBeInTheDocument();
    expect(screen.getByText(/\$1,500\.00/)).toBeInTheDocument();
  });

  it('correctly matches the service item weight update event', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        estimated_weight: 1000,
      },
      context,
      eventName: o.updateMTOShipment,
      tableName: t.mto_service_items,
    };

    const template = getTemplate(historyRecord);
    expect(template.getEventNameDisplay(historyRecord)).toEqual('Service item estimated weight updated');

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Service item')).toBeInTheDocument();
    expect(screen.getByText(/International shipping & linehaul/)).toBeInTheDocument();
    expect(screen.getByText('Estimated weight')).toBeInTheDocument();
    expect(screen.getByText(/1,000 lbs/)).toBeInTheDocument();
  });

  it('returns "Service item updated" for unknown changes', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        unknownField: 'some value',
      },
      context,
      eventName: o.updateMTOShipment,
      tableName: t.mto_service_items,
    };

    const template = getTemplate(historyRecord);
    expect(template.getEventNameDisplay(historyRecord)).toEqual('Service item updated');
    expect(template.getDetails(historyRecord)).toBeNull();
  });
});
