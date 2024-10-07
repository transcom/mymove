import React from 'react';
import { render, screen, within, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import ServicesCounselingTabNav from './ServicesCounselingTabNav';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

const basicNavProps = {
  unapprovedShipmentCount: 0,
  missingOrdersInfoCount: 0,
  moveCode: 'TESTCO',
};

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('Move details tag rendering', () => {
  it('should render the move details tab container without a tag', () => {
    render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).queryByTestId('tag')).not.toBeInTheDocument();
  });

  it('should render the move details tab container with a tag that shows the count of action items', () => {
    const moveDetailsShipmentAndAmendedOrders = {
      ...basicNavProps,
      missingOrdersInfoCount: 4,
    };
    render(<ServicesCounselingTabNav {...moveDetailsShipmentAndAmendedOrders} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('4');
  });

  it('should render the move details tab container with a tag that shows the count of an warn and error count', () => {
    const moveDetailsShipmentAndAmendedOrders = {
      ...basicNavProps,
      missingOrdersInfoCount: 4,
      shipmentWarnConcernCount: 2,
      shipmentErrorConcernCount: 1,
    };
    render(<ServicesCounselingTabNav {...moveDetailsShipmentAndAmendedOrders} />, { wrapper: MemoryRouter });

    const moveDetailsTab = screen.getByTestId('MoveDetails-Tab');
    expect(within(moveDetailsTab).getByTestId('tag')).toHaveTextContent('7');
  });
});

describe('MTO tag rendering', () => {
  it('should render the move task order tab container without a tag', () => {
    render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
    expect(within(moveTaskOrderTab).queryByTestId('tag')).not.toBeInTheDocument();
  });
});

describe('Supporting Documents tag rendering', () => {
  it('should render the Supporting Documents tab container without a tag IF the feature flag is turned on', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    await waitFor(() => {
      const supportingDocumentsTab = screen.getByTestId('SupportingDocuments-Tab');
      expect(within(supportingDocumentsTab).queryByTestId('tag')).not.toBeInTheDocument();
    });
  });

  it('should not render the Supporting Documents tab if the feature flag is turned off', async () => {
    render(<ServicesCounselingTabNav {...basicNavProps} />, { wrapper: MemoryRouter });

    await waitFor(() => {
      expect(screen.queryByTestId('SupportingDocuments-Tab')).not.toBeInTheDocument();
    });
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

      const moveTaskOrderTab = screen.getByTestId('MoveTaskOrder-Tab');
      expect(moveTaskOrderTab.getAttribute('href')).toBe(`/counseling/moves/${basicNavProps.moveCode}/mto`);

      const moveHistoryTab = screen.getByTestId('MoveHistory-Tab');
      expect(moveHistoryTab.getAttribute('href')).toBe(`/counseling/moves/${basicNavProps.moveCode}/history`);
    });
  });
});
