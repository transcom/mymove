import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateCloseoutOffice/updateCloseoutOffice';

describe('when given an updateCloseoutOffice history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateCloseoutOffice',
    tableName: 'moves',
  };

  it('correctly matches updateCloseoutOffice template', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });

  it('renders the default details in the details column when updateCloseoutOffice key is not present ', () => {
    const template = getTemplate(historyRecord);
    const ren = render(template.getDetails(historyRecord));
    ren.debug();
    expect(screen.getByText('-')).toBeInTheDocument();
  });

  it('renders the proper message in the details column when updateCloseoutOffice is present ', () => {
    const newHistoryRecord = {
      ...historyRecord,
      context: [{ closeout_office_name: 'this is the closeout_office_name' }],
    };
    const template = getTemplate(newHistoryRecord);
    const ren = render(template.getDetails(newHistoryRecord));
    ren.debug();
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });
});
