import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { object, text } from '@storybook/addon-knobs';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

const info = {
  primeActualWeight: 2000,
  storageFacility: {
    address: {
      city: 'Anytown',
      country: 'USA',
      postalCode: '90210',
      state: 'OK',
      streetAddress1: '555 Main Ave',
      streetAddress2: 'Apartment 900',
    },
    facilityName: 'my storage',
    lotNumber: '2222',
  },
  serviceOrderNumber: '12341234',
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  agents: [
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
  tacType: 'HHG',
  sacType: 'NTS',
  tac: '1234',
  sac: '1234123412',
};

const labels = {
  primeActualWeight: 'Shipment weight',
  storageFacility: 'Storage facility info',
  serviceOrderNumber: 'Service order #',
  storageFacilityAddress: 'Storage facility address',
  requestedDeliveryDate: 'Preferred delivery date',
  destinationAddress: 'Delivery address',
  secondaryDeliveryAddress: 'Second delivery address',
  agents: 'Receiving agent',
  customerRemarks: 'Customer remarks',
  counselorRemarks: 'Counselor remarks',
  tacType: 'TAC',
  sacType: 'SAC',
};

describe('Shipment Info List', () => {
  it('renders all fields when provided and expanded', () => {
    render(<NTSRShipmentInfoList isExpanded shipment={info} />);

    const primeActualWeight = screen.getByText(labels.primeActualWeight);
    expect(within(primeActualWeight.parentElement).getByText('2000 lbs')).toBeInTheDocument();

    const storageFacility = screen.getByText(labels.storageFacility);
    expect(within(storageFacility.parentElement).getByText(info.storageFacility.facilityName)).toBeInTheDocument();

    const serviceOrderNumber = screen.getByText(labels.serviceOrderNumber);
    expect(within(serviceOrderNumber.parentElement).getByText(info.serviceOrderNumber)).toBeInTheDocument();

    const storageFacilityAddress = screen.getByText(labels.storageFacilityAddress);
    expect(
      within(storageFacilityAddress.parentElement).getByText(info.storageFacility.address.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const destinationAddress = screen.getByText(labels.destinationAddress);
    expect(
      within(destinationAddress.parentElement).getByText(info.destinationAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const secondaryDeliveryAddress = screen.getByText(labels.secondaryDeliveryAddress);
    expect(
      within(secondaryDeliveryAddress.parentElement).getByText(info.secondaryDeliveryAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const receivingAgent = screen.getByText(labels.agents);
    expect(within(receivingAgent.parentElement).getByText(info.agents[0].email, { exact: false })).toBeInTheDocument();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(info.counselorRemarks)).toBeInTheDocument();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(info.customerRemarks)).toBeInTheDocument();

    const tacType = screen.getByText(labels.tacType);
    expect(within(tacType.parentElement).getByText('1234 (HHG)')).toBeInTheDocument();

    const sacType = screen.getByText(labels.sacType);
    expect(within(sacType.parentElement).getByText('1234123412 (NTS)')).toBeInTheDocument();
  });

  // it('renders a dash for for non-required missing items', () => {
  //   render(
  //     <NTSRShipmentInfoList
  //       isExpanded
  //       shipment={{
  //         requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
  //         storageFacility: object('storageFacility', info.storageFacility),
  //         destinationAddress: object('destinationAddress', info.destinationAddress),
  //       }}
  //     />,
  //   );
  //
  //   const counselorRemarks = screen.getByText(labels.counselorRemarks);
  //   console.log(counselorRemarks)
  //   expect(within(counselorRemarks.parentElement).getByText('—')).toBeInTheDocument();
  //   expect(counselorRemarks.parentElement.toHaveClass('OfficeDefinitionLists_warning__Z3q_v'));
  //
  //   const tacType = screen.getByText(labels.tacType);
  //   expect(within(tacType.parentElement).getByText('—')).toBeInTheDocument();
  //   expect(tacType.parentElement.classList.contains('OfficeDefinitionLists_warning__Z3q_v'));
  //
  //   const sacType = screen.getByText(labels.sacType);
  //   expect(within(sacType.parentElement).getByText('—')).toBeInTheDocument();
  //   expect(sacType.parentElement.classList.contains('OfficeDefinitionLists_warning__Z3q_v'));
  //
  //   const primeActualWeight = screen.getByText(labels.primeActualWeight);
  //   expect(within(primeActualWeight.parentElement).getByText('—')).toBeInTheDocument();
  //   expect(primeActualWeight.parentElement.classList.contains('OfficeDefinitionLists_warning__Z3q_v'));
  //
  //   const serviceOrderNumber = screen.getByText(labels.serviceOrderNumber);
  //   expect(within(serviceOrderNumber.parentElement).getByText('—')).toBeInTheDocument();
  //   expect(serviceOrderNumber.parentElement.classList.contains('OfficeDefinitionLists_warning__Z3q_v'));
  // });
  //
  // it('shows Missing for required missing items', () => {
  //   render(
  //     <NTSRShipmentInfoList
  //       shipment={{
  //         counselorRemarks: text('counselorRemarks', info.counselorRemarks),
  //         requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
  //         storageFacility: object('storageFacility', info.storageFacility),
  //         destinationAddress: object('destinationAddress', info.destinationAddress),
  //       }}
  //     />,
  //   );
  //
  //   const storageFacility = screen.getByText(labels.storageFacility);
  //   expect(within(storageFacility.parentElement).getByText('Missing').toBeInTheDocument());
  //   expect(within(storageFacility.parentElement).classList.contains('OfficeDefinitionLists_missingInfoError__3ckQO'));
  //
  //   const storageFacilityAddress = screen.getByText(labels.storageFacilityAddress);
  //   expect(within(storageFacilityAddress.parentElement).getByText('Missing')).toBeInTheDocument();
  //   expect(within(storageFacilityAddress.parentElement).classList.contains('OfficeDefinitionLists_missingInfoError__3ckQO'));
  // });

  it('hides fields when collapsed', () => {
    render(<NTSRShipmentInfoList isExpanded={false} shipment={info} />);

    expect(screen.queryByText(labels.primeActualWeight)).toBeNull();
    expect(screen.queryByText(labels.storageFacility)).toBeNull();
    expect(screen.queryByText(labels.serviceOrderNumber)).toBeNull();
    expect(screen.queryByText(labels.secondaryDeliveryAddress)).toBeNull();
    expect(screen.queryByText(labels.agents)).toBeNull();
    expect(screen.queryByText(labels.customerRemarks)).toBeNull();
  });
});
