import React from 'react';
import { render, screen, within } from '@testing-library/react';

import ServicesCounselorCustomerSupportRemarks from './ServicesCounselorCustomerSupportRemarks';

import { MockProviders } from 'testUtils';
import { useCustomerSupportRemarksQueries } from 'hooks/queries';

const mockMoveLocator = 'SPTRMK';

jest.mock('hooks/queries', () => ({
  useCustomerSupportRemarksQueries: jest.fn(),
}));

const remarkData = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
  },
  customerSupportRemarks: [
    {
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      moveID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      officeUserID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      content: 'This is a comment.',
      officeUserFirstName: 'Grace',
      officeUserLastName: 'Griffin',
      createdAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
};

describe('Customer support remarks page', () => {
  it('renders the page successfully', async () => {
    useCustomerSupportRemarksQueries.mockReturnValue(remarkData);
    render(
      <MockProviders initialEntries={[`moves/${mockMoveLocator}/customer-support-remarks`]}>
        <ServicesCounselorCustomerSupportRemarks />
      </MockProviders>,
    );

    // Verify headings
    expect(screen.getByRole('heading', { name: 'Customer support remarks', level: 1 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Remarks', level: 2 })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Past remarks', level: 3 })).toBeInTheDocument();

    // Should have "create" remark" form
    const form = await screen.findByTestId('form');
    expect(form).toBeInTheDocument();

    const createLabel = 'Use this form to document any customer support provided for this move.';
    expect(within(form).getByText(createLabel)).toBeInTheDocument();

    expect(await within(form).findByTestId('textarea')).toBeInTheDocument();
    expect(await within(form).findByTestId('button')).toBeInTheDocument();

    // Should show existing remarks
    const remark = await screen.getByText(/This is a comment./);
    expect(remark).toBeInTheDocument();
  });
});
