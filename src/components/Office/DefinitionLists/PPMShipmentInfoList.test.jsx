import React from 'react';
import { render, screen } from '@testing-library/react';

import PPMShipmentInfoList from './PPMShipmentInfoList';

import affiliation from 'content/serviceMemberAgencies';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const renderWithPermissions = (shipment) => {
  render(
    <MockProviders permissions={[permissionTypes.viewCloseoutOffice]}>
      <PPMShipmentInfoList shipment={shipment} />
    </MockProviders>,
  );
};

describe('PPMShipmentInfoList', () => {
  it('renders closeout display for Marines', () => {
    renderWithPermissions({ agency: affiliation.MARINES });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('TVCB');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display for Navy', () => {
    renderWithPermissions({ agency: affiliation.NAVY });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('NAVY');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display Coast guard', () => {
    renderWithPermissions({ agency: affiliation.COAST_GUARD });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('USCG');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display for Army and Air Force', () => {
    renderWithPermissions({ closeoutOffice: 'Test office' });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('Test office');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display when there is no closeout office', () => {
    renderWithPermissions({ closeoutOffice: '-' });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('-');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });
});
