import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom-old';

import TXOTabNav from './TXOTabNav';

const basicNavProps = {
  order: {},
  moveCode: 'TESTCO',
};

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
