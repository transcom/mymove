import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CustomerSupportRemarkText from './CustomerSupportRemarkText';

import { MockProviders } from 'testUtils';

const customerSupportRemark = {
  id: '672ff379-f6e3-48b4-a87d-796713f8f997',
  moveID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
  officeUserID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
  content: 'This is a remark.',
  officeUserFirstName: 'Grace',
  officeUserLastName: 'Griffin',
  createdAt: '2020-06-10T15:58:02.404031Z',
  updatedAt: '2020-06-10T15:58:02.404031Z',
};

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: 'LR4T8V' }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('CustomerSupportRemarkText', () => {
  it('can render successfully', () => {
    render(<CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />);
    expect(screen.getByText('This is a remark.')).toBeInTheDocument();
  });

  it('renders edited text if a remark has been edited', () => {
    const editedRemark = { ...customerSupportRemark, updatedAt: '2020-06-13:T13:45:03.49593Z' };

    render(<CustomerSupportRemarkText customerSupportRemark={editedRemark} />);
    expect(screen.getByText('(edited)')).toBeInTheDocument();
  });

  it('does not render edited text if a remark has been edited', () => {
    render(<CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />);
    expect(screen.queryByText('(edited)')).not.toBeInTheDocument();
  });

  it('renders the see more button when long text is present', () => {
    const longTextRemark = {
      ...customerSupportRemark,
      content:
        'This is a really long remark that will need a see more button to view it all. This is a really long remark that will need a see more button to view it all. This is a really long remark that will need a see more button to view it all. This is a really long remark that will need a see more button to view it all. This is a really long remark that will need a see more button to view it all. This is a really long remark that will need a see more button to view it all. This is a really long remark that will need a see more button to view it all.',
    };
    render(<CustomerSupportRemarkText customerSupportRemark={longTextRemark} />);
    expect(screen.getByText('(see more)')).toBeInTheDocument();
  });

  it('does not render the see more button when short text is present', () => {
    const longTextRemark = {
      ...customerSupportRemark,
      content: 'This is a short remark',
    };
    render(<CustomerSupportRemarkText customerSupportRemark={longTextRemark} />);
    expect(screen.queryByText('(see more)')).not.toBeInTheDocument();
  });

  it('renders an edit and delete buttons when the on the users own remark', () => {
    render(
      <MockProviders currentUserId={customerSupportRemark.officeUserID}>
        <CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />
      </MockProviders>,
    );
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Delete')).toBeInTheDocument();
  });

  it("does not render the edit and delete buttons on other user's remarks", () => {
    render(
      <MockProviders currentUserId="BadID">
        <CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />
      </MockProviders>,
    );
    expect(screen.queryByText('Edit')).not.toBeInTheDocument();
    expect(screen.queryByText('Delete')).not.toBeInTheDocument();
  });

  it('renders edit remark buttons and textarea', async () => {
    render(
      <MockProviders currentUserId={customerSupportRemark.officeUserID}>
        <CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />
      </MockProviders>,
    );

    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Delete')).toBeInTheDocument();
    expect(screen.queryByText('Save')).not.toBeInTheDocument();
    expect(screen.queryByText('Cancel')).not.toBeInTheDocument();
    expect(screen.getByText(customerSupportRemark.content)).toBeInTheDocument();
    expect(screen.queryByTestId('edit-remark-textarea')).not.toBeInTheDocument();

    // Open editing of the remark
    await userEvent.click(screen.getByText('Edit'));

    expect(screen.getByText('Save')).toBeInTheDocument();
    expect(screen.getByText('Cancel')).toBeInTheDocument();
    expect(screen.queryByText('Edit')).not.toBeInTheDocument();
    expect(screen.queryByText('Delete')).not.toBeInTheDocument();
    expect(screen.getByText(customerSupportRemark.content)).toBeInTheDocument();
    expect(screen.getByTestId('edit-remark-textarea')).toBeInTheDocument();
  });
});
