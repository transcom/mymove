import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { storiesOf } from '@storybook/react';

import { Modal, Button } from '@trussworks/react-uswds';
import { Form } from '../components/form';
import { TextInput } from '../components/form/fields';

storiesOf('Components|Modals', module)
  .add('modal with content', () => (
    <div
      style={{
        margin: '100px',
      }}
    >
      <Modal
        className="modal container container--popout"
        title={<h4>Are you sure you want to reject this request?</h4>}
      >
        <Formik
          initialValues={{ rejectionReason: '' }}
          validationSchema={Yup.object({
            rejectionReason: Yup.string().min(15, 'Must be 15 characters or more').required('Required'),
          })}
        >
          <Form>
            <TextInput name="rejectionReason" label="Reason for rejection" type="text" />
          </Form>
        </Formik>
        <div className="display-flex">
          <Button type="button">Confirm</Button>
          <Button secondary type="button">
            Cancel
          </Button>
        </div>
      </Modal>
    </div>
  ))
  .add('empty modal', () => (
    <div
      style={{
        margin: '100px',
      }}
    >
      <Modal className="modal container container--popout" title={<h4>Modal title</h4>}>
        <div className="display-flex">
          <Button type="button">Submit</Button>
          <Button secondary type="button">
            Back
          </Button>
        </div>
      </Modal>
    </div>
  ));
