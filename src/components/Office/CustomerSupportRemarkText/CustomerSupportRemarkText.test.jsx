import React from 'react';
import { render, screen } from '@testing-library/react';

import CustomerSupportRemarkText from './CustomerSupportRemarkText';

const customerSupportRemark = {
  id: '672ff379-f6e3-48b4-a87d-796713f8f997',
  moveID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
  officeUserID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
  content: 'This is a comment.',
  officeUserFirstName: 'Grace',
  officeUserLastName: 'Griffin',
  createdAt: '2020-06-10T15:58:02.404031Z',
};

describe('CustomerSupportRemarkText', () => {
  it('can render successfully', () => {
    render(<CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />);
    expect(screen.getByText('This is a comment.')).toBeInTheDocument();
  });
});
