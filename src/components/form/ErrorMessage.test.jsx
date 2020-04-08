import React from 'react';
import { shallow, mount } from 'enzyme';
import { ErrorMessage as UswdsErrorMessage } from '@trussworks/react-uswds';
import { ErrorMessage } from '.';

describe('ErrorMessage', () => {
  describe('with error and display true', () => {
    const wrapper = mount(
      <ErrorMessage display className="sample-class">
        Error
      </ErrorMessage>,
    );

    it('should render the USWDS ErrorMessage', () => {
      expect(wrapper.find(UswdsErrorMessage).length).toBe(1);
    });

    it('should accept className', () => {
      expect(wrapper.prop('className')).toBe('sample-class');
    });

    it('should display the error message', () => {
      expect(wrapper.text()).toBe('Error');
    });
  });
  describe('with display false', () => {
    it('should NOT render the USWDS ErrorMessage', () => {
      const wrapper = shallow(
        <ErrorMessage display={false} className="sample-class">
          Error
        </ErrorMessage>,
      );

      expect(wrapper.find(UswdsErrorMessage).length).toBe(0);
    });
  });
});
