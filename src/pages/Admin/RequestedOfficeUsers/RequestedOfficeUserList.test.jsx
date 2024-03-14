import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import { AdminContext } from 'react-admin';

import RequestedOfficeUserList from './RequestedOfficeUserList';

describe('RequestedOfficeUserList', () => {
  it('renders requested office user fields correctly', () => {
    const dataProvider = {
      getList: Promise.resolve({
        id: 1,
        name: 'Leila',
      }),
    };
    // Render the component
    const { getByTestId, getByText } = render(
      <AdminContext dataProvider={dataProvider}>
        <RequestedOfficeUserList />
      </AdminContext>,
    );

    // Verify that the requested office user fields are present
    expect(getByTestId('requested-office-user-fields')).toBeInTheDocument();

    // You can add more specific tests for other elements if needed
    expect(getByText('Transportation Office')).toBeInTheDocument();
    expect(getByText('Requested on')).toBeInTheDocument();
    // Add more assertions as needed
  });
});
