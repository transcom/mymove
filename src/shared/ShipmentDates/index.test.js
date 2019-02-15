import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';

import Dates from '.';

describe('DatesPanel tests', () => {
  let wrapper;

  const shipment = {
    pm_survey_conducted_date: '',
    pm_survey_planned_pack_date: '',
    pm_survey_planned_pickup_date: '',
    pm_survey_planned_delivery_date: '',
    requested_pickup_date: '',
    actual_pack_date: '',
    actual_pickup_date: '',
    actual_delivery_date: '',
    requested_delivery_date: '',
    pm_survey_notes: '',
    pm_survey_method: '',
  };
  const title = 'Some Dates title';
  const update = () => {};

  const mockStore = configureStore();

  let store;

  beforeEach(() => {
    store = mockStore({ requests: { lastErrors: [] } });
    //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
    wrapper = mount(
      <Provider store={store}>
        <Dates title={title} shipment={shipment} update={update} />
      </Provider>,
    );
  });

  it('PanelField renders without crashing', () => {
    expect(wrapper.find('.editable-panel').length).toEqual(1);
  });

  it('includes the title of the document', () => {
    expect(wrapper.find('.editable-panel').text()).toContain(title);
  });
});
