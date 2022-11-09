import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import createServiceMemberBackupContact from 'constants/MoveHistory/EventTemplates/CreateServiceMemberBackupContact/createServiceMemberBackupContact';

const BACKUP_CONTACT = {
  name: 'Tom Haverford',
  email: 'tommyHavie@example.com',
  phone: '555-555-5555',
};

describe('When a service members updates their profile', () => {
  const template = {
    action: 'INSERT',
    eventName: o.createServiceMemberBackupContact,
    tableName: t.backup_contacts,
    eventNameDisplay: 'Updated profile',
  };
  it('correctly matches the patch service member event template', () => {
    const result = getTemplate(template);
    expect(result).toMatchObject(createServiceMemberBackupContact);
  });
  describe('it correctly renders the details component for an updated backup contact', () => {
    it.each([
      ['Backup contact name', ': Tom Haverford'],
      ['Backup contact email', ': tommyHavie@example.com'],
      ['Backup contact phone', ': 555-555-5555'],
    ])('displays the correct details value for %s', async (label, value) => {
      const historyRecord = { changedValues: BACKUP_CONTACT };
      const result = getTemplate(template);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
