import { formatOfficeUser, formatAvailableOfficeUsers } from './queues';

const users = [
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

describe('formatOfficeUser', () => {
  it('should format a single office user', () => {
    const formattedUser = formatOfficeUser(users[0]);
    expect(formattedUser.label).toBe('Doe, John');
    expect(formattedUser.value).toBe('1234');
  });
  it('should format the office users dropdown where current user is not a supervisor', () => {
    const formattedAvailableUsers = formatAvailableOfficeUsers(users, false, '123456');
    expect(formattedAvailableUsers[0].label).toBe('—');
    expect(formattedAvailableUsers[0].value).toBe(null);
    expect(formattedAvailableUsers[1].label).toBe('User, Current');
    expect(formattedAvailableUsers[1].value).toBe('123456');
  });
  it('should format the office users dropdown where current user is a supervisor', () => {
    const formattedAvailableUsers = formatAvailableOfficeUsers(users, true, '123456');
    expect(formattedAvailableUsers[0].label).toBe('—');
    expect(formattedAvailableUsers[0].value).toBe(null);
    expect(formattedAvailableUsers[1].label).toBe('Doe, John');
    expect(formattedAvailableUsers[1].value).toBe('1234');
    expect(formattedAvailableUsers[2].label).toBe('Ipsum, Lorem');
    expect(formattedAvailableUsers[2].value).toBe('5678');
    expect(formattedAvailableUsers[3].label).toBe('User, Current');
    expect(formattedAvailableUsers[3].value).toBe('123456');
  });
});
