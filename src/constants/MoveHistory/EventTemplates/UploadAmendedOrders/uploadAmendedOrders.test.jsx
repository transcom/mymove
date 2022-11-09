import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UploadAmendedOrders/uploadAmendedOrders';

describe('when given an upload amended orders history record', () => {
  const item = {
    action: 'INSERT',
    eventName: 'uploadAmendedOrders',
    tableName: 'user_uploads',
    context: [{ filename: 'test.txt', upload_type: 'amendedOrders' }],
  };
  it('correctly matches the upload amended orders event', () => {
    const details = `Uploaded amended orders document ${item.context[0]?.filename}`;
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    expect(screen.getByText(details)).toBeInTheDocument();
  });
});
