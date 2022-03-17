import React from 'react';
import { screen } from '@testing-library/react';

import ServiceInfoDisplay from './ServiceInfoDisplay';

import { renderWithRouter } from 'testUtils';

describe('ServiceInfoDisplay component', () => {
  const testProps = {
    firstName: 'Jason',
    lastName: 'Ash',
    affiliation: 'Air Force',
    rank: 'E-5',
    edipi: '9999999999',
    originDutyLocationName: 'Buckley AFB',
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

    const rank = screen.getByText('Rank');

    expect(rank).toBeInTheDocument();

    expect(rank.nextElementSibling.textContent).toBe(testProps.rank);

    const dodId = screen.getByText('DoD ID#');

    expect(dodId).toBeInTheDocument();

    expect(dodId.nextElementSibling.textContent).toBe(testProps.edipi);

    const currentDutyStation = screen.getByText('Current duty location');

    expect(currentDutyStation).toBeInTheDocument();

    expect(currentDutyStation.nextElementSibling.textContent).toBe(testProps.originDutyLocationName);

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
});
