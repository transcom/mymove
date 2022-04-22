import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import ServicesCounselingTabNav from './ServicesCounselingTabNav';

const basicNavProps = {
  unapprovedShipmentCount: 0,
  moveCode: 'TESTCO',
};

describe('Move details tag rendering', () => {
  it('should render the move details tab container without a tag', () => {
    render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).queryByTestId('tag')).not.toBeInTheDocument();
  });

  it('should render the move details tab container with a tag that shows the count of unapproved shipments', () => {
    const moveDetailsShipmentAndAmendedOrders = {
      ...basicNavProps,
      unapprovedShipmentCount: 6,
    };
    render(<ServicesCounselingTabNav {...moveDetailsShipmentAndAmendedOrders} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('6');
  });
});

describe('Move history tab', () => {
  it('should render the move history tab container without a tag', () => {
    render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveHistory-Tab');
    expect(within(moveTaskOrderTab).queryByTestId('tag')).not.toBeInTheDocument();
  });

  describe('Tab Links', () => {
    it('should should have the correct hrefs', () => {
      render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

      const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
      expect(moveDetailsTab.getAttribute('href')).toBe(`/counseling/moves/${basicNavProps.moveCode}/details`);

      const moveHistoryTab = screen.getByTestId('MoveHistory-Tab');
      expect(moveHistoryTab.getAttribute('href')).toBe(`/counseling/moves/${basicNavProps.moveCode}/history`);
    });
  });
});
