import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateUpload/proofOfServiceDocUpload';

describe('when given a proof of service document upload history record', () => {
  const historyRecord = {
    action: 'INSERT',
    eventName: 'createUpload',
    tableName: 'proof_of_service_docs',
    context: [{ payment_request_number: '1234-5678-1' }],
  };
  it('correctly matches the proof of service document upload event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper details message with the correct request number', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Proof of service document uploaded for payment request 1234-5678-1')).toBeInTheDocument();
  });
});
