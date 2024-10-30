import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when customer uploads additional documents', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'uploadAdditionalDocuments',
    tableName: 'moves',
  };
  it('correctly matches the Upload additional documents event', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Uploaded additional document')).toBeInTheDocument();
  });
});
