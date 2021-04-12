/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import ConnectedProfile from './Profile';

import { MockProviders } from 'testUtils';

describe('Profile component', () => {
  const testProps = {};

  it('renders the Profile Page', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        orders: {
          test: {
            new_duty_station: {
              name: 'Test Duty Station',
            },
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            residential_address: {
              city: 'San Diego',
              state: 'CA',
              postal_code: '92131',
              street_address_1: 'Some Street',
              country: 'USA',
            },
            backup_mailing_address: {
              city: 'San Diego',
              state: 'CA',
              postal_code: '92131',
              street_address_1: 'Some Backup Street',
              country: 'USA',
            },
            current_station: {
              origin_duty_station: {
                name: 'Current Station',
              },
              grade: 'E-5',
            },
            backup_contacts: [
              {
                name: 'Backup Contact',
                telephone: '555-555-5555',
                email: 'backup@test.com',
              },
            ],
            orders: ['test'],
          },
        },
      },
    };
    render(
      <MockProviders initialState={mockState}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );

    const mainHeader = await screen.findByRole('heading', { name: 'Profile', level: 1 });

    expect(mainHeader).toBeInTheDocument();

    const contactInfoHeader = await screen.findByRole('heading', { name: 'Contact info', level: 2 });

    expect(contactInfoHeader).toBeInTheDocument();

    const serviceInfoHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(serviceInfoHeader).toBeInTheDocument();
  });

  it('renders the Profile Page when there are no orders', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        orders: {
          test: {},
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            residential_address: {
              city: 'San Diego',
              state: 'CA',
              postal_code: '92131',
              street_address_1: 'Some Street',
              country: 'USA',
            },
            backup_mailing_address: {
              city: 'San Diego',
              state: 'CA',
              postal_code: '92131',
              street_address_1: 'Some Backup Street',
              country: 'USA',
            },
            current_station: {
              origin_duty_station: {
                name: 'Current Station',
              },
              grade: 'E-5',
            },
            backup_contacts: [
              {
                name: 'Backup Contact',
                telephone: '555-555-5555',
                email: 'backup@test.com',
              },
            ],
            orders: ['test'],
          },
        },
      },
    };
    render(
      <MockProviders initialState={mockState}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );

    const mainHeader = await screen.findByRole('heading', { name: 'Profile', level: 1 });

    expect(mainHeader).toBeInTheDocument();

    const contactInfoHeader = await screen.findByRole('heading', { name: 'Contact info', level: 2 });

    expect(contactInfoHeader).toBeInTheDocument();

    const serviceInfoHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(serviceInfoHeader).toBeInTheDocument();
  });
});
