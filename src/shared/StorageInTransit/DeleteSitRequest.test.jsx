import React from 'react';
import { shallow } from 'enzyme';

import DeleteSitRequest from './DeleteSitRequest';

let wrapper;
const requestedStorageInTransit = {
  status: 'REQUESTED',
  location: 'DESTINATION',
  estimated_start_date: '2019-05-15',
};
const onCloseHandler = jest.fn();
const onDeleteHandler = jest.fn();

describe('Delete SIT request', () => {
  beforeEach(() => {
    onCloseHandler.mockClear();
    wrapper = shallow(
      <DeleteSitRequest
        storageInTransit={requestedStorageInTransit}
        onClose={onCloseHandler}
        onDelete={onDeleteHandler}
      />,
    );
  });

  it('renders without crashing', () => {
    expect(wrapper.find('.sit-delete-warning').length).toEqual(1);
  });

  it('clicking calls close handler', () => {
    wrapper.find('.sit-delete-cancel').simulate('click');
    expect(onCloseHandler.mock.calls.length).toBe(1);
  });

  it('clicking calls delete handler', () => {
    wrapper.find('.sit-delete-delete').simulate('click');
    expect(onDeleteHandler.mock.calls.length).toBe(1);
  });
});
