import React from 'react';
import { render, screen } from '@testing-library/react';

import AllowancesList from './AllowancesList';

const info = {
  branch: 'NAVY',
  rank: 'E_6',
  weightAllowance: 12000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
};

describe('AllowancesList', () => {
  it('renders formatted branch and rank', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('Navy, E-6')).toBeInTheDocument();
  });

  it('renders formatted weight allowance', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('12,000 lbs')).toBeInTheDocument();
  });

  it('renders formatted authorized weight', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('11,000 lbs')).toBeInTheDocument();
  });

  it('renders storage in transit', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('90 days')).toBeInTheDocument();
  });

  it('renders authorized dependents', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByTestId('dependents').textContent).toEqual('Authorized');
  });

  it('renders unauthorized dependents', () => {
    const withUnauthorizedDependents = { ...info, dependents: false };
    render(<AllowancesList info={withUnauthorizedDependents} />);
    expect(screen.getByTestId('dependents').textContent).toEqual('Unauthorized');
  });

  it('renders formatted pro-gear', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('2,000 lbs')).toBeInTheDocument();
  });

  it('renders formatted spouse pro-gear', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('500 lbs')).toBeInTheDocument();
  });

  it('renders formatted rme', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('1,000 lbs')).toBeInTheDocument();
  });

  it('renders authorized ocie', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByTestId('ocie').textContent).toEqual('Authorized');
  });

  it('renders unauthorized ocie', () => {
    const withUnauthorizedOcie = { ...info, organizationalClothingAndIndividualEquipment: false };
    render(<AllowancesList info={withUnauthorizedOcie} />);
    expect(screen.getByTestId('ocie').textContent).toEqual('Unauthorized');
  });

  it('renders visual cues classname', () => {
    render(<AllowancesList info={info} showVisualCues />);
    expect(screen.getByText('Pro-gear').parentElement.className).toContain('rowWithVisualCue');
    expect(screen.getByText('Spouse pro-gear').parentElement.className).toContain('rowWithVisualCue');
    expect(screen.getByText('RME').parentElement.className).toContain('rowWithVisualCue');
    expect(screen.getByText('OCIE').parentElement.className).toContain('rowWithVisualCue');
  });
});
