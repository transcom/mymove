import React from 'react';
import { AdminContext, Resource } from 'react-admin';
import { render } from '@testing-library/react';

import RejectedOfficeUserList from './RejectedOfficeUserList';
import RejectedOfficeUserShow from './RejectedOfficeUserShow';

const RejectedOfficeUser = {
  active: true,
  createdAt: '2024-10-25T20:04:29.658Z',
  edipi: null,
  email: 'Siobhan@example.com',
  firstName: 'Siobhan',
  id: '49136521-a02a-43dd-8884-9b8fcca198d3',
  lastName: "O'Testoghue",
  middleInitials: null,
  otherUniqueId: null,
  privileges: null,
  rejectionReason: 'Testing rejection reason',
  roles: [
    {
      createdAt: '2024-10-25T19:55:13.204Z',
      id: '2458e82f-b1ab-4eca-84a6-f39666b778fd',
      roleName: 'Task Ordering Officer',
      roleType: 'task_ordering_officer',
      updatedAt: '2024-10-25T19:55:13.204Z',
    },
  ],
  status: 'REJECTED',
  telephone: '787-787-7878',
  transportationOfficeAssignments: [
    {
      createdAt: '2024-10-25T20:04:29.670Z',
      officeUserId: '49136521-a02a-43dd-8884-9b8fcca198d3',
      primaryOffice: true,
      transportationOfficeId: '171b54fa-4c89-45d8-8111-a2d65818ff8c',
      updatedAt: '2024-10-25T20:04:29.670Z',
    },
  ],
  transportationOfficeId: '171b54fa-4c89-45d8-8111-a2d65818ff8c',
  updatedAt: '2024-10-25T20:04:29.658Z',
  userId: '465a47b2-8ce1-46d8-86b5-bdf6639ced22',
  rejectedOn: '2025-10-25T20:04:29.658Z',
};

const dataProvider = {
  getOne: () => Promise.resolve({ data: RejectedOfficeUser }),
  getList: () => Promise.resolve({ data: [RejectedOfficeUser], total: 1 }),
};

describe('RejectedOfficeUserList page', () => {
  it('renders without crashing', async () => {
    render(
      <AdminContext dataProvider={dataProvider} basename="/system">
        <Resource name="office-users" options={{ label: 'Office Users' }} list={RejectedOfficeUserList} />
      </AdminContext>,
    );
  });
});

describe('RejectedOfficeUserShow page', () => {
  it('renders without crashing', async () => {
    render(
      <AdminContext dataProvider={dataProvider} basename="/system">
        <Resource name="office-users" options={{ label: 'Office Users' }} show={RejectedOfficeUserShow} />
      </AdminContext>,
    );
  });
});
