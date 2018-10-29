import React from 'react';
import { shallow } from 'enzyme';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';
import store from 'shared/store';
import Dates from '.';

describe('DatesPanel tests', () => {
  let shipment = {
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
  const renderDatesPanel = () =>
    shallow(
      <Provider store={store}>
        <MockRouter>
          <Dates title={'Some Dates title'} shipment={shipment} update={''} />,
        </MockRouter>
      </Provider>,
    );
  it('renders without crashing', () => {
    const wrapper = renderDatesPanel();
    expect(wrapper.find('.editable-panel').length).toEqual(1);
  });

  // it('includes the title of the document', () => {
  //   const title = 'My Title';
  //   const documentPanel = renderDocumentPanel({ title });
  //   expect(documentPanel.find('.panel-subhead').text()).toContain(title);
  // });
});
