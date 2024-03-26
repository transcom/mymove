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
    const newHistoryRecord = {
      ...historyRecord,
      changedValues: {
        // closeout_office_id: '123'
      },
      oldValues: {
        closeout_office_id: '123',
      },
      context: [{ closeout_office_name: 'this is the closeout_office_name' }],
    };
    const template = getTemplate(newHistoryRecord);
    const { baseElement } = render(template.getDetails(newHistoryRecord));
    expect(baseElement.textContent).toEqual('');
  });

  it('renders the proper message in the details column when updateCloseoutOffice is present ', () => {
    const newHistoryRecord = {
      ...historyRecord,
      changedValues: { closeout_office_id: '123' },
      context: [{ closeout_office_name: 'this is the closeout_office_name' }],
    };
    const template = getTemplate(newHistoryRecord);
    render(template.getDetails(newHistoryRecord));
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });
});
