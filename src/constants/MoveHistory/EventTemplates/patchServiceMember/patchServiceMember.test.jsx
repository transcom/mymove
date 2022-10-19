import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import patchServiceMember from 'constants/MoveHistory/EventTemplates/PatchServiceMember/patchServiceMember';

const PROFILE = {
  branch: {
    affiliation: 'AIR_FORCE',
    rank: 'E_4',
    edipi: '123456789',
  },
  name: { first_name: 'Leslie', last_name: 'Knope', middle_name: 'Barbara' },
  contact: {
    telephone: '555-444-5555',
    secondary_telephone: '555-555-5555',
    personal_email: 'leslieknope@example.com',
    phone_is_preferred: 'true',
    email_is_preferred: 'false',
  },
  backupContact: {
    name: 'Ben Wyatt',
    email: 'benwyatt@example.com',
    phone: '555-555-2222',
  },
};

describe('When a service members updates their profile', () => {
  const template = {
    action: 'UPDATE',
    eventName: o.patchServiceMember,
    tableName: t.service_members,
    eventNameDisplay: 'Updated profile',
  };
  it('correctly matches the patch service member event template', () => {
    const result = getTemplate(template);
    expect(result).toMatchObject(patchServiceMember);
  });
  describe('it correctly renders the details component for the branch form', () => {
    it.each([
      ['Branch', ': Air Force'],
      ['Rank', ': E-4'],
    ])('displays the correct details value for %s', async (label, value) => {
      const historyRecord = { changedValues: PROFILE.branch };
      const result = getTemplate(template);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  describe('it correctly renders the details component for the name form', () => {
    it.each([
      ['First name', ': Leslie'],
      ['Last name', ': Knope'],
      ['Middle name', ': Barbara'],
    ])('displays the correct details value for %s', async (label, value) => {
      const historyRecord = { changedValues: PROFILE.name };
      const result = getTemplate(template);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  describe('it correctly renders the details component for the contact information form', () => {
    it.each([
      ['Telephone', ': 555-444-5555'],
      ['Secondary telephone', ': 555-555-5555'],
      ['Personal email', ': leslieknope@example.com'],
      ['Phone preferred', ': true'],
      ['Email preferred', ': false'],
    ])('displays the correct details value for %s', async (label, value) => {
      const historyRecord = { changedValues: PROFILE.contact };
      const result = getTemplate(template);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  describe('it correctly renders the details component for duty location', () => {
    const historyRecord = {
      changedValues: { duty_location_id: '2222' },
      context: [{ current_duty_location_name: 'Fort Hood' }],
    };
    const result = getTemplate(template);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('Current duty location name')).toBeInTheDocument();
    expect(screen.getByText(': Fort Hood')).toBeInTheDocument();
  });
});
