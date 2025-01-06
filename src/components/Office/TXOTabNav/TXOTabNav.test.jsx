import React from 'react';
import { render, screen, within, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import TXOTabNav from './TXOTabNav';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

const basicNavProps = {
  order: {},
  moveCode: 'TESTCO',
};

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('Move details tag rendering', () => {
  it('should render the move details tab container without a tag', () => {
    render(<TXOTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).queryByTestId('tag')).toBeFalsy();
  });

  it('should render the move details tab container with a tag that shows the count of items that need attention when the orders have been amended', () => {
    const moveDetailsAmendedOrders = {
      ...basicNavProps,
      order: {
        uploadedAmendedOrderID: '1234',
      },
    };
    render(<TXOTabNav {...moveDetailsAmendedOrders} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('1');
  });

  it('should render the move details tab container with a tag that shows the count of items that need attention when there are unapproved shipments', () => {
    const moveDetailsOneShipment = {
      ...basicNavProps,
      unapprovedShipmentCount: 1,
      missingOrdersInfoCount: 4,
    };
    render(<TXOTabNav {...moveDetailsOneShipment} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('5');
  });

  it('should render the move details tab container with a tag that shows the count of items that need attention when there are approved shipments with a delivery address update requiring TXO review', () => {
    const moveDetailsOneShipment = {
      ...basicNavProps,
      shipmentsWithDeliveryAddressUpdateRequestedCount: 1,
    };
    render(<TXOTabNav {...moveDetailsOneShipment} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('1');
  });

  it('should render the move details tab container with a tag that shows the count of items that need attention when the orders have been amended and there are unapproved shipments', () => {
    const moveDetailsShipmentAndAmendedOrders = {
      ...basicNavProps,
      order: {
        uploadedAmendedOrderID: '1234',
      },
      unapprovedShipmentCount: 1,
    };
    render(<TXOTabNav {...moveDetailsShipmentAndAmendedOrders} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('2');
  });
});

describe('Move task order tag rendering', () => {
  it('should render the move task order tab container without a tag', () => {
    render(<TXOTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
    expect(within(moveTaskOrderTab).queryByTestId('tag')).toBeFalsy();
  });

  it('should render the move task order tab container with a tag that shows the count of shipments that have excess weight risk', () => {
    const moveTaskOrderWithExcessRisk = {
      ...basicNavProps,
      excessWeightRiskCount: 1,
    };
    render(<TXOTabNav {...moveTaskOrderWithExcessRisk} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
    expect(within(moveTaskOrderTab).getByTestId('tag')).toHaveTextContent('1');
  });
  it('should render the move task order tab container with a tag that shows the count of items that need attention when there are unapproved shipments', () => {
    const moveTaskOrderWithUnapprovedServiceItem = {
      ...basicNavProps,
      unapprovedServiceItemCount: 1,
    };
    render(<TXOTabNav {...moveTaskOrderWithUnapprovedServiceItem} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
    expect(within(moveTaskOrderTab).getByTestId('tag')).toHaveTextContent('1');
  });
  it('should render the move task order tab container with a tag that shows the count of shipments with SIT extensions needing review', () => {
    const moveTaskOrderWithUnapprovedServiceItem = {
      ...basicNavProps,
      unapprovedSITExtensionCount: 1,
    };
    render(<TXOTabNav {...moveTaskOrderWithUnapprovedServiceItem} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
    expect(within(moveTaskOrderTab).getByTestId('tag')).toHaveTextContent('1');
  });
  it('should render the move task order tab container with a tag that shows the count of items that need attention when there are unapproved ServiceItems and an excessive weight risk', () => {
    const moveTaskOrderWithUnapprovedServiceItemAndExcessWeight = {
      ...basicNavProps,
      excessWeightRiskCount: 1,
      unapprovedServiceItemCount: 1,
    };
    render(<TXOTabNav {...moveTaskOrderWithUnapprovedServiceItemAndExcessWeight} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
    expect(within(moveTaskOrderTab).getByTestId('tag')).toHaveTextContent('2');
  });
});

describe('Supporting Documents tag rendering', () => {
  it('should render the Supporting Documents tab container without a tag IF the feature flag is turned on', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    render(<TXOTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    await waitFor(() => {
      const supportingDocumentsTab = screen.getByTestId('SupportingDocuments-Tab');
      expect(within(supportingDocumentsTab).queryByTestId('tag')).not.toBeInTheDocument();
    });
  });

  it('should not render the Supporting Documents tab if the feature flag is turned off', async () => {
    render(<TXOTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    await waitFor(() => {
      expect(screen.queryByTestId('SupportingDocuments-Tab')).not.toBeInTheDocument();
    });
  });
});
