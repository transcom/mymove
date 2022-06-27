import React from 'react';
import { shallow } from 'enzyme';
import { Button } from '@trussworks/react-uswds';

import { DocsButton, EditButton } from './index';

describe('DocsButton', () => {
  it('should render the button', () => {
    const wrapper = shallow(<DocsButton label="my docs button" />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).html()).toBe(
      '<button class="usa-button" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="file" class="svg-inline--fa fa-file " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><path fill="currentColor" d="M0 64C0 28.65 28.65 0 64 0H224V128C224 145.7 238.3 160 256 160H384V448C384 483.3 355.3 512 320 512H64C28.65 512 0 483.3 0 448V64zM256 128V0L384 128H256z"></path></svg></span><span>my docs button</span></button>',
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
      '<button class="usa-button" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="pen" class="svg-inline--fa fa-pen " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M362.7 19.32C387.7-5.678 428.3-5.678 453.3 19.32L492.7 58.75C517.7 83.74 517.7 124.3 492.7 149.3L444.3 197.7L314.3 67.72L362.7 19.32zM421.7 220.3L188.5 453.4C178.1 463.8 165.2 471.5 151.1 475.6L30.77 511C22.35 513.5 13.24 511.2 7.03 504.1C.8198 498.8-1.502 489.7 .976 481.2L36.37 360.9C40.53 346.8 48.16 333.9 58.57 323.5L291.7 90.34L421.7 220.3z"></path></svg></span><span>Edit</span></button>',
    );
  });
  it('should render the button with custom label', () => {
    const wrapper = shallow(<EditButton label="my custom edit" />);
    expect(wrapper.find(Button).length).toBe(1);
    expect(wrapper.find(Button).html()).toBe(
      '<button class="usa-button" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="pen" class="svg-inline--fa fa-pen " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M362.7 19.32C387.7-5.678 428.3-5.678 453.3 19.32L492.7 58.75C517.7 83.74 517.7 124.3 492.7 149.3L444.3 197.7L314.3 67.72L362.7 19.32zM421.7 220.3L188.5 453.4C178.1 463.8 165.2 471.5 151.1 475.6L30.77 511C22.35 513.5 13.24 511.2 7.03 504.1C.8198 498.8-1.502 489.7 .976 481.2L36.37 360.9C40.53 346.8 48.16 333.9 58.57 323.5L291.7 90.34L421.7 220.3z"></path></svg></span><span>my custom edit</span></button>',
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
