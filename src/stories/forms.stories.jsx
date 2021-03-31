import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import {
  Form,
  Fieldset,
  Checkbox,
  Radio,
  TextInput,
  Label,
  FormGroup,
  ErrorMessage,
  Grid,
  Button,
  Textarea,
  Dropdown,
  GridContainer,
} from '@trussworks/react-uswds';
import { action } from '@storybook/addon-actions';

import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Hint from 'components/Hint/index';

export default {
  title: 'Components/Forms',
};

export const TextFieldset = () => (
  <Form className={formStyles.form}>
    <Fieldset>
      <FormGroup>
        <Label>Text Input Label</Label>
        <TextInput />
      </FormGroup>
    </Fieldset>
  </Form>
);

export const TextFieldsetDisabled = () => (
  <Form className={formStyles.form}>
    <Fieldset>
      <FormGroup>
        <Label>Text Input Label</Label>
        <TextInput value="This cannot be edited" disabled />
      </FormGroup>
    </Fieldset>
  </Form>
);

export const TextFieldsetWithError = () => (
  <Form className={formStyles.form}>
    <Fieldset>
      <FormGroup error>
        <Label>Text Input Label</Label>
        <ErrorMessage>This input has an error</ErrorMessage>
        <TextInput error />
      </FormGroup>
    </Fieldset>
  </Form>
);

export const InlineFieldset = () => (
  <Form className={formStyles.form}>
    <p>This form uses Grid components to position multiple inputs next to each other</p>
    <Fieldset>
      <Grid row gap>
        <Grid tablet={{ col: 'fill' }}>
          <FormGroup>
            <Label>Signature</Label>
            <TextInput />
          </FormGroup>
        </Grid>
        <Grid tablet={{ col: 'auto' }}>
          <FormGroup>
            <Label>Date</Label>
            <TextInput value="1/20/2021" disabled />
          </FormGroup>
        </Grid>
      </Grid>
    </Fieldset>
  </Form>
);

export const FormElements = () => (
  <GridContainer>
    <Formik
      initialValues={{ rejectionReason: '' }}
      validationSchema={Yup.object({
        rejectionReason: Yup.string().min(15, 'Must be 15 characters or more').required('Required'),
      })}
      onSubmit={action('Form Submit')}
      onReset={action('Form Canceled')}
    >
      <Form>
        <Label htmlFor="input-type-text-example-1">
          Text input label
          <TextInput id="input-type-text-example-1" name="input-type-text-example-1" type="text" />
        </Label>

        <Label htmlFor="input-focus">
          Text input focused
          <TextInput className="usa-focus" id="input-focus" name="input-focus" type="text" />
        </Label>

        <FormGroup error>
          <Label error htmlFor="input-error">
            Text input error
            <ErrorMessage id="input-error-message" role="alert">
              Helpful error message
            </ErrorMessage>
            <TextInput error id="input-error" name="input-error" type="text" aria-describedby="input-error-message" />
          </Label>
        </FormGroup>

        <Label htmlFor="input-type-textarea">
          Text area label
          <Textarea id="input-type-textarea" name="input-type-textarea" />
        </Label>

        <Label htmlFor="options">
          Dropdown label
          <Dropdown name="options" id="options">
            <option value>- Select -</option>
            <option value="value1">Option A</option>
            <option value="value2">Option B</option>
            <option value="value3">Option C</option>
          </Dropdown>
        </Label>

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
      </Form>
    </Formik>
  </GridContainer>
);

export const KitchenSinkForm = () => (
  <Form className={formStyles.form}>
    <h1>This is an example form</h1>
    <h2>It has lots of different components</h2>
    <Fieldset legend="Fieldset containing a single field">
      <p>Some arbitrary text about this fieldset</p>
      <Label>Text Input Label</Label>
      <TextInput />
    </Fieldset>
    <Fieldset legend="Fieldset containing multiple fields">
      <Label>Text Input Label</Label>
      <TextInput />

      <Checkbox label="A checkbox" />

      <Label>Text Input 2 Label</Label>
      <TextInput />

      <Hint>
        <p>Some arbitrary hint text at the end of this fieldset</p>
      </Hint>
    </Fieldset>
    <Fieldset legend="Fieldset containing multiple fields in form groups">
      <FormGroup>
        <p>These radio inputs are stacked on top of each other and grouped together in a form group</p>
        <Radio label="Yes" />
        <Radio label="No" />
      </FormGroup>

      <FormGroup>
        <p>These radio inputs are displayed next to each other in a radio group</p>
        <div className={formStyles.radioGroup}>
          <Radio label="Yes" />
          <Radio label="No" />
        </div>
      </FormGroup>

      <FormGroup>
        <Label>Text Input Label</Label>
        <TextInput />
      </FormGroup>
      <FormGroup>
        <Label>Text Input 2 Label</Label>
        <Hint>
          <p>Some arbitrary hint text about this input</p>
        </Hint>
        <TextInput />
      </FormGroup>
    </Fieldset>
    <Fieldset>
      <p>
        This fieldset just contains static text and doesn’t even have a legend. Maybe it represents data that cannot
        change right now.
      </p>
    </Fieldset>
    <Fieldset
      legend={
        <div className={formStyles.legendContent}>
          Legend with optional label <span className={formStyles.optional}>Optional</span>
        </div>
      }
    >
      <FormGroup>
        <Label>Text Input Label</Label>
        <TextInput />
      </FormGroup>
    </Fieldset>
    <SectionWrapper className={formStyles.formSection}>
      <h1>Multiple sections</h1>
      <p>This form has some sub-sections</p>

      <Fieldset legend="Fieldset containing a single field">
        <p>Some arbitrary text about this fieldset</p>
        <Label>Text Input Label</Label>
        <TextInput />
      </Fieldset>

      <Fieldset legend="Fieldset containing multiple fields">
        <Label>Text Input Label</Label>
        <TextInput />

        <Checkbox label="A checkbox" />

        <Label>Text Input 2 Label</Label>
        <TextInput />

        <Hint>
          <p>Some arbitrary hint text at the end of this fieldset</p>
        </Hint>
      </Fieldset>

      <Fieldset legend="Fieldset containing multiple fields in form groups">
        <FormGroup>
          <p>These radio inputs are stacked on top of each other and grouped together in a form group</p>
          <Radio label="Yes" />
          <Radio label="No" />
        </FormGroup>

        <FormGroup>
          <p>These radio inputs are displayed next to each other in a radio group</p>
          <div className={formStyles.radioGroup}>
            <Radio label="Yes" />
            <Radio label="No" />
          </div>
        </FormGroup>

        <FormGroup>
          <Label>Text Input Label</Label>
          <TextInput />
        </FormGroup>
        <FormGroup>
          <Label>Text Input 2 Label</Label>
          <Hint>
            <p>Some arbitrary hint text about this input</p>
          </Hint>
          <TextInput />
        </FormGroup>
      </Fieldset>
    </SectionWrapper>

    <SectionWrapper className={formStyles.formSection}>
      <h1>Section 2!</h1>
      <p>Here is everything again, but inside a section wrapper this time.</p>

      <Fieldset>
        <p>This fieldset just contains static text. Maybe it represents data that cannot change right now.</p>
      </Fieldset>

      <Fieldset
        legend={
          <div className={formStyles.legendContent}>
            Legend with optional label <span className={formStyles.optional}>Optional</span>
          </div>
        }
      >
        <FormGroup>
          <Label>Text Input Label</Label>
          <TextInput />
        </FormGroup>
      </Fieldset>

      <FormGroup>
        <Label>Text Input not in a fieldset</Label>
        <TextInput />
      </FormGroup>
    </SectionWrapper>
    <div className={formStyles.formActions}>
      <Button type="button">Submit</Button>
      <Button type="button" secondary outline>
        Cancel
      </Button>
    </div>
  </Form>
);

export const KitchenSinkFormWithErrors = () => (
  <Form className={formStyles.form}>
    <h1>This is an example form</h1>
    <h2>It has lots of different components</h2>
    <Fieldset legend="Fieldset containing a single field">
      <p>Some arbitrary text about this fieldset</p>
      <Label error>Text Input Label</Label>
      <TextInput error />
    </Fieldset>
    <Fieldset legend="Fieldset containing multiple fields">
      <Label error>Text Input Label</Label>
      <TextInput error />

      <Checkbox error label="A checkbox" />

      <Label error>Text Input 2 Label</Label>
      <TextInput error />

      <Hint>
        <p>Some arbitrary hint text at the end of this fieldset</p>
      </Hint>
    </Fieldset>
    <Fieldset legend="Fieldset containing multiple fields in form groups">
      <FormGroup error>
        <p>These radio inputs are stacked on top of each other and grouped together in a form group</p>
        <Radio error label="Yes" />
        <Radio error label="No" />
      </FormGroup>

      <FormGroup error>
        <p>These radio inputs are displayed next to each other in a radio group</p>
        <div className={formStyles.radioGroup}>
          <Radio error label="Yes" />
          <Radio error label="No" />
        </div>
      </FormGroup>

      <FormGroup error>
        <Label error>Text Input Label</Label>
        <TextInput error />
      </FormGroup>
      <FormGroup error>
        <Label error>Text Input 2 Label</Label>
        <Hint>
          <p>Some arbitrary hint text about this input</p>
        </Hint>
        <TextInput error />
      </FormGroup>
    </Fieldset>
    <Fieldset>
      <p>
        This fieldset just contains static text and doesn’t even have a legend. Maybe it represents data that cannot
        change right now.
      </p>
    </Fieldset>
    <Fieldset
      legend={
        <div className={formStyles.legendContent}>
          Legend with optional label <span className={formStyles.optional}>Optional</span>
        </div>
      }
    >
      <FormGroup error>
        <Label error>Text Input Label</Label>
        <TextInput error />
      </FormGroup>
    </Fieldset>
    <SectionWrapper className={formStyles.formSection}>
      <h1>Multiple sections</h1>
      <p>This form has some sub-sections</p>

      <Fieldset legend="Fieldset containing a single field">
        <p>Some arbitrary text about this fieldset</p>
        <Label error>Text Input Label</Label>
        <TextInput error />
      </Fieldset>

      <Fieldset legend="Fieldset containing multiple fields">
        <Label error>Text Input Label</Label>
        <TextInput error />

        <Checkbox error label="A checkbox" />

        <Label error>Text Input 2 Label</Label>
        <TextInput error />

        <Hint>
          <p>Some arbitrary hint text at the end of this fieldset</p>
        </Hint>
      </Fieldset>

      <Fieldset legend="Fieldset containing multiple fields in form groups">
        <FormGroup error>
          <p>These radio inputs are stacked on top of each other and grouped together in a form group</p>
          <Radio error label="Yes" />
          <Radio error label="No" />
        </FormGroup>

        <FormGroup error>
          <p>These radio inputs are displayed next to each other in a radio group</p>
          <div className={formStyles.radioGroup}>
            <Radio error label="Yes" />
            <Radio error label="No" />
          </div>
        </FormGroup>

        <FormGroup error>
          <Label error>Text Input Label</Label>
          <TextInput error />
        </FormGroup>
        <FormGroup error>
          <Label error>Text Input 2 Label</Label>
          <Hint>
            <p>Some arbitrary hint text about this input</p>
          </Hint>
          <TextInput error />
        </FormGroup>
      </Fieldset>
    </SectionWrapper>

    <SectionWrapper className={formStyles.formSection}>
      <h1>Section 2!</h1>
      <p>Here is everything again, but inside a section wrapper this time.</p>

      <Fieldset>
        <p>This fieldset just contains static text. Maybe it represents data that cannot change right now.</p>
      </Fieldset>

      <Fieldset
        legend={
          <div className={formStyles.legendContent}>
            Legend with optional label <span className={formStyles.optional}>Optional</span>
          </div>
        }
      >
        <FormGroup error>
          <Label error>Text Input Label</Label>
          <TextInput error />
        </FormGroup>
      </Fieldset>

      <FormGroup error>
        <Label error>Text Input not in a fieldset</Label>
        <TextInput error />
      </FormGroup>
    </SectionWrapper>
    <div className={formStyles.formActions}>
      <Button type="button">Submit</Button>
      <Button type="button" secondary outline>
        Cancel
      </Button>
    </div>
  </Form>
);
