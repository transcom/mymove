import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Fieldset, Checkbox, Radio } from '@trussworks/react-uswds';
import { action } from '@storybook/addon-actions';

export default {
  title: 'Components|Form',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/3f6f957c-aa9a-4ea4-a064-430e6624fd62?mode=design',
    },
  },
};

export const elements = () => (
  <div style={{ padding: '20px' }}>
    <hr />
    <h3>Form Elements</h3>
    <Formik
      initialValues={{ rejectionReason: '' }}
      validationSchema={Yup.object({
        rejectionReason: Yup.string().min(15, 'Must be 15 characters or more').required('Required'),
      })}
      onSubmit={action('Form Submit')}
      onReset={action('Form Canceled')}
    >
      <form className="usa-form">
        <label className="usa-label" htmlFor="input-type-text-example-1">
          Text input label
          <input className="usa-input" id="input-type-text-example-1" name="input-type-text-example-1" type="text" />
        </label>

        <label className="usa-label" htmlFor="input-focus">
          Text input focused
          <input className="usa-input usa-focus" id="input-focus" name="input-focus" type="text" />
        </label>

        <div className="usa-form-group usa-form-group--error">
          <label className="usa-label usa-label--error" htmlFor="input-error">
            Text input error
            <span className="usa-error-message" id="input-error-message" role="alert">
              Helpful error message
            </span>
            <input
              className="usa-input usa-input--error"
              id="input-error"
              name="input-error"
              type="text"
              aria-describedby="input-error-message"
            />
          </label>
        </div>

        <label className="usa-label" htmlFor="input-type-textarea">
          Text area label
          <textarea className="usa-textarea" id="input-type-textarea" name="input-type-textarea" />
        </label>

        <label className="usa-label" htmlFor="options">
          Dropdown label
          <select className="usa-select" name="options" id="options">
            <option value>- Select -</option>
            <option value="value1">Option A</option>
            <option value="value2">Option B</option>
            <option value="value3">Option C</option>
          </select>
        </label>
        <br />
        <Fieldset legend="Historical figures 1" legendSrOnly id="input-type-fieldset">
          <Checkbox defaultChecked id="truth" label="Sojourner Truth" name="historical-figures-1" value="truth" />
          <Checkbox id="douglass" label="Frederick Douglass" name="historical-figures-1" value="douglass" />
          <Checkbox id="washington" label="Booker T. Washington" name="historical-figures-1" value="washington" />
          <Checkbox disabled id="carver" label="George Washington Carver" name="historical-figures-1" />
        </Fieldset>

        <Fieldset legend="Historical figures 2" legendSrOnly id="radios-fieldset">
          <Radio
            defaultChecked
            id="stanton"
            label="Elizabeth Cady Stanton"
            name="historical-figures-2"
            value="stanton"
          />
          <Radio id="anthony" label="Susan B. Anthony" name="historical-figures-2" value="anthony" />
          <Radio id="tubman" label="Harriet Tubman" name="historical-figures-2" value="tubman" />
          <Radio disabled id="invalid" label="Invalid option" name="historical-figures-2" value="invalid" />
        </Fieldset>
      </form>
    </Formik>
  </div>
);
