import React from 'react';
import { screen } from '@testing-library/react';
import { cloneDeep } from 'lodash';

import OktaInfoDisplay from './OktaInfoDisplay';

import { renderWithRouter } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

describe('OktaInfoDisplay component', () => {
  const testProps = {
    isEditable: true,
    oktaUsername: 'dummy@okta.mil',
    oktaEmail: 'dummy@okta.mil',
    oktaFirstName: 'Dummy',
    oktaLastName: 'Dumber',
    oktaEdipi: '1234567890',
    editURL: '/moves/review/edit-okta-profile',
  };

  it('renders the data', async () => {
    renderWithRouter(<OktaInfoDisplay {...testProps} />);

    const oktaUsername = screen.getByText('Username');
    expect(oktaUsername).toBeInTheDocument();
    expect(oktaUsername.nextElementSibling.textContent).toBe(testProps.oktaUsername);

    const oktaEmail = screen.getByText('Email');
    expect(oktaEmail).toBeInTheDocument();
    expect(oktaEmail.nextElementSibling.textContent).toBe(testProps.oktaEmail);

    const oktaFirstName = screen.getByText('First Name');
    expect(oktaFirstName).toBeInTheDocument();
    expect(oktaFirstName.nextElementSibling.textContent).toBe(testProps.oktaFirstName);

    const oktaLastName = screen.getByText('Last Name');
    expect(oktaLastName).toBeInTheDocument();
    expect(oktaLastName.nextElementSibling.textContent).toBe(testProps.oktaLastName);

    const oktaEdipi = screen.getByText('DoD ID Number | EDIPI');
    expect(oktaEdipi).toBeInTheDocument();
    expect(oktaEdipi.nextElementSibling.textContent).toBe(testProps.oktaEdipi);
  });

  it('Goes to editURL when Edit link is clicked', async () => {
    renderWithRouter(<OktaInfoDisplay {...testProps} />, { path: '/moves/review/edit-okta-profile' });

    const editLink = screen.getByRole('link', { name: 'Edit' });

    expect(editLink).toBeInTheDocument();

    expect(editLink.href).toContain(testProps.editURL);
  });

  it('Disables edit link when isEditable prop is false', async () => {
    const testPropsNotEditable = cloneDeep(testProps);
    testPropsNotEditable.isEditable = false;
    renderWithRouter(<OktaInfoDisplay {...testPropsNotEditable} />, { path: '/moves/review/edit-okta-profile' });

    const editLink = screen.queryByRole('link', { name: 'Edit' });

    expect(editLink).toBeNull();
  });
});
