import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ContactInfoDisplay from './ContactInfoDisplay';

import { renderWithRouter } from 'testUtils';

describe('ContactInfoDisplay component', () => {
  const testProps = {
    telephone: '703-555-4578',
    personalEmail: 'test@example.com',
    emailIsPreferred: true,
    residentialAddress: {
      street_address_1: '1292 Orchard Terrace',
      street_address_2: 'Building C, Unit 10',
      city: 'El Paso',
      state: 'TX',
      postal_code: '79912',
    },
    backupMailingAddress: {
      street_address_1: '448 Washington Blvd NE',
      street_address_2: '',
      city: 'El Paso',
      state: 'TX',
      postal_code: '79936',
    },
    backupContact: {
      name: 'Gabriela Sáenz Perez',
      telephone: '206-555-8989',
      email: 'gsp@example.com',
    },
    editURL: '/moves/review/edit-profile',
  };

  it('renders the data', async () => {
    renderWithRouter(<ContactInfoDisplay {...testProps} />);

    const mainHeader = await screen.findByRole('heading', { name: 'Contact info', level: 2 });

    expect(mainHeader).toBeInTheDocument();

    const phoneTerm = screen.getByText('Best contact phone');

    expect(phoneTerm).toBeInTheDocument();

    expect(phoneTerm.nextElementSibling.textContent).toBe(testProps.telephone);

    const emailTerm = screen.getByText('Personal email');

    expect(emailTerm).toBeInTheDocument();

    expect(emailTerm.nextElementSibling.textContent).toBe(testProps.personalEmail);

    const addressTerm = screen.getByText('Current mailing address');

    expect(addressTerm).toBeInTheDocument();

    Object.values(testProps.residentialAddress).forEach((value) => {
      expect(addressTerm.nextElementSibling.textContent).toContain(value);
    });

    const backupAddressTerm = screen.getByText('Backup mailing address');

    expect(backupAddressTerm).toBeInTheDocument();

    Object.values(testProps.backupMailingAddress).forEach((value) => {
      expect(backupAddressTerm.nextElementSibling.textContent).toContain(value);
    });

    const backupHeader = screen.getByRole('heading', { name: 'Backup contact', level: 3 });

    expect(backupHeader).toBeInTheDocument();

    const backupNameTerm = screen.getByText('Name');

    expect(backupNameTerm).toBeInTheDocument();

    expect(backupNameTerm.nextElementSibling.textContent).toBe(testProps.backupContact.name);

    const backupEmailTerm = screen.getAllByText('Email')[1];

    expect(backupEmailTerm).toBeInTheDocument();

    expect(backupEmailTerm.nextElementSibling.textContent).toBe(testProps.backupContact.email);

    const backupPhoneTerm = screen.getByText('Phone');

    expect(backupPhoneTerm).toBeInTheDocument();

    expect(backupPhoneTerm.nextElementSibling.textContent).toBe(testProps.backupContact.telephone);
  });

  it.each([
    ['', '–'],
    ['703-555-9999', '703-555-9999'],
  ])('Shows alt phone (%s) as expected (%s)', async (secondaryTelephone, expectedDisplay) => {
    const contactProps = { ...testProps, secondaryTelephone };

    renderWithRouter(<ContactInfoDisplay {...contactProps} />);

    const altPhoneTerm = await screen.findByText('Alt. phone');

    expect(altPhoneTerm).toBeInTheDocument();

    expect(altPhoneTerm.nextElementSibling.textContent).toBe(expectedDisplay);
  });

  it.each([
    [true, false, 'Phone'],
    [false, true, 'Email'],
    [true, true, 'Phone, Email'],
  ])(
    'Shows preferred contact (Phone: %s | Email: %s) as expected: %s',
    async (phoneIsPreferred, emailIsPreferred, expectedDisplay) => {
      const contactProps = { ...testProps, phoneIsPreferred, emailIsPreferred };

      renderWithRouter(<ContactInfoDisplay {...contactProps} />);

      const contactMethodTerm = await screen.findByText('Preferred contact method');

      expect(contactMethodTerm).toBeInTheDocument();

      expect(contactMethodTerm.nextElementSibling.textContent).toBe(expectedDisplay);
    },
  );

  it('Goes to editURL when Edit link is clicked', async () => {
    const { history } = renderWithRouter(<ContactInfoDisplay {...testProps} />);

    const editLink = await screen.findByRole('link');

    expect(editLink).toBeInTheDocument();

    userEvent.click(editLink);

    await waitFor(() => {
      expect(history.location.pathname).toEqual(testProps.editURL);
    });
  });
});
