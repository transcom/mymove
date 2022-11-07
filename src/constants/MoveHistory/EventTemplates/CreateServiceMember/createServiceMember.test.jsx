import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import t from 'constants/MoveHistory/Database/Tables';
import createServiceMember from 'constants/MoveHistory/EventTemplates/CreateServiceMember/createServiceMember';

describe('When a service members creates a profile', () => {
  const item = {
    action: 'INSERT',
    tableName: t.service_members,
    eventNameDisplay: 'Created Profile',
  };
  it('correctly matches the create service member event template', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(createServiceMember);
    expect(result.getEventNameDisplay()).toEqual('Created Profile');
  });
  it('correctly displays the details component', () => {
    const result = getTemplate(item);
    render(result.getDetails(item));
    expect(screen.getByText('-')).toBeInTheDocument();
  });
});
