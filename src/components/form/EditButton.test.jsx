import React from 'react';
import { shallow } from 'enzyme';
import { Button } from '@trussworks/react-uswds';
import { EditButton } from '.';

describe('EditButton', () => {
  it('should render the button', () => {
    const wrapper = shallow(<EditButton />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).html()).toBe(
      '<button class="usa-button usa-button--icon" data-testid="button"><span class="icon"><svg>edit.svg</svg></span><span>Edit</span></button>',
    );
  });
  it('should pass props down', () => {
    const wrapper = shallow(<EditButton className="sample-class-name" />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).prop('className')).toBe('sample-class-name');
  });
  it('onClick works', () => {
    const mockFn = jest.fn();
    const wrapper = shallow(<EditButton onClick={mockFn} />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).prop('onClick')).toBe(mockFn);
    wrapper.simulate('click');
    expect(mockFn).toHaveBeenCalled();
  });
});
