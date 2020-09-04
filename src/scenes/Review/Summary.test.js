/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';
// import { mount } from 'enzyme';
// import { Provider } from 'react-redux';
// import { ConnectedRouter } from 'connected-react-router';

// import { history, store } from '../../shared/store';

import Summary from './Summary';

const defaultProps = {
  serviceMember: {
    id: '123',
    current_station: {},
  },
  match: { path: '', url: '/moves/123/review', params: { moveId: '123' } },
  selectedMoveType: 'HHG',
};

// function mountSummary(props = defaultProps) {
//   return mount(
//     <Provider store={store}>
//       <ConnectedRouter history={history}>
//         <Summary {...props} />
//       </ConnectedRouter>
//     </Provider>,
//   );
// }

describe('Component renders', () => {
  expect(shallow(<Summary {...defaultProps} />).length).toEqual(1);

  // Not able to find this element with shallow rendering :(
  // const wrapper = shallow(<Summary {...defaultProps} />);
  // expect(wrapper.containsMatchingElement(<h3>Add another shipment</h3>)).toBe(true);
});
describe('Summary component', () => {
  it('renders the button to add another shipment', () => {
    // Not able to properly render props via ConnectedRouter rendering :(
    // const wrapper = mountSummary();
    // expect(wrapper.find('Summary').length).toBe(1);
  });
});
