import React from 'react';
import { render, screen } from '@testing-library/react';

import OrdersList from './OrdersList';

const ordersInfo = {
  currentDutyStation: { name: 'JBSA Lackland' },
  newDutyStation: { name: 'JB Lewis-McChord' },
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
  currentDutyStation: 'JBSA Lackland',
  newDutyStation: 'JB Lewis-McChord',
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
  currentDutyStation: { name: 'JBSA Lackland' },
  newDutyStation: { name: 'JB Lewis-McChord' },
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
    render(<OrdersList ordersInfo={ordersInfo} showLOAWarnings />);
    Object.keys(expectedRenderedOrdersInfo).forEach((key) => {
      expect(screen.getByText(expectedRenderedOrdersInfo[key])).toBeInTheDocument();
    });
  });

  it('renders missing orders info', () => {
    render(<OrdersList ordersInfo={ordersInfoMissing} showLOAWarnings />);
    expect(screen.getByTestId('departmentIndicator').textContent).toEqual('Missing');
    expect(screen.getByTestId('ordersNumber').textContent).toEqual('Missing');
    expect(screen.getByTestId('ordersType').textContent).toEqual('Missing');
    expect(screen.getByTestId('ordersTypeDetail').textContent).toEqual('Missing');
  });

  it('renders HHG TAC as missing if showLOAWarnings is included', () => {
    render(<OrdersList ordersInfo={ordersInfoMissing} showLOAWarnings />);
    expect(screen.getByTestId('tacMDC').textContent).toEqual('Missing');
  });

  it('renders HHG TAC as options if showLOAWarnings is false', () => {
    render(<OrdersList ordersInfo={ordersInfoMissing} showLOAWarnings={false} />);
    expect(screen.getByTestId('tacMDC').textContent).toEqual('—');
  });
});
