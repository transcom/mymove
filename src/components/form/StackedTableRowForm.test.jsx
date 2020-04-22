import React from 'react';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { mount } from 'enzyme';
import { EditButton, StackedTableRowForm } from '.';

describe('StackedTableRowForm', () => {
  const renderStackedTableRowForm = (submit, reset) => {
    return mount(
      <table>
        <tbody>
          <StackedTableRowForm
            initialValues={{ fieldName: 'value' }}
            validationSchema={Yup.object({
              ordersNumber: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
            })}
            onSubmit={submit}
            onReset={reset}
            name="fieldName"
            type="text"
            label="Field Name"
          />
        </tbody>
      </table>,
    );
  };

  it('renders a tr with correct html', () => {
    const component = renderStackedTableRowForm();
    expect(component.html()).toBe(
      '<table><tbody><tr class="default-table-row-classes"><th scope="row" class="default-table-header-class-names">Field Name</th><td class="default-table-data-class-names"><span>value</span><button type="button" class="usa-button usa-button--icon usa-button--unstyled" data-testid="button"><span class="icon"><svg>edit.svg</svg></span><span>Edit</span></button></td></tr></tbody></table>',
    );
  });

  it('renders a tr with correct html form', () => {
    const setShow = jest.fn();
    const setErrors = jest.fn();
    jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, setErrors]);
    const component = renderStackedTableRowForm();
    expect(component.html()).toBe(
      '<table><tbody><tr class="default-table-row-classes"><th scope="row" class="default-table-header-class-names"><label data-testid="label" class="usa-label" for="fieldName">Field Name</label></th><td class="default-table-data-class-names"><form data-testid="form" class="usa-form"><input data-testid="textInput" class="usa-input" name="fieldName" type="text" value="value"><div class="display-flex"><button type="submit" class="usa-button" data-testid="button">Submit</button><button type="reset" class="usa-button usa-button--secondary" data-testid="button">Cancel</button></div></form></td></tr></tbody></table>',
    );
  });

  it('toggles show on click of edit button', () => {
    const setShow = jest.fn();
    const setErrors = jest.fn();
    jest.spyOn(React, 'useState').mockReturnValueOnce([false, setShow]).mockReturnValueOnce([{}, setErrors]);
    const component = renderStackedTableRowForm();
    expect(component.find(EditButton).length).toBe(1);
    component.find(EditButton).simulate('click');
    expect(setShow).toHaveBeenCalledWith(true);
  });

  describe('event submit', () => {
    const setShow = jest.fn();
    const setErrors = jest.fn();
    jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, setErrors]);
    const submit = jest.fn();
    const reset = jest.fn();

    const component = renderStackedTableRowForm(submit, reset);

    it('does not show the EditButton with the form', () => {
      expect(component.find(EditButton).length).toBe(0);
    });

    it('triggers submit handler', () => {
      const onSubmitFunc = component.find(Formik).props().onSubmit;
      onSubmitFunc({ sample: 'submit data' });
      expect(submit).toHaveBeenCalledWith({
        sample: 'submit data',
      });
      expect(reset).not.toHaveBeenCalled();
    });

    it('resets show', () => {
      const onSubmitFunc = component.find(Formik).props().onSubmit;
      onSubmitFunc({ sample: 'submit data' });
      expect(setShow).toHaveBeenCalledWith(false);
    });
  });

  describe('event reset', () => {
    const setShow = jest.fn();
    const setErrors = jest.fn();
    jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, setErrors]);
    const submit = jest.fn();
    const reset = jest.fn();

    const component = renderStackedTableRowForm(submit, reset);

    it('does not show the EditButton with the form', () => {
      expect(component.find(EditButton).length).toBe(0);
    });

    it('triggers reset handler', () => {
      const onResetFunc = component.find(Formik).props().onReset;
      onResetFunc({ sample: 'reset data' });
      expect(reset).toHaveBeenCalledWith({
        sample: 'reset data',
      });
      expect(submit).not.toHaveBeenCalled();
    });

    it('resets show', () => {
      const onResetFunc = component.find(Formik).props().onReset;
      onResetFunc({ sample: 'reset data' });
      expect(setShow).toHaveBeenCalledWith(false);
    });
  });

  afterEach(jest.clearAllMocks);
});
