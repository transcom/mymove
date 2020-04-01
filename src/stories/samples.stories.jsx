import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { storiesOf } from '@storybook/react';

import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';
import { Form } from '../components/form';
import { TextInput } from '../components/form/fields';

storiesOf('Samples|Form', module).add('personal info', () => (
  <div id="samples-orders-container" style={{ padding: '20px' }}>
    <div className="container container--accent--blue">
      <Formik
        initialValues={{ firstName: '', lastName: '', email: '' }}
        validationSchema={Yup.object({
          firstName: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          lastName: Yup.string().max(20, 'Must be 20 characters or less').required('Required'),
          email: Yup.string().email('Invalid email address').required('Required'),
        })}
        onSubmit={action('Form Submit')}
        onReset={action('Form Canceled')}
      >
        <Form>
          <TextInput name="firstName" label="First Name" type="text" />
          <TextInput name="lastName" label="Last Name" type="text" />
          <TextInput name="email" label="Email Address" type="email" />
          <Button type="submit">Submit</Button>
          <Button type="reset" secondary>
            Cancel
          </Button>
        </Form>
      </Formik>
    </div>
  </div>
));
