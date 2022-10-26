import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/DeleteUpload/deleteUpload';

describe('when given a delete upload history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'deleteUpload',
    tableName: 'user_uploads',
    context: [{ filename: 'test.txt' }],
  };
  it('correctly matches the delete upload event for orders deletions', () => {
    const details = `Deleted orders document test.txt`;
    item.context[0].upload_type = 'orders';
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    render(result.getEventNameDisplay(item));
    expect(screen.getByText(details)).toBeInTheDocument();
    expect(screen.getByText('Updated orders')).toBeInTheDocument();
  });
  it('correctly matches the delete upload event for amended orders deletions', () => {
    const details = `Deleted amended orders document test.txt`;
    item.context[0].upload_type = 'amendedOrders';
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    render(result.getEventNameDisplay(item));
    expect(screen.getByText(details)).toBeInTheDocument();
    expect(screen.getByText('Updated orders')).toBeInTheDocument();
  });
  it('correctly matches the default display when the type is not recognized', () => {
    item.context[0].upload_type = 'default';
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    render(result.getEventNameDisplay(item));
    expect(screen.getByText('-')).toBeInTheDocument();
    expect(screen.getByText('Updated move')).toBeInTheDocument();
  });
});
