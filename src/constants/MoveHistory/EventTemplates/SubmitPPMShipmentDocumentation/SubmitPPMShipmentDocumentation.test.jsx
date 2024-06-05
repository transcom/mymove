import { render, screen } from '@testing-library/react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given a PPM shipment update', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: o.submitPPMShipmentDocumentation,
    tableName: t.ppm_shipments,
    changedValues: {
      status: 'NEEDS_CLOSEOUT',
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: '12992',
      },
    ],
  };

  it('displays the correct label for shipment', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
  });

  it('displays that the shipment was submitted', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': NEEDS_CLOSEOUT')).toBeInTheDocument();
  });
});
