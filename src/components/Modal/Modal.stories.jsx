import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions } from './Modal';

import { Form } from 'components/form';
import { TextInput } from 'components/form/fields';

export default {
  title: 'Components|Modals',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/e6db2b1a-4d3e-40e7-89a8-39a25ab28b9a?mode=design',
    },
  },
};

export const withContent = () => (
  <Modal>
    <ModalTitle>
      <ModalClose handleClick={action('Close modal')} />
      <h4>Are you sure you want to reject this request?</h4>
    </ModalTitle>
    <Formik
      initialValues={{ rejectionReason: '' }}
      validationSchema={Yup.object({
        rejectionReason: Yup.string().min(15, 'Must be 15 characters or more').required('Required'),
      })}
      onSubmit={action('Form Submit')}
      onReset={action('Form Canceled')}
    >
      <Form>
        <TextInput name="rejectionReason" label="Reason for rejection" type="text" />
        <ModalActions>
          <Button type="submit">Confirm</Button>
          <Button secondary type="reset">
            Cancel
          </Button>
        </ModalActions>
      </Form>
    </Formik>
  </Modal>
);

export const empty = () => (
  <Modal>
    <ModalTitle>
      <ModalClose handleClick={action('Close modal')} />
      <h4>Modal title</h4>
    </ModalTitle>
    <ModalActions>
      <Button type="button">Submit</Button>
      <Button secondary type="button">
        Back
      </Button>
    </ModalActions>
  </Modal>
);
