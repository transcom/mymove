import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateServiceMemberBackupContact from 'constants/MoveHistory/EventTemplates/UpdateServiceMemberBackupContact/updateServiceMemberBackupContact';

const BACKUP_CONTACT = {
  name: 'Ben Wyatt',
  email: 'benwyatt@example.com',
  phone: '555-555-2222',
};

describe('When a service members updates their profile', () => {
  const template = {
    action: 'UPDATE',
    eventName: 'updateServiceMemberBackupContact',
    tableName: 'backup_contacts',
    eventNameDisplay: 'Updated profile',
  };
  it('correctly matches the patch service member event template', () => {
    const result = getTemplate(template);
    expect(result).toMatchObject(updateServiceMemberBackupContact);
  });
  describe('it correctly renders the details component for an updated backup contact', () => {
    it.each([
      ['Backup contact name', ': Ben Wyatt'],
      ['Backup contact email', ': benwyatt@example.com'],
      ['Backup contact phone', ': 555-555-2222'],
    ])('displays the correct details value for %s', async (label, value) => {
      const historyRecord = { changedValues: BACKUP_CONTACT };
      const result = getTemplate(template);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

describe('When a service counselor updates a backup contact name', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateCustomer',
    tableName: 'backup_contacts',
    eventNameDisplay: 'Updated profile',
    changedValues: {
      name: 'Ron Swanson',
    },
  };

  it('correctly matches the updateCustomer event to the template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(updateServiceMemberBackupContact);
  });

  describe('Renders the correct details for updated customer name', () => {
    it.each([['Backup contact name', ': Ron Swanson']])('expect label %s to have value %s', async (label, value) => {
      const result = getTemplate(historyRecord);

      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
