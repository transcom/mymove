import React from 'react';
import { render } from '@testing-library/react';

import AddressUpdatePreview from './AddressUpdatePreview';

describe('AddressUpdatePreview', () => {
  it('does a thing', () => {
    render(
      <AddressUpdatePreview
        deliveryAddressUpdate={{
          originalAddress: { city: '', state: '', postalCode: '', streetAddress1: '', streetAddress2: '' },
          newAddress: { city: '', state: '', postalCode: '', streetAddress1: '', streetAddress2: '' },
          contractorRemarks: '',
        }}
      />,
    );

    expect(true).toBe(true);
  });
});
