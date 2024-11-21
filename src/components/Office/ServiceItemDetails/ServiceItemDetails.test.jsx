import React from 'react';
import { render, screen } from '@testing-library/react';

import ServiceItemDetails from './ServiceItemDetails';

const sitStatus = {
  currentSIT: {
    sitAuthorizedEndDate: '2024-03-17',
  },
};

const shipment = {
  sitDaysAllowance: 90,
};

const details = {
  description: 'some description',
  pickupPostalCode: '90210',
  SITPostalCode: '12345',
  sitEntryDate: '2024-03-11T00:00:00.000Z',
  reason: 'some reason',
  itemDimensions: { length: 1000, width: 2500, height: 3000 },
  crateDimensions: { length: 2000, width: 3500, height: 4000 },
  customerContacts: [
    { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15', dateOfContact: '2020-09-15' },
    { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21', dateOfContact: '2020-09-21' },
  ],
  estimatedWeight: 2500,
  sitCustomerContacted: '2024-03-14T00:00:00.000Z',
  sitRequestedDelivery: '2024-03-15T00:00:00.000Z',
  sitDepartureDate: '2024-03-16T00:00:00.000Z',
  sitDeliveryMiles: 50,
  sitOriginHHGOriginalAddress: {
    city: 'Origin Original Tampa',
    eTag: 'MjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '7fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  sitOriginHHGActualAddress: {
    city: 'Origin Actual MacDill',
    eTag: 'HjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '8fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  sitDestinationOriginalAddress: {
    city: 'Destination Original Tampa',
    eTag: 'MjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '7fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  sitDestinationFinalAddress: {
    city: 'Destination Final MacDill',
    eTag: 'HjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '8fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  estimatedPrice: 2800,
};

const submittedServiceItemDetails = {
  description: 'some description',
  pickupPostalCode: '90210',
  SITPostalCode: '12345',
  sitEntryDate: '2024-03-11T00:00:00.000Z',
  reason: 'some reason',
  itemDimensions: { length: 1000, width: 2500, height: 3000 },
  crateDimensions: { length: 2000, width: 3500, height: 4000 },
  customerContacts: [
    { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15', dateOfContact: '2020-09-15' },
    { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21', dateOfContact: '2020-09-21' },
  ],
  estimatedWeight: 2500,
  sitCustomerContacted: '2024-03-14T00:00:00.000Z',
  sitRequestedDelivery: '2024-03-15T00:00:00.000Z',
  sitDepartureDate: '2024-03-16T00:00:00.000Z',
  sitDeliveryMiles: 50,
  sitOriginHHGOriginalAddress: {
    city: 'Origin Original Tampa',
    eTag: 'MjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '7fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  sitOriginHHGActualAddress: {
    city: 'Origin Actual MacDill',
    eTag: 'HjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '8fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  sitDestinationOriginalAddress: {
    city: 'Destination Original Tampa',
    eTag: 'MjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '7fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  sitDestinationFinalAddress: {
    city: 'Destination Final MacDill',
    eTag: 'HjAyNC0wMy0xMlQxOTo1OTowOC41NjkxMzla',
    id: '8fd6cb90-54cd-44d8-8735-102e28734d84',
    postalCode: '33621',
    state: 'FL',
    streetAddress1: 'MacDill',
  },
  status: 'SUBMITTED',
};

const serviceRequestDocs = [
  {
    uploads: [
      {
        filename: '/mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-receipt.pdf',
      },
    ],
  },
];

const detailsRejectedServiceItem = { ...details, rejectionReason: 'some rejection reason' };

const nilDetails = {
  estimatedWeight: null,
};

describe('ServiceItemDetails Domestic Destination SIT', () => {
  it('renders DDASIT details', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DDASIT"
        details={details}
        shipment={shipment}
        sitStatus={sitStatus}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );
    expect(screen.getByText('Original Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('Destination Original Tampa, FL 33621')).toBeInTheDocument();

    expect(screen.getByText("Add'l SIT Start Date:")).toBeInTheDocument();
    expect(screen.getByText('12 Mar 2024')).toBeInTheDocument();

    expect(screen.queryByText('Customer contacted homesafe:')).not.toBeInTheDocument();
    expect(screen.queryByText('14 Mar 2024')).not.toBeInTheDocument();

    expect(screen.getByText('# of days approved for:')).toBeInTheDocument();
    expect(screen.getByText('89 days')).toBeInTheDocument();

    expect(screen.getByText('SIT expiration date:')).toBeInTheDocument();
    expect(screen.getByText('17 Mar 2024')).toBeInTheDocument();

    expect(screen.queryByText('Customer requested delivery date:')).not.toBeInTheDocument();
    expect(screen.queryByText('15 Mar 2024')).not.toBeInTheDocument();

    expect(screen.queryByText('SIT departure date:')).not.toBeInTheDocument();
    expect(screen.queryByText('16 Mar 2024')).not.toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('renders DDDSIT details', () => {
    render(<ServiceItemDetails id="1" code="DDDSIT" details={details} serviceRequestDocs={serviceRequestDocs} />);
    expect(screen.getByText('Original Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('Destination Original Tampa, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Final Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('Destination Final MacDill, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Delivery miles out of SIT:')).toBeInTheDocument();
    expect(screen.getByText('50')).toBeInTheDocument();

    expect(screen.getByText('Customer contacted homesafe:')).toBeInTheDocument();
    expect(screen.getByText('14 Mar 2024')).toBeInTheDocument();

    expect(screen.getByText('Customer requested delivery date:')).toBeInTheDocument();
    expect(screen.getByText('15 Mar 2024')).toBeInTheDocument();

    expect(screen.getByText('SIT departure date:')).toBeInTheDocument();
    expect(screen.getByText('16 Mar 2024')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
  it('renders DDDSIT details with - for the final delivery address is service item is in submitted state', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DDDSIT"
        details={submittedServiceItemDetails}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Final Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('-')).toBeInTheDocument();
  });
  it('renders DDFSIT details', () => {
    render(<ServiceItemDetails id="1" code="DDFSIT" details={details} serviceRequestDocs={serviceRequestDocs} />);
    expect(screen.getByText('Original Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('Destination Original Tampa, FL 33621')).toBeInTheDocument();
  });
  it('renders DDSFSC details', () => {
    render(<ServiceItemDetails id="1" code="DDSFSC" details={details} serviceRequestDocs={serviceRequestDocs} />);
    expect(screen.getByText('Original Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('Destination Original Tampa, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Final Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('Destination Final MacDill, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Delivery miles out of SIT:')).toBeInTheDocument();
    expect(screen.getByText('50')).toBeInTheDocument();
  });
  it('renders DDSFSC details with - for the final delivery address is service item is in submitted state', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DDSFSC"
        details={submittedServiceItemDetails}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Final Delivery Address:')).toBeInTheDocument();
    expect(screen.getByText('-')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails Domestic Origin SIT', () => {
  it(`renders DOASIT details`, () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DOASIT"
        details={details}
        shipment={shipment}
        sitStatus={sitStatus}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Original Pickup Address:')).toBeInTheDocument();
    expect(screen.getByText('Origin Original Tampa, FL 33621')).toBeInTheDocument();

    expect(screen.getByText("Add'l SIT Start Date:")).toBeInTheDocument();
    expect(screen.getByText('12 Mar 2024')).toBeInTheDocument();

    expect(screen.getByText('# of days approved for:')).toBeInTheDocument();
    expect(screen.getByText('89 days')).toBeInTheDocument();

    expect(screen.getByText('SIT expiration date:')).toBeInTheDocument();
    expect(screen.getByText('17 Mar 2024')).toBeInTheDocument();

    expect(screen.getByText('Customer contacted homesafe:')).toBeInTheDocument();
    expect(screen.getByText('14 Mar 2024')).toBeInTheDocument();

    expect(screen.getByText('Customer requested delivery date:')).toBeInTheDocument();
    expect(screen.getByText('15 Mar 2024')).toBeInTheDocument();

    expect(screen.getByText('SIT departure date:')).toBeInTheDocument();
    expect(screen.getByText('16 Mar 2024')).toBeInTheDocument();
  });

  it(`renders DOPSIT details`, () => {
    render(<ServiceItemDetails id="1" code="DOPSIT" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('Original Pickup Address:')).toBeInTheDocument();
    expect(screen.getByText('Origin Original Tampa, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Actual Pickup Address:')).toBeInTheDocument();
    expect(screen.getByText('Origin Actual MacDill, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Delivery miles into SIT:')).toBeInTheDocument();
    expect(screen.getByText('50')).toBeInTheDocument();
  });

  it(`renders DOSFSC details`, () => {
    render(<ServiceItemDetails id="1" code="DOSFSC" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('Original Pickup Address:')).toBeInTheDocument();
    expect(screen.getByText('Origin Original Tampa, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Actual Pickup Address:')).toBeInTheDocument();
    expect(screen.getByText('Origin Actual MacDill, FL 33621')).toBeInTheDocument();

    expect(screen.getByText('Delivery miles into SIT:')).toBeInTheDocument();
    expect(screen.getByText('50')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails for DOFSIT', () => {
  it('renders SIT entry date, ZIP, original pickup address, and reason', () => {
    render(<ServiceItemDetails id="1" code="DOFSIT" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('Original Pickup Address:')).toBeInTheDocument();
    expect(screen.getByText('Origin Original Tampa, FL 33621')).toBeInTheDocument();
    expect(screen.getByText('SIT entry date:')).toBeInTheDocument();
    expect(screen.getByText('11 Mar 2024')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Domestic Destination SIT', () => {
  it('renders first and second customer contact and available delivery date', () => {
    render(<ServiceItemDetails id="1" code="DDFSIT" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('Customer contact attempt 1:')).toBeInTheDocument();
    expect(screen.getByText('15 Sep 2020, 1200Z')).toBeInTheDocument();
    expect(screen.getByText('15 Sep 2020')).toBeInTheDocument();
    expect(screen.getByText('Customer contact attempt 2:')).toBeInTheDocument();
    expect(screen.getByText('21 Sep 2020, 2300Z')).toBeInTheDocument();
    expect(screen.getByText('21 Sep 2020')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Crating', () => {
  it('renders description and dimensions', () => {
    render(<ServiceItemDetails id="1" code="DCRT" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Domestic Shuttling', () => {
  it.each([['DOSHUT'], ['DDSHUT']])('renders formatted estimated weight and reason', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('2,500 lbs')).toBeInTheDocument();
    expect(screen.getByText('estimated weight')).toBeInTheDocument();
    expect(screen.getByText('Reason:')).toBeInTheDocument();
    expect(screen.getByText('some reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it.each([['DOSHUT'], ['DDSHUT']])('renders estimated weight nil values with an em dash', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={nilDetails} />);

    expect(screen.getByText('— lbs')).toBeInTheDocument();
    expect(screen.getByText('estimated weight')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails Crating Rejected', () => {
  it('renders the rejection reason field when it is populated with information', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DUCRT"
        details={detailsRejectedServiceItem}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Description:')).toBeInTheDocument();
    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('some rejection reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Crating Rejected', () => {
  it('renders the rejection reason field when it is populated with information', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DCRT"
        details={detailsRejectedServiceItem}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Description:')).toBeInTheDocument();
    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('some rejection reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Estimated Price for DLH, DSH, FSC, DOP, DDP, DPK, DUPK', () => {
  it.each([['DLH'], ['DSH'], ['FSC'], ['DOP'], ['DDP'], ['DPK'], ['DUPK']])(
    'renders the formatted estimated price field for the service items',
    (code) => {
      render(
        <ServiceItemDetails
          id="1"
          code={code}
          details={details}
          shipment={shipment}
          serviceRequestDocs={serviceRequestDocs}
        />,
      );

      expect(screen.getByText('Estimated Price:')).toBeInTheDocument();
      expect(screen.getByText('$28.00')).toBeInTheDocument();
    },
  );

  const noEstimatePriceDetails = {};

  it.each([['DLH'], ['DSH'], ['FSC'], ['DOP'], ['DDP'], ['DPK'], ['DUPK']])(
    'renders - for estimated price when price is not in details',
    (code) => {
      render(
        <ServiceItemDetails
          id="1"
          code={code}
          details={noEstimatePriceDetails}
          shipment={shipment}
          serviceRequestDocs={serviceRequestDocs}
        />,
      );

      expect(screen.getByText('Estimated Price:')).toBeInTheDocument();
      expect(screen.getByText('-')).toBeInTheDocument();
    },
  );
});

describe('ServiceItemDetails Price for MS, CS', () => {
  it.each([['MS'], ['CS']])('renders the formatted price field for the service items', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={details} />);

    expect(screen.getByText('Price:')).toBeInTheDocument();
    expect(screen.getByText('$28.00')).toBeInTheDocument();
  });
});
