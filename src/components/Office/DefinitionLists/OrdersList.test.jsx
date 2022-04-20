import React from 'react';
import { render, screen } from '@testing-library/react';

import OrdersList from './OrdersList';

const ordersInfo = {
  currentDutyLocation: { name: 'JBSA Lackland' },
  newDutyLocation: { name: 'JB Lewis-McChord' },
  issuedDate: '2020-03-08',
  reportByDate: '2020-04-01',
  departmentIndicator: 'NAVY_AND_MARINES',
  ordersNumber: '999999999',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'HHG_PERMITTED',
  tacMDC: '9999',
  sacSDN: '999 999999 999',
};

// what ordersInfo from above should be rendered as
const expectedRenderedOrdersInfo = {
  currentDutyLocation: 'JBSA Lackland',
  newDutyLocation: 'JB Lewis-McChord',
  issuedDate: '08 Mar 2020',
  reportByDate: '01 Apr 2020',
  departmentIndicator: '17 Navy and Marine Corps',
  ordersNumber: '999999999',
  ordersType: 'Permanent Change Of Station (PCS)',
  ordersTypeDetail: 'Shipment of HHG Permitted',
  tacMDC: '9999',
  sacSDN: '999 999999 999',
};

const ordersInfoMissing = {
  currentDutyLocation: { name: 'JBSA Lackland' },
  newDutyLocation: { name: 'JB Lewis-McChord' },
  issuedDate: '2020-03-08',
  reportByDate: '2020-04-01',
  departmentIndicator: '',
  ordersNumber: '',
  ordersType: '',
  ordersTypeDetail: '',
  tacMDC: '',
  sacSDN: '999 999999 999',
};

describe('OrdersList', () => {
  it('renders formatted orders info', () => {
    render(<OrdersList ordersInfo={ordersInfo} />);
    Object.keys(expectedRenderedOrdersInfo).forEach((key) => {
      expect(screen.getByText(expectedRenderedOrdersInfo[key])).toBeInTheDocument();
    });
  });

  it('renders missing orders info as warning if showMissingWarnings is included', () => {
    render(<OrdersList ordersInfo={ordersInfoMissing} />);
    expect(screen.getByTestId('departmentIndicator').textContent).toEqual('Missing');
    expect(screen.getByTestId('ordersNumber').textContent).toEqual('Missing');
    expect(screen.getByTestId('ordersType').textContent).toEqual('Missing');
    expect(screen.getByTestId('ordersTypeDetail').textContent).toEqual('Missing');
    expect(screen.getByTestId('tacMDC').textContent).toEqual('Missing');
  });

  it('renders missing orders info as dashes if showMissingWarnings is false', () => {
    render(<OrdersList ordersInfo={ordersInfoMissing} showMissingWarnings={false} />);
    expect(screen.getByTestId('departmentIndicator').textContent).toEqual('—');
    expect(screen.getByTestId('ordersNumber').textContent).toEqual('—');
    expect(screen.getByTestId('ordersType').textContent).toEqual('—');
    expect(screen.getByTestId('ordersTypeDetail').textContent).toEqual('—');
    expect(screen.getByTestId('tacMDC').textContent).toEqual('—');
  });
});
