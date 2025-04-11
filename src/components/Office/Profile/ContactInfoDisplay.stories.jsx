import React from 'react';
import { MemoryRouter } from 'react-router-dom';

import ContactInfoDisplay from './ContactInfoDisplay';

export default {
  title: 'Office Components/ContactInfoDisplay',
  component: ContactInfoDisplay,
  decorators: [
    (Story) => (
      <MemoryRouter>
        <div style={{ padding: '20px', maxWidth: '400px' }}>
          <Story />
        </div>
      </MemoryRouter>
    ),
  ],
};

const mockUserInfo = {
  name: 'John Doe',
  telephone: '123-456-7890',
  email: 'john.doe@example.com',
};

const mockEditURL = '/edit-profile';

export const Basic = () => <ContactInfoDisplay officeUserInfo={mockUserInfo} editURL={mockEditURL} />;
