import React from 'react';
import { Provider } from 'react-redux';
import store from 'shared/store';
import { mount } from 'enzyme';

import StoragePanel from './StoragePanel';
import Alert from 'shared/Alert';
import { PanelSwaggerField } from 'shared/EditablePanel';

describe('StoragePanel', () => {
  describe('Receipts Awaiting Review Alert', () => {
    it('is non-existent when all receipts have OK status', () => {
      const moveDocuments = [
        { moving_expense_type: 'STORAGE', status: 'OK' },
        { moving_expense_type: 'STORAGE', status: 'OK' },
      ];
      const moveId = 'some ID';

      const wrapper = mount(
        <Provider store={store}>
          <StoragePanel title="Storage" moveId={moveId} moveDocuments={moveDocuments} />
        </Provider>,
      );
      expect(wrapper.find(PanelSwaggerField)).toHaveLength(2);
      expect(wrapper.find(Alert)).toHaveLength(0);
    });
    it('is existent when any receipt does not have an have OK status', () => {
      const moveDocuments = [
        { moving_expense_type: 'STORAGE', status: 'OK' },
        { moving_expense_type: 'STORAGE', status: 'HAS_ISSUE' },
      ];
      const moveId = 'some ID';

      const wrapper = mount(
        <Provider store={store}>
          <StoragePanel title="Storage" moveId={moveId} moveDocuments={moveDocuments} />
        </Provider>,
      );
      expect(wrapper.find(PanelSwaggerField)).toHaveLength(2);
      expect(wrapper.find(Alert)).toHaveLength(1);
    });
  });
});
