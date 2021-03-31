import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions } from './Modal';

import { Form } from 'components/form';
import TextField from 'components/form/fields/TextField';

export default {
  title: 'Components/Modals',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/e6db2b1a-4d3e-40e7-89a8-39a25ab28b9a?mode=design',
    },
  },
};

export const withContent = () => (
  <Modal>
    <ModalClose handleClick={action('Close modal')} />
    <ModalTitle>
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
        <TextField name="rejectionReason" label="Reason for rejection" type="text" />
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
    <ModalClose handleClick={action('Close modal')} />
    <ModalTitle>
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

export const contentNoTitle = () => (
  <Modal>
    <ModalClose handleClick={action('Close modal')} />
    <h4>
      <strong>Long-term storage (NTS)</strong>
    </h4>
    <p>
      Put some or all of your things into storage as part of one move, and get it out of storage on a future move. Your
      move counselor can verify whether or not you qualify to put things into long-term storage on this move.
    </p>
    <ul>
      <li>The weight of this shipment counts against your weight allowance</li>
      <li>Useful when you can’t take all your things to your new location</li>
      <li>Common in OCONUS moves, but may not be available in CONUS</li>
      <li>Stored in a government-approved facility, typically near your starting location</li>
    </ul>
    <p>
      NTS (short for “non-temp storage”) lasts 6 months or longer. Do not count on easy access to things in storage. You
      can retrieve them during a future move.
    </p>
    <ModalActions>
      <Button type="button">Submit</Button>
      <Button secondary type="button">
        Back
      </Button>
    </ModalActions>
  </Modal>
);
