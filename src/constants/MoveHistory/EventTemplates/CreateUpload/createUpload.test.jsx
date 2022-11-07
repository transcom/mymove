import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateUpload/createUpload';

describe('when given a create upload history record', () => {
  const item = {
    action: 'INSERT',
    eventName: 'createUpload',
    tableName: 'user_uploads',
    context: [{ filename: 'test.txt', upload_type: 'orders' }],
  };
  it('correctly matches the create upload event', () => {
    const details = `Uploaded orders document ${item.context[0]?.filename}`;
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    expect(screen.getByText(details)).toBeInTheDocument();
  });
});
