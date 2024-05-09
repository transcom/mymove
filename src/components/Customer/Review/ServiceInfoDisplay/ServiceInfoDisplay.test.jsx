import React from 'react';
import { screen } from '@testing-library/react';

import ServiceInfoDisplay from './ServiceInfoDisplay';

import { renderWithRouter } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('ServiceInfoDisplay component', () => {
  const testProps = {
    firstName: 'Jason',
    lastName: 'Ash',
    affiliation: 'Air Force',
    edipi: '9999999999',
    emplid: '1234567',
    originTransportationOfficeName: 'Buckley AFB',
    originTransportationOfficePhone: '555-555-5555',
  };

  it('renders the data', async () => {
    renderWithRouter(<ServiceInfoDisplay {...testProps} />);

    const mainHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(mainHeader).toBeInTheDocument();

    const name = screen.getByText('Name');

    expect(name).toBeInTheDocument();

    expect(name.nextElementSibling.textContent).toBe(`${testProps.firstName} ${testProps.lastName}`);

    const branch = screen.getByText('Branch');

    expect(branch).toBeInTheDocument();

    expect(branch.nextElementSibling.textContent).toBe(testProps.affiliation);

    const dodId = screen.getByText('DoD ID#');

    expect(dodId).toBeInTheDocument();

    expect(dodId.nextElementSibling.textContent).toBe(testProps.edipi);

    const editLink = screen.getByText('Edit');

    expect(editLink).toBeInTheDocument();
  });

  it('renders who to contact when the service info is no longer editable and it should notify the transportation office', async () => {
    renderWithRouter(<ServiceInfoDisplay {...testProps} isEditable={false} showMessage />);

    expect(screen.queryByText('Edit')).toBeNull();

    const whoToContact = screen.getByText(
      `To change information in this section, contact the ${testProps.originTransportationOfficeName} transportation office at ${testProps.originTransportationOfficePhone}.`,
    );

    expect(whoToContact).toBeInTheDocument();
  });

  it('renders a non editable service info display wuth no message', () => {
    renderWithRouter(<ServiceInfoDisplay {...testProps} isEditable={false} />);

    expect(screen.queryByText('Edit')).toBeNull();

    expect(
      screen.queryByText(
        `To change information in this section, contact the ${testProps.originTransportationOfficeName} transportation office at ${testProps.originTransportationOfficePhone}.`,
      ),
    ).toBeNull();
  });

  it('Coast Guard Customers', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    testProps.affiliation = 'Coast Guard';
    renderWithRouter(<ServiceInfoDisplay {...testProps} />);

    const mainHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(mainHeader).toBeInTheDocument();

    const name = screen.getByText('Name');

    expect(name).toBeInTheDocument();

    expect(name.nextElementSibling.textContent).toBe(`${testProps.firstName} ${testProps.lastName}`);

    const branch = screen.getByText('Branch');

    expect(branch).toBeInTheDocument();

    expect(branch.nextElementSibling.textContent).toBe(testProps.affiliation);

    const dodId = screen.getByText('DoD ID#');

    expect(dodId).toBeInTheDocument();

    expect(dodId.nextElementSibling.textContent).toBe(testProps.edipi);

    const emplid = screen.getByText('EMPLID');

    expect(emplid).toBeInTheDocument();

    expect(emplid.nextElementSibling.textContent).toBe(testProps.emplid);

    const editLink = screen.getByText('Edit');

    expect(editLink).toBeInTheDocument();
  });
});
