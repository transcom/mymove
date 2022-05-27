import React from 'react';
import { shallow } from 'enzyme';
import { Button } from '@trussworks/react-uswds';

import { DocsButton, EditButton } from './index';

describe('DocsButton', () => {
  it('should render the button', () => {
    const wrapper = shallow(<DocsButton label="my docs button" />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).html()).toBe(
      '<button class="usa-button" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="file" class="svg-inline--fa fa-file fa-w-12 " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><path fill="currentColor" d="M224 136V0H24C10.7 0 0 10.7 0 24v464c0 13.3 10.7 24 24 24h336c13.3 0 24-10.7 24-24V160H248c-13.2 0-24-10.8-24-24zm160-14.1v6.1H256V0h6.1c6.4 0 12.5 2.5 17 7l97.9 98c4.5 4.5 7 10.6 7 16.9z"></path></svg></span><span>my docs button</span></button>',
    );
  });
  it('should pass props down', () => {
    const wrapper = shallow(<DocsButton label="my docs button" className="sample-class-name" />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).prop('className')).toBe('sample-class-name');
  });
  it('onClick works', () => {
    const mockFn = jest.fn();
    const wrapper = shallow(<DocsButton label="my docs button" onClick={mockFn} />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).prop('onClick')).toBe(mockFn);
    wrapper.simulate('click');
    expect(mockFn).toHaveBeenCalled();
  });
});

describe('EditButton', () => {
  it('should render the button', () => {
    const wrapper = shallow(<EditButton />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).html()).toBe(
      '<button class="usa-button" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="pen" class="svg-inline--fa fa-pen fa-w-16 " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M290.74 93.24l128.02 128.02-277.99 277.99-114.14 12.6C11.35 513.54-1.56 500.62.14 485.34l12.7-114.22 277.9-277.88zm207.2-19.06l-60.11-60.11c-18.75-18.75-49.16-18.75-67.91 0l-56.55 56.55 128.02 128.02 56.55-56.55c18.75-18.76 18.75-49.16 0-67.91z"></path></svg></span><span>Edit</span></button>',
    );
  });
  it('should render the button with custom label', () => {
    const wrapper = shallow(<EditButton label="my custom edit" />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).html()).toBe(
      '<button class="usa-button" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="pen" class="svg-inline--fa fa-pen fa-w-16 " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M290.74 93.24l128.02 128.02-277.99 277.99-114.14 12.6C11.35 513.54-1.56 500.62.14 485.34l12.7-114.22 277.9-277.88zm207.2-19.06l-60.11-60.11c-18.75-18.75-49.16-18.75-67.91 0l-56.55 56.55 128.02 128.02 56.55-56.55c18.75-18.76 18.75-49.16 0-67.91z"></path></svg></span><span>my custom edit</span></button>',
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
