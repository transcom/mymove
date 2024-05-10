import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given a deleted trip history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      final_incentive: 46855239,
    },
    context: [
      {
        shipment_id_abbr: '7f559',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.deleteWeightTicket,
    tableName: t.ppm_shipments,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated shipment')).toBeInTheDocument();
  });

  it('displays details of shipment type, shipment ID', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
  });

  it('displays incentive details', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Final incentive')).toBeInTheDocument();
    expect(screen.getByText(': $468,552.39')).toBeInTheDocument();
  });
});
