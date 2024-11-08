import { render, screen } from '@testing-library/react';

import addAppealToViolation from 'constants/MoveHistory/EventTemplates/AddAppeal/addAppealToViolation';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given an Appeal Decision on Violation history record', () => {
  const historyRecord = {
    action: 'INSERT',
    changedValues: {
      evaluation_report_id: '12345',
      remarks: 'Appeal remark',
      appeal_status: 'SUSTAINED',
    },
    eventName: 'addAppealToViolation',
    tableName: 'gsr_appeals',
    context: [
      {
        evaluation_report_type: 'SHIPMENT',
      },
    ],
  };

  it('correctly matches to the Appeal Decision on Violation template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(addAppealToViolation);
  });

  it('displays the proper value in the details field', async () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));

    const reportIdElement = screen.getByTestId('violationAppealInfo');
    expect(reportIdElement).toBeInTheDocument();
  });
});
