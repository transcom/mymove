import React from 'react';
import { render, screen } from '@testing-library/react';

import ServicesCounselingOrdersList from './ServicesCounselingOrdersList';

const ordersInfo = {
  currentDutyStation: { name: 'JBSA Lackland' },
  newDutyStation: { name: 'JB Lewis-McChord' },
  issuedDate: '2020-03-08',
  reportByDate: '2020-04-01',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
};

// what ordersInfo from above should be rendered as
const expectedRenderedOrdersInfo = {
  currentDutyStation: 'JBSA Lackland',
  newDutyStation: 'JB Lewis-McChord',
  issuedDate: '08 Mar 2020',
  reportByDate: '01 Apr 2020',
  ordersType: 'Permanent Change Of Station (PCS)',
};

const ordersInfoOnlyForTOO = {
  departmentIndicator: '17 Navy and Marine Corps',
  ordersNumber: '999999999',
  ordersTypeDetail: 'Shipment of HHG Permitted',
  tacMDC: '9999',
  sacSDN: '999 999999 999',
};

describe('ServicesCounselingOrdersList', () => {
  it('renders formatted orders info', () => {
    render(<ServicesCounselingOrdersList ordersInfo={ordersInfo} />);

    Object.keys(expectedRenderedOrdersInfo).forEach((key) => {
      expect(screen.getByText(expectedRenderedOrdersInfo[key])).toBeInTheDocument();
    });
  });

  it('does not render orders info that are only relevant to a TOO', () => {
    render(<ServicesCounselingOrdersList ordersInfo={ordersInfo} />);

    Object.keys(ordersInfoOnlyForTOO).forEach((key) => {
      expect(screen.queryByText(ordersInfoOnlyForTOO[key])).not.toBeInTheDocument();
    });
  });
});
