import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/proofOfServiceDocUpload';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a proof of service document upload history record', () => {
  const item = {
    action: a.INSERT,
    eventName: 'createUpload',
    tableName: t.proof_of_service_docs,
    context: [{ payment_request_number: '1234-5678-1' }],
  };
  it('correctly matches the proof of service document upload event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual(
      'Proof of service document uploaded for payment request 1234-5678-1',
    );
  });
});
