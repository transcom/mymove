import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import FeedbackForm from './FeedbackForm';

// Enzyme.configure({ adapter: new Adapter() });

describe('<FeedbackForm />', () => {
  it('displays the feedback form, testing with Enzyme/Jest', () => {
    const feedbackFormWrapper = shallow(<FeedbackForm />);
    expect(feedbackFormWrapper).toHaveLength(1);
  });
});
