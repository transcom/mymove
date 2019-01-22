import React from 'react';
import { shallow } from 'enzyme';

import { Editor } from './Editor';

const tariff400ng_items = [
  {
    id: 'sdlfkj',
    code: 'F9D',
    item: 'Long Haul',
  },
  {
    id: 'badfka',
    code: '19D',
    item: 'Crate',
  },
];
const shipmentLineItem = {
  id: 'sldkjf',
  tariff400ng_item: tariff400ng_items[0],
  location: 'D',
  quantity_1: 167000,
  notes: '',
  created_at: '2018-09-24T14:05:38.847Z',
  status: 'SUBMITTED',
};
let wrapper;
const cancel = jest.fn();
const save = jest.fn();
const saveComplete = jest.fn();
const submit = jest.fn();
const clear = jest.fn();
describe('given an Editor', () => {
  describe('when the form is disabled', () => {
    beforeEach(() => {
      cancel.mockClear();
      wrapper = shallow(
        <Editor
          tariff400ngItems={tariff400ng_items}
          shipmentLineItem={shipmentLineItem}
          saveEdit={save}
          cancelEdit={cancel}
          onSaveComplete={saveComplete}
          formEnabled={false}
          hasSubmitSucceeded={false}
          submitForm={submit}
          clearForm={clear}
        />,
      );
    });

    it('renders without crashing', () => {
      expect(wrapper.exists('div')).toBe(true);
    });
    it('cancel closes the form', () => {
      wrapper.find('.cancel-link a').simulate('click');
      expect(cancel.mock.calls.length).toBe(1);
    });
    it('buttons are disabled', () => {
      expect(wrapper.find('button.usa-button-primary').prop('disabled')).toBeTruthy();
    });
  });
  describe('when the form is enabled', () => {
    beforeEach(() => {
      submit.mockClear();
      wrapper = shallow(
        <Editor
          tariff400ngItems={tariff400ng_items}
          shipmentLineItem={shipmentLineItem}
          saveEdit={save}
          cancelEdit={cancel}
          onSaveComplete={saveComplete}
          formEnabled={true}
          hasSubmitSucceeded={false}
          submitForm={submit}
          clearForm={clear}
        />,
      );
    });

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });
    it('cancel closes the form', () => {
      wrapper.find('.cancel-link a').simulate('click');
      expect(cancel.mock.calls.length).toBe(1);
    });
    it('buttons are enabled', () => {
      expect(wrapper.find('button.usa-button-primary').prop('disabled')).toBe(false);
    });
    it('clicking save calls saveEdit', () => {
      wrapper.find('button.usa-button-primary').simulate('click');
      expect(submit.mock.calls.length).toBe(1);
    });
  });
});
