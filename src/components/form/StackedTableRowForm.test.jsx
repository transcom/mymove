import React from 'react';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { mount } from 'enzyme';

import { EditButton, ErrorMessage, Form, StackedTableRowForm } from './index';

describe('StackedTableRowForm', () => {
  const renderStackedTableRowForm = (submit = jest.fn(), reset, value = 'value') => {
    return mount(
      <table className="table--stacked">
        <tbody>
          <StackedTableRowForm
            initialValues={{ fieldName: value }}
            validationSchema={Yup.object({
              ordersNumber: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
            })}
            onSubmit={submit}
            onReset={reset}
            name="fieldName"
            type="text"
            label="Field Name"
            id="fieldName"
          />
        </tbody>
      </table>,
    );
  };

  describe('when NOT showing form', () => {
    it('renders a tr with correct html', () => {
      const component = renderStackedTableRowForm();
      expect(component.html()).toBe(
        '<table class="table--stacked"><tbody><tr class="stacked-table-row"><th scope="row" class="label ">Field Name</th><td><span>value</span><button type="button" class="usa-button usa-button--icon usa-button--unstyled float-right" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="pen" class="svg-inline--fa fa-pen fa-w-16 " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M290.74 93.24l128.02 128.02-277.99 277.99-114.14 12.6C11.35 513.54-1.56 500.62.14 485.34l12.7-114.22 277.9-277.88zm207.2-19.06l-60.11-60.11c-18.75-18.75-49.16-18.75-67.91 0l-56.55 56.55 128.02 128.02 56.55-56.55c18.75-18.76 18.75-49.16 0-67.91z"></path></svg></span><span>Edit</span></button></td></tr></tbody></table>',
      );
    });

    it('renders a span with nbsp when no value', () => {
      const component = renderStackedTableRowForm(jest.fn(), null, null);
      expect(component.html()).toBe(
        '<table class="table--stacked"><tbody><tr class="stacked-table-row"><th scope="row" class="label ">Field Name</th><td><span>&nbsp;</span><button type="button" class="usa-button usa-button--icon usa-button--unstyled float-right" data-testid="button"><span class="icon"><svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="pen" class="svg-inline--fa fa-pen fa-w-16 " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M290.74 93.24l128.02 128.02-277.99 277.99-114.14 12.6C11.35 513.54-1.56 500.62.14 485.34l12.7-114.22 277.9-277.88zm207.2-19.06l-60.11-60.11c-18.75-18.75-49.16-18.75-67.91 0l-56.55 56.55 128.02 128.02 56.55-56.55c18.75-18.76 18.75-49.16 0-67.91z"></path></svg></span><span>Edit</span></button></td></tr></tbody></table>',
      );
    });

    it('toggles show on click of edit button', () => {
      const setShow = jest.fn();
      jest.spyOn(React, 'useState').mockReturnValueOnce([false, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      const component = renderStackedTableRowForm();
      expect(component.find(EditButton).length).toBe(1);
      component.find(EditButton).simulate('click');
      expect(setShow).toHaveBeenCalledWith(true);
    });

    it('does not add error class th', () => {
      const component = renderStackedTableRowForm();
      expect(component.find('th').props().className).toBe('label ');
    });

    it('does set display on ErrorMessage', () => {
      const component = renderStackedTableRowForm();
      expect(component.find(ErrorMessage).length).toBe(1);
      expect(component.find(ErrorMessage).props().display).toBe(false);
    });
    describe('with errors', () => {
      beforeEach(() => {
        jest
          .spyOn(React, 'useState')
          .mockReturnValueOnce([false, jest.fn()])
          .mockReturnValueOnce([{ fieldName: 'Required' }, jest.fn()]);
      });

      it('adds error class th', () => {
        const component = renderStackedTableRowForm();
        expect(component.find('th').props().className).toBe('label error');
      });

      it('displays the error message and value', () => {
        const component = renderStackedTableRowForm();
        expect(component.find(ErrorMessage).length).toBe(1);
        expect(component.find(ErrorMessage).props().className).toBe('display-inline');
        expect(component.find(ErrorMessage).props().display).toBe(true);
        expect(component.find(ErrorMessage).props().children).toBe('Required');
      });
    });
  });

  describe('when showing form', () => {
    const setShow = jest.fn();
    const submit = jest.fn();
    const reset = jest.fn();

    it('renders a tr with correct html form', () => {
      jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      const component = renderStackedTableRowForm();
      expect(component.html()).toBe(
        '<table class="table--stacked"><tbody><tr class="stacked-table-row"><th scope="row" class="label ">Field Name</th><td><form data-testid="form" class="usa-form" role="form"><div data-testid="formGroup" class="usa-form-group"><label data-testid="label" class="usa-label" for="fieldName"></label><input data-testid="textInput" class="usa-input" id="fieldName" name="fieldName" value="value"></div><div class="form-buttons"><button type="submit" class="usa-button" data-testid="button">Submit</button><button type="reset" class="usa-button usa-button--secondary" data-testid="button">Cancel</button></div></form></td></tr></tbody></table>',
      );
    });

    describe('event submit', () => {
      beforeEach(() => {
        jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      });

      it('does not show the EditButton with the form', () => {
        const component = renderStackedTableRowForm(submit, reset);
        expect(component.find(EditButton).length).toBe(0);
      });

      it('triggers submit handler', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onSubmitFunc = component.find(Formik).props().onSubmit;
        onSubmitFunc({ sample: 'submit data' });
        expect(submit).toHaveBeenCalledWith({
          sample: 'submit data',
        });
        expect(reset).not.toHaveBeenCalled();
      });

      it('resets show', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onSubmitFunc = component.find(Formik).props().onSubmit;
        onSubmitFunc({ sample: 'submit data' });
        expect(setShow).toHaveBeenCalledWith(false);
      });
    });

    describe('event reset', () => {
      beforeEach(() => {
        jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      });

      it('does not show the EditButton with the form', () => {
        const component = renderStackedTableRowForm(submit, reset);
        expect(component.find(EditButton).length).toBe(0);
      });

      it('triggers reset handler', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onResetFunc = component.find(Formik).props().onReset;
        onResetFunc({ sample: 'reset data' });
        expect(reset).toHaveBeenCalledWith({
          sample: 'reset data',
        });
        expect(submit).not.toHaveBeenCalled();
      });

      it('resets show', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onResetFunc = component.find(Formik).props().onReset;
        onResetFunc({ sample: 'reset data' });
        expect(setShow).toHaveBeenCalledWith(false);
      });
    });
    describe('with errors', () => {
      const setErrors = jest.fn();

      beforeEach(() => {
        jest
          .spyOn(React, 'useState')
          .mockReturnValueOnce([true, jest.fn()])
          .mockReturnValueOnce([{ fieldName: 'Required' }, setErrors]);
      });

      it('passes in an errorCallback to Form', () => {
        const component = renderStackedTableRowForm();
        expect(component.find(Form).length).toBe(1);
        component.find(Form).props().errorCallback({ fieldName: 'Required' });
        expect(setErrors).toHaveBeenCalledWith({ fieldName: 'Required' });
      });

      it('adds error class th', () => {
        const component = renderStackedTableRowForm();
        expect(component.find('th').props().className).toBe('label error');
      });
    });
  });

  afterEach(jest.clearAllMocks);
});
