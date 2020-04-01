import React from 'react';
import { shallow } from 'enzyme';
import { Form as UswdsForm, Button } from '@trussworks/react-uswds';

describe('Form', () => {
  // mock out formik hook as we are not testing formik
  // needs to be before first describe
  jest.mock('formik', () => {
    return {
      useFormikContext: jest.fn().mockReturnValue({
        handleReset: jest.fn().mockName('handleReset'),
        handleSubmit: jest.fn().mockName('handleSubmit'),
      }),
    };
  });
  // require the above mock for expectations
  // eslint-disable-next-line global-require
  const mock = require('formik');

  // import component we are testing after mock created
  // eslint-disable-next-line global-require
  const { Form } = require('.');

  it('should render the USWDS Form', () => {
    const wrapper = shallow(
      <Form className="sample-class">
        <Button type="submit">Submit</Button>
      </Form>,
    );

    expect(wrapper.find(UswdsForm).length).toBe(1);
    expect(wrapper.find(Button).length).toBe(1);
    expect(mock.useFormikContext).toHaveBeenCalled();
    expect(wrapper.prop('onSubmit').getMockName()).toBe('handleSubmit');
    expect(wrapper.prop('onReset').getMockName()).toBe('handleReset');
    expect(wrapper.prop('className')).toBe('sample-class');
  });
});
