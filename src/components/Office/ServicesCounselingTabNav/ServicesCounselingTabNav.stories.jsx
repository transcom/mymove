import React from 'react';
import { MemoryRouter } from 'react-router';

import ServicesCounselingTabNav from './ServicesCounselingTabNav';

export default {
  title: 'components/Services Counseling Tab Navigation',
  component: ServicesCounselingTabNav,
  decorators: [
    (Story) => {
      return (
        <MemoryRouter initialEntries={['/']}>
          <div style={{ padding: '20px' }}>
            <Story />
          </div>
        </MemoryRouter>
      );
    },
  ],
};

const basicNavProps = {
  moveCode: 'TESTCO',
};

const moveWithUnapprovedShipments = {
  ...basicNavProps,
  unapprovedShipmentCount: 6,
};

export const Default = () => <ServicesCounselingTabNav {...basicNavProps} />;

export const withMoveDetailsTag = () => <ServicesCounselingTabNav {...moveWithUnapprovedShipments} />;
