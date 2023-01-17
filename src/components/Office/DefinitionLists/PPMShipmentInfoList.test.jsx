import React from 'react';
import { render, screen } from '@testing-library/react';

import PPMShipmentInfoList from './PPMShipmentInfoList';

import affiliation from 'content/serviceMemberAgencies';

describe('PPMShipmentInfoList', () => {
  it('renders closeout display for Marines', () => {
    render(<PPMShipmentInfoList shipment={{ agency: affiliation.MARINES }} />);
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('TVCB');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display for Navy', () => {
    render(<PPMShipmentInfoList shipment={{ agency: affiliation.MARINES }} />);
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('NAVY');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display Coast guard', () => {
    render(<PPMShipmentInfoList shipment={{ agency: affiliation.MARINES }} />);
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('UMCG');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display for Army and Air Force', () => {
    render(<PPMShipmentInfoList shipment={{ closeoutOffice: 'Test office' }} />);
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('Test office');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });
});
