import React from 'react';
import { shallow } from 'enzyme';

import { OfficeEditor } from './OfficeEditor';

const storageInTransit = {
  estimated_start_date: '2019-02-12',
  id: '5cd370a1-ac3d-4fb3-86a3-c4f23e289687',
  location: 'ORIGIN',
  shipment_id: 'dd67cec5-334a-4209-a9d9-a14485414052',
  status: 'REQUESTED',
  warehouse_address: {
    city: 'Beverly Hills',
    postal_code: '90210',
    state: 'CA',
    street_address_1: '123 Any Street',
  },
  warehouse_id: '76567867',
  warehouse_name: 'haus',
};

let wrapper;
const cancel = jest.fn();
const approveStorageInTransit = jest.fn();
const denyStorageInTransit = jest.fn();
const onClose = jest.fn();
const saveAndClose = jest.fn();

describe('given an editor', () => {
  describe('when the form is disabled', () => {
    beforeEach(() => {
      cancel.mockClear();
      wrapper = shallow(
        <OfficeEditor
          approveStorageInTransit={approveStorageInTransit}
          denyStorageInTransit={denyStorageInTransit}
          onClose={onClose}
          storageInTransit={storageInTransit}
          submitForm={saveAndClose}
          formEnabled={false}
          hasSubmitSucceeded={false}
        />,
      );
    });

    it('renders without crashing', () => {
      expect(wrapper.exists('.storage-in-transit-panel-modal')).toBe(true);
    });

    it('buttons are disabled', () => {
      expect(wrapper.find('button.usa-button-primary').prop('disabled')).toBeTruthy();
    });
  });

  describe('when the form is enabled', () => {
    beforeEach(() => {
      cancel.mockClear();
      wrapper = shallow(
        <OfficeEditor
          approveStorageInTransit={approveStorageInTransit}
          denyStorageInTransit={denyStorageInTransit}
          onClose={onClose}
          storageInTransit={storageInTransit}
          submitForm={saveAndClose}
          formEnabled={true}
          hasSubmitSucceeded={false}
        />,
      );
    });

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('.storage-in-transit-panel-modal')).toBe(true);
    });

    it('buttons are enabled', () => {
      expect(wrapper.find('button.usa-button-primary').prop('disabled')).toBe(false);
    });

    it('clicking save calls saveEdit', () => {
      wrapper.find('button.usa-button-primary').simulate('click');
      expect(saveAndClose.mock.calls.length).toBe(1);
    });
  });
});
