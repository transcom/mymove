import React from 'react';
import { shallow } from 'enzyme';
import { ErrorMessage as UswdsErrorMessage } from '@trussworks/react-uswds';

import { ErrorMessage } from './index';

describe('ErrorMessage', () => {
  describe('with error and display true', () => {
    const wrapper = shallow(
      <ErrorMessage display className="sample-class">
        This field is required
      </ErrorMessage>,
    );

    it('should render the USWDS ErrorMessage', () => {
      expect(wrapper.find(UswdsErrorMessage).length).toBe(1);
    });

    it('should accept className', () => {
      expect(wrapper.prop('className')).toBe('sample-class');
    });

    it('should display the error message', () => {
      expect(wrapper.dive().text()).toBe('This field is required');
    });
  });
  describe('with display false', () => {
    it('should NOT render the USWDS ErrorMessage', () => {
      const wrapper = shallow(
        <ErrorMessage display={false} className="sample-class">
          This field is required
        </ErrorMessage>,
      );

      expect(wrapper.find(UswdsErrorMessage).length).toBe(0);
    });
  });
  describe('with no error message text', () => {
    it('should NOT render the USWDS ErrorMessage', () => {
      const wrapper = shallow(<ErrorMessage display className="sample-class" />);

      expect(wrapper.find(UswdsErrorMessage).length).toBe(0);
    });
  });
});
