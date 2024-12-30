import React from 'react';
import { AdminContext, Resource } from 'react-admin';
import { render } from '@testing-library/react';

import OfficeUserList from './OfficeUserList';
import OfficeUserShow from './OfficeUserShow';
import OfficeUserCreate from './OfficeUserCreate';
import OfficeUserEdit from './OfficeUserEdit';

const officeUser = {
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
  rejectionReason: null,
  roles: [
    {
      createdAt: '2024-10-25T19:55:13.204Z',
      id: '2458e82f-b1ab-4eca-84a6-f39666b778fd',
      roleName: 'Task Ordering Officer',
      roleType: 'task_ordering_officer',
      updatedAt: '2024-10-25T19:55:13.204Z',
    },
  ],
  status: 'APPROVED',
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
};

const dataProvider = {
  getOne: () => Promise.resolve({ data: officeUser }),
  getList: () => Promise.resolve({ data: [officeUser], total: 1 }),
};

describe('OfficeUserList page', () => {
  it('renders without crashing', async () => {
    render(
      <AdminContext dataProvider={dataProvider} basename="/system">
        <Resource name="office-users" options={{ label: 'Office Users' }} list={OfficeUserList} />
      </AdminContext>,
    );
  });
});

describe('OfficeUserShow page', () => {
  it('renders without crashing', async () => {
    render(
      <AdminContext dataProvider={dataProvider} basename="/system">
        <Resource name="office-users" options={{ label: 'Office Users' }} show={OfficeUserShow} />
      </AdminContext>,
    );
  });
});

describe('OfficeUserCreate page', () => {
  it('renders without crashing', async () => {
    render(
      <AdminContext dataProvider={dataProvider} basename="/system">
        <Resource name="office-users" options={{ label: 'Office Users' }} create={OfficeUserCreate} />
      </AdminContext>,
    );
  });
});

describe('OfficeUserEdit page', () => {
  it('renders without crashing', async () => {
    render(
      <AdminContext dataProvider={dataProvider} basename="/system">
        <Resource name="office-users" options={{ label: 'Office Users' }} edit={OfficeUserEdit} />
      </AdminContext>,
    );
  });
});
