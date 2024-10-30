import { formatOfficeUser, formatAvailableOfficeUsers, formatAvailableOfficeUsersForRow } from './queues';

import { MOVE_STATUSES } from 'shared/constants';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';

const availableOfficeUsers = [
  {
    firstName: 'John',
    lastName: 'Doe',
    officeUserId: '1234',
  },
  {
    firstName: 'Lorem',
    lastName: 'Ipsum',
    officeUserId: '5678',
  },
  {
    firstName: 'Current',
    lastName: 'User',
    officeUserId: '123456',
  },
];
const needsCounselingMoves = {
  isLoading: false,
  isError: false,
  queueResult: {
    totalCount: 3,
    data: [
      {
        id: 'move1',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.ARMY,
          first_name: 'test first',
          last_name: 'test last',
          dodID: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
        assignedTo: {
          officeUserId: '1234',
          firstname: 'John',
          lastname: 'Doe',
        },
      },
      {
        id: 'move2',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.COAST_GUARD,
          first_name: 'test another first',
          last_name: 'test another last',
          dodID: '4444444444',
          emplid: '4521567',
        },
        locator: 'T12AR',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
        counselingOffice: '',
        assignedTo: {
          officeUserId: '5678',
          firstname: 'Lorem',
          lastname: 'Ipsum',
        },
      },
      {
        id: 'move3',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.MARINES,
          first_name: 'test third first',
          last_name: 'test third last',
          dodID: '4444444444',
        },
        locator: 'T12MP',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Denver, 80136',
        },
        originGBLOC: 'LKNQ',
        assignedTo: {
          officeUserId: 'exampleId1',
          firstname: 'Jimmy',
          lastname: 'John',
        },
      },
    ],
  },
};

describe('formatOfficeUser', () => {
  it('should format a single office user', () => {
    const formattedUser = formatOfficeUser(availableOfficeUsers[0]);
    expect(formattedUser.label).toBe('Doe, John');
    expect(formattedUser.value).toBe('1234');
  });
  it('should format the office users dropdown where current user is not a supervisor', () => {
    const formattedAvailableUsers = formatAvailableOfficeUsers(availableOfficeUsers, false, '123456');
    expect(formattedAvailableUsers[0].label).toBe('—');
    expect(formattedAvailableUsers[0].value).toBe(null);
    expect(formattedAvailableUsers[1].label).toBe('User, Current');
    expect(formattedAvailableUsers[1].value).toBe('123456');
  });
  it('should format the office users dropdown where current user is a supervisor', () => {
    const formattedAvailableUsers = formatAvailableOfficeUsers(availableOfficeUsers, true, '123456');
    expect(formattedAvailableUsers[0].label).toBe('—');
    expect(formattedAvailableUsers[0].value).toBe(null);
    expect(formattedAvailableUsers[1].label).toBe('Doe, John');
    expect(formattedAvailableUsers[1].value).toBe('1234');
    expect(formattedAvailableUsers[2].label).toBe('Ipsum, Lorem');
    expect(formattedAvailableUsers[2].value).toBe('5678');
    expect(formattedAvailableUsers[3].label).toBe('User, Current');
    expect(formattedAvailableUsers[3].value).toBe('123456');
  });
  it('should format the office users row where the assigned user of the row is not already part of availableOfficeUsers', () => {
    // get the row data
    // row assigned to should contain a different value than available
    const needsCounselingRowWithNewAssignedUser = { ...needsCounselingMoves.queueResult.data[0] };
    needsCounselingRowWithNewAssignedUser.availableOfficeUsers = availableOfficeUsers;
    needsCounselingRowWithNewAssignedUser.assignedTo = {
      firstName: 'Adam',
      lastName: 'Sandler',
      officeUserId: '9101112',
    };

    const { formattedAvailableOfficeUsers, assignedToUser } = formatAvailableOfficeUsersForRow(
      needsCounselingRowWithNewAssignedUser,
      true,
      'CurrentUser123',
    );
    expect(formattedAvailableOfficeUsers[4].props.value).toBe('9101112');
    expect(formattedAvailableOfficeUsers[4].props.children).toBe('Sandler, Adam');
    expect(assignedToUser.value).toBe('9101112');
    expect(assignedToUser.label).toBe('Sandler, Adam');
  });
  it('should format the office users row where the assigned user of the row is already part of availableOfficeUsers', () => {
    // get the row data
    // row assigned to should contain a different value than available
    const needsCounselingRowWithNewAssignedUser = { ...needsCounselingMoves.queueResult.data[0] };
    needsCounselingRowWithNewAssignedUser.availableOfficeUsers = availableOfficeUsers;
    needsCounselingRowWithNewAssignedUser.assignedTo = {
      firstName: 'Current',
      lastName: 'User',
      officeUserId: '123456',
    };

    const { formattedAvailableOfficeUsers, assignedToUser } = formatAvailableOfficeUsersForRow(
      needsCounselingRowWithNewAssignedUser,
      true,
      'CurrentUser123',
    );
    expect(formattedAvailableOfficeUsers[3].props.value).toBe('123456');
    expect(formattedAvailableOfficeUsers[3].props.children).toBe('User, Current');
    expect(assignedToUser.value).toBe('123456');
    expect(assignedToUser.label).toBe('User, Current');
  });
});
