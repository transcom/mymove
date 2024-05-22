import { render, screen } from '@testing-library/react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given a PPM shipment update', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: o.submitMoveForApproval,
    tableName: t.mto_shipments,
    changedValues: {
      status: 'SUBMITTED',
    },
    context: [
      {
        shipment_type: 'HHG',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: '12992',
      },
    ],
  };

  it('displays the correct label for shipment', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #RQ38D4-01')).toBeInTheDocument();
  });

  it('displays that the shipment was submitted', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': SUBMITTED')).toBeInTheDocument();
  });

  it('displays event name', () => {
    const result = getTemplate(historyRecord);
    render(result.getEventNameDisplay());

    expect(screen.getByText('Submitted Move for Approval')).toBeInTheDocument();
  });
});
